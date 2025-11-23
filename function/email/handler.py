import json
import os
import sys
import smtplib
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
from pathlib import Path
from datetime import datetime, timezone
from kubernetes import client, config
from kubernetes.client.rest import ApiException

def log(message, level="INFO"):
    """Log message to stderr with timestamp and level."""
    timestamp = datetime.now(timezone.utc).isoformat()
    print(f"[{timestamp}] [{level}] {message}", file=sys.stderr, flush=True)

def handle(event, context):
    """
    OpenFaaS function to send email based on Email CR JSON representation.
    Expects POST body with Email spec fields: toAddress, fromName, fromAddress, subject, body.
    Patches the Email CR status directly using Kubernetes API.
    """
    log("Email handler invoked")
    
    try:
        # Validate HTTP method is POST
        method = event.method if hasattr(event, 'method') else os.getenv('Http_Method', 'POST')
        log(f"HTTP method: {method}")
        
        if method.upper() != 'POST':
            log(f"Invalid method {method}, returning 405", "WARN")
            return {
                "statusCode": 405,
                "body": json.dumps({"error": f"Method {method} not allowed. Only POST is supported."})
            }
        
        # Parse incoming JSON body
        body = event.body
        if isinstance(body, bytes):
            body = body.decode('utf-8')
        
        log(f"Received body length: {len(body)} bytes")
        
        email_data = json.loads(body)
        log(f"Successfully parsed JSON")
        
        # Extract Email metadata
        metadata = email_data.get('metadata', {})
        name = metadata.get('name')
        namespace = metadata.get('namespace', 'default')
        generation = metadata.get('generation', 0)
        
        log(f"Email CR: {namespace}/{name}, generation: {generation}")
        
        if not name:
            error_msg = "Missing metadata.name in Email CR"
            log(error_msg, "ERROR")
            return {
                "statusCode": 400,
                "body": json.dumps({"error": error_msg})
            }
        
        # Extract Email spec fields
        spec = email_data.get('spec', {})
        to_address = spec.get('toAddress')
        from_name = spec.get('fromName')
        from_address = spec.get('fromAddress')
        subject = spec.get('subject')
        body_content = spec.get('body')
        
        log(f"Email details: to={to_address}, from={from_name} <{from_address}>, subject={subject}")
        
        # Validate required fields
        if not all([to_address, from_name, from_address, subject, body_content]):
            error_msg = "Missing required fields in Email spec"
            log(error_msg, "ERROR")
            log(f"Present fields: toAddress={bool(to_address)}, fromName={bool(from_name)}, fromAddress={bool(from_address)}, subject={bool(subject)}, body={bool(body_content)}", "ERROR")
            patch_email_status(
                name=name,
                namespace=namespace,
                generation=generation,
                error_message=error_msg,
                error_timestamp=datetime.now(timezone.utc)
            )
            return {
                "statusCode": 400,
                "body": json.dumps({
                    "error": error_msg,
                    "required": ["toAddress", "fromName", "fromAddress", "subject", "body"]
                })
            }
        
        # Send email
        try:
            log(f"Attempting to send email via SMTP")
            send_email(
                to_address=to_address,
                from_name=from_name,
                from_address=from_address,
                subject=subject,
                body=body_content
            )
            log("Email sent successfully via SMTP")
            
            # Patch Email CR status with success
            log(f"Patching Email CR status with success")
            patch_email_status(
                name=name,
                namespace=namespace,
                generation=generation,
                sent_timestamp=datetime.now(timezone.utc),
                error_message="",
                error_timestamp=None
            )
            log("Email CR status patched successfully")
            
            return {
                "statusCode": 200,
                "body": json.dumps({
                    "message": "Email sent successfully",
                    "to": to_address,
                    "subject": subject
                })
            }
            
        except Exception as smtp_error:
            error_msg = f"Failed to send email: {str(smtp_error)}"
            log(error_msg, "ERROR")
            log(f"SMTP error details: {type(smtp_error).__name__}", "ERROR")
            patch_email_status(
                name=name,
                namespace=namespace,
                generation=generation,
                error_message=error_msg,
                error_timestamp=datetime.now(timezone.utc)
            )
            return {
                "statusCode": 500,
                "body": json.dumps({"error": error_msg})
            }
        
    except json.JSONDecodeError as e:
        error_msg = f"Invalid JSON: {str(e)}"
        log(error_msg, "ERROR")
        return {
            "statusCode": 400,
            "body": json.dumps({"error": error_msg})
        }
    except Exception as e:
        error_msg = f"Unexpected error: {str(e)}"
        log(error_msg, "ERROR")
        log(f"Exception type: {type(e).__name__}", "ERROR")
        import traceback
        log(f"Traceback: {traceback.format_exc()}", "ERROR")
        return {
            "statusCode": 500,
            "body": json.dumps({"error": error_msg})
        }

