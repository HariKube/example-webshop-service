import json
import os
import smtplib
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
from pathlib import Path
from datetime import datetime, timezone
from kubernetes import client, config
from kubernetes.client.rest import ApiException


def handle(event, context):
    """
    OpenFaaS function to send email based on Email CR JSON representation.
    Expects POST body with Email spec fields: toAddress, fromName, fromAddress, subject, body.
    Patches the Email CR status directly using Kubernetes API.
    """
    try:
        # Validate HTTP method is POST
        method = event.method if hasattr(event, 'method') else os.getenv('Http_Method', 'POST')
        if method.upper() != 'POST':
            return {
                "statusCode": 405,
                "body": json.dumps({"error": f"Method {method} not allowed. Only POST is supported."})
            }
        
        # Parse incoming JSON body
        body = event.body
        if isinstance(body, bytes):
            body = body.decode('utf-8')
        
        email_data = json.loads(body)
        
        # Extract Email metadata
        metadata = email_data.get('metadata', {})
        name = metadata.get('name')
        namespace = metadata.get('namespace', 'default')
        generation = metadata.get('generation', 0)
        
        if not name:
            error_msg = "Missing metadata.name in Email CR"
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
        
        # Validate required fields
        if not all([to_address, from_name, from_address, subject, body_content]):
            error_msg = "Missing required fields in Email spec"
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
            send_email(
                to_address=to_address,
                from_name=from_name,
                from_address=from_address,
                subject=subject,
                body=body_content
            )
            
            # Patch Email CR status with success
            patch_email_status(
                name=name,
                namespace=namespace,
                generation=generation,
                sent_timestamp=datetime.now(timezone.utc),
                error_message="",
                error_timestamp=None
            )
            
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
        return {
            "statusCode": 400,
            "body": json.dumps({"error": error_msg})
        }
    except Exception as e:
        error_msg = f"Unexpected error: {str(e)}"
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
        # Load in-cluster config from service account
        config.load_incluster_config()
        
        # Create custom objects API client
        api = client.CustomObjectsApi()
        
        # Build status patch
        status_patch = {
            "lastGeneration": generation
        }
        
        if sent_timestamp:
            status_patch["sentTimestamp"] = sent_timestamp.isoformat().replace('+00:00', 'Z')
        
        if error_message is not None:
            status_patch["errorMessage"] = error_message
        
        if error_timestamp:
            status_patch["errorTimestamp"] = error_timestamp.isoformat().replace('+00:00', 'Z')
        
        # Patch the Email CR status subresource
        api.patch_namespaced_custom_object_status(
            group="product.webshop.harikube.info",
            version="v1",
            namespace=namespace,
            plural="emails",
            name=name,
            body={"status": status_patch}
        )
        
    except ApiException as e:
        # Log error but don't fail the function
        print(f"Failed to patch Email status: {e.status} {e.reason}")
        print(f"Response body: {e.body}")
    except Exception as e:
        print(f"Unexpected error patching Email status: {str(e)}")


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
    
    # Create message
    msg = MIMEMultipart('alternative')
    msg['Subject'] = subject
    msg['From'] = f"{from_name} <{from_address}>"
    msg['To'] = to_address
    
    # Add body (support both plain text and HTML)
    if body.strip().startswith('<'):
        # Likely HTML
        msg.attach(MIMEText(body, 'html'))
    else:
        msg.attach(MIMEText(body, 'plain'))
    
    # Send email
    with smtplib.SMTP(smtp_host, smtp_port) as server:
        if smtp_use_tls:
            server.starttls()
        
        if smtp_username and smtp_password:
            server.login(smtp_username, smtp_password)
        
        server.send_message(msg)


def read_secret(key):
    """
    Read secret value from OpenFaaS secret mount path.
    Returns None if secret file doesn't exist.
    """
    secret_path = Path(f'/var/openfaas/secrets/{key}')
    try:
        if secret_path.exists():
            return secret_path.read_text().strip()
    except Exception:
        pass
    return None