def patch_email_status(name, namespace, generation, sent_timestamp=None, error_message=None, error_timestamp=None):
    """
    Patch the Email CR status using Kubernetes API.
    Uses in-cluster config from service account injection.
    """
    try:
        log(f"Loading Kubernetes in-cluster config")
        # Load in-cluster config from service account
        config.load_incluster_config()
        log("In-cluster config loaded successfully")
        
        # Create custom objects API client
        api = client.CustomObjectsApi()
        
        # Build status patch
        status_patch = {
            "lastGeneration": generation
        }
        
        if sent_timestamp:
            status_patch["sentTimestamp"] = sent_timestamp.isoformat().replace('+00:00', 'Z')
            log(f"Setting sentTimestamp: {status_patch['sentTimestamp']}")
        
        if error_message is not None:
            status_patch["errorMessage"] = error_message
            log(f"Setting errorMessage: {error_message}")
        
        if error_timestamp:
            status_patch["errorTimestamp"] = error_timestamp.isoformat().replace('+00:00', 'Z')
            log(f"Setting errorTimestamp: {status_patch['errorTimestamp']}")
        
        log(f"Patching Email CR {namespace}/{name} status subresource")
        # Patch the Email CR status subresource
        api.patch_namespaced_custom_object_status(
            group="product.webshop.harikube.info",
            version="v1",
            namespace=namespace,
            plural="emails",
            name=name,
            body={"status": status_patch}
        )
        log("Email CR status patch successful")
        
    except ApiException as e:
        # Log error but don't fail the function
        log(f"Failed to patch Email status: {e.status} {e.reason}", "ERROR")
        log(f"Response body: {e.body}", "ERROR")
    except Exception as e:
        log(f"Unexpected error patching Email status: {str(e)}", "ERROR")
        log(f"Exception type: {type(e).__name__}", "ERROR")

def send_email(to_address, from_name, from_address, subject, body):
    """
    Send email via SMTP.
    SMTP credentials are read from Kubernetes secret mounted at /var/openfaas/secrets/
    Expected secret keys:
    - smtp-host
    - smtp-port
    - smtp-username
    - smtp-password
    - smtp-use-tls (optional, defaults to "true")
    
    Fallback to environment variables if secrets are not found.
    """
    # Try to read from mounted secrets first, fallback to env vars
    smtp_host = read_secret('smtp-host') or os.getenv('SMTP_HOST', 'localhost')
    smtp_port = int(read_secret('smtp-port') or os.getenv('SMTP_PORT', '587'))
    smtp_username = read_secret('smtp-username') or os.getenv('SMTP_USERNAME')
    smtp_password = read_secret('smtp-password') or os.getenv('SMTP_PASSWORD')
    smtp_use_tls = (read_secret('smtp-use-tls') or os.getenv('SMTP_USE_TLS', 'true')).lower() == 'true'
    
    log(f"SMTP config: host={smtp_host}, port={smtp_port}, username={smtp_username}, use_tls={smtp_use_tls}")
    
    # Create message
    msg = MIMEMultipart('alternative')
    msg['Subject'] = subject
    msg['From'] = f"{from_name} <{from_address}>"
    msg['To'] = to_address
    
    # Add body (support both plain text and HTML)
    if body.strip().startswith('<'):
        # Likely HTML
        log("Attaching HTML body")
        msg.attach(MIMEText(body, 'html'))
    else:
        log("Attaching plain text body")
        msg.attach(MIMEText(body, 'plain'))
    
    log(f"Connecting to SMTP server {smtp_host}:{smtp_port}")
    # Send email
    with smtplib.SMTP(smtp_host, smtp_port) as server:
        if smtp_use_tls:
            log("Starting TLS")
            server.starttls()
        
        if smtp_username and smtp_password:
            log(f"Logging in as {smtp_username}")
            server.login(smtp_username, smtp_password)
        else:
            log("No SMTP credentials provided, skipping authentication")
        
        log("Sending message")
        server.send_message(msg)
        log("Message sent successfully")

def read_secret(key):
    """
    Read secret value from OpenFaaS secret mount path.
    Returns None if secret file doesn't exist.
    """
    secret_path = Path(f'/var/openfaas/secrets/{key}')
    try:
        if secret_path.exists():
            value = secret_path.read_text().strip()
            log(f"Read secret: {key} (length: {len(value)})")
            return value
        else:
            log(f"Secret not found: {key}", "WARN")
    except Exception as e:
        log(f"Error reading secret {key}: {str(e)}", "ERROR")
    return None
