# HariKube: The Cloud-Native-Platform-as-a-Service Reality

Welcome! We are HariKube, and we believe the future of service development hinges on transforming Kubernetes into a true Cloud-Native-Platform-as-a-Service (CNPaaS).

The challenge today is clear: Kubernetes' scalability is bottlenecked by its centralized data store, and development remains fragmented across functions, operators, and APIs.

## The HariKube Solution: A Revolutionary Platform Shift

HariKube introduces a fundamental architectural change that solves both problems simultaneously:

- Hyperscale: We replace the native data store with an innovative Dynamic Data Layer. This agnostic layer allows you to use a user-defined topology with various databases (like MySQL, PgSQL), eliminating Kubernetes' main scalability bottleneck and unlocking true hyperscale.
- Unified Service Design: We transform Kubernetes into the Single Source of Truth for your entire Cloud-Native estate. This means your functions (serverless), microservices (operators), and traditional REST APIs are all unified under the extended Kubernetes API.

## The Strategic Impact

This unified, intelligent operating system for your cloud estate leads to a powerful double-win:

- Developers focus purely on business logic, consuming platform capabilities through the Kubernetes API without managing boilerplate YAML or self-managed databases.
- Infrastructure complexity becomes an API, accelerating innovation and dramatically cutting time-to-market.

The result is clear: We are making the Kubernetes platform disappear, leaving only pure business value behind.

Join us as we demonstrate how HariKube turns this vision into a deployable, functional reality, changing how your teams build and manage cloud-native services forever.

## Infrastructure Setup: Building the CNPaaS

- Creates a local Kubernetes cluster using the Kind tool for the demonstration environment. Bring your own Kubernetes HariKube works seamlessly on all.
- Waits until the control plane node in the new cluster reports that it is ready.
- Installs Cert-Manager to handle TLS certificates within the cluster.
- Applies the necessary Custom Resource Definitions (CRDs) for the Prometheus Operator.
- Waits for the Cert-Manager webhook component to be fully available.
- Creates the dedicated harikube namespace for the platform components.
- Creates a Docker registry secret to allow the platform components to pull images. We have an open source editioin of HariKube, which is limited to one single database, but the business edition we use for this demo is subscription based.
- Deploys the HariKube Operator, which manages the dynamic data layer.
- Deploys the middleware/vcluster workload, which runs the HariKube data fabric. The middleware is totally transparent for the view of Kubernetes, using HariKube doesn't require to modify your applications. This setup uses a virtual Kubernetes instance for serving application data. It helps separate infrastructure and business APIs without stepping on each other's foot.
- Waits for the HariKube Operator deployment to be ready.
- Waits for the vCluster to be ready.
- Creates a TopologyConfig for the Tenant custom resource. Topology config defines the routing policy for HariKube, so all Tenant objects will be stored in this database.
- Creates a TopologyConfig for the User custom resource.
- Creates a TopologyConfig for the Email custom resource.
- Checks HariKube middleware logs, looking for registered databases.
- Connects the local environment to the virtual cluster provisioned by HariKube, which will house the service development workloads. For simplicity we run all application services on the virtual cluster, for more advanced usecases you can use an API only version of the virtual cluster.
- Installs Cert-Manager inside the connected virtual cluster.
- Waits for the Cert-Manager webhook to be ready in the virtual cluster.
- Before we install Kube Watch Trigger, we create RBAC for the operator. Kube Watch Trigger is our open source tool to trigger serverless functions, webhooks based on Kubernetes domain changes.
- Installs the Serverless Kube Watch Trigger to turn Kubernetes events into function invocations.
- Builds and deploys our example the Webshop Operator, which handles the complex, stateful business logic. It uses Serverless function, Operators, and Kubernetes Aggregation API to demonstrate unified service design concept of HariKube.
- Waits for the Webshop Operator to be ready.
- Adds the OpenFaaS Helm chart repository for the serverless layer.
- Fetches the latest updates from the Helm repositories.
- Creates the necessary namespaces for OpenFaaS.
- Installs the OpenFaaS serverless platform in the cluster.
- Waits for the OpenFaaS gateway to be fully ready.
- Pulls the necessary Python 3 function template for the serverless logic.
- Builds the serverless function image. The function is simple, sends emails and updates resource statuses.
- Pushes the function image to the configured registry.
- Logs into the OpenFaaS gateway.
- Deploys the serverless function to the OpenFaaS gateway.
- Patches the serverless function deployment to use the appropriate service account for Kubernetes API access. Unfortunately Community Edition doesn't support service account injection, which is necessary to connect to the Kubernetes API.
- Waits for the patched function deployment to complete its rollout and become ready.
- Creates a HTTPTrigger config, to trigger Email function when an Email object has been created.

Tha platform is ready to be the source of truth and serve millions and millions of records. Each layer is individually scalable and designed to handle huge amount of data, like a webshop. The services are not separated entities on top of Kubernetes any more, they became first class citizens, and they can leverage of all the built-in Kubernetes features.

## HariKube Architecture and User Registration Flow

The HariKube architecture, as depicted in the sequence diagram, fundamentally transforms Kubernetes into a Cloud-Native-Platform-as-a-Service (CNPaaS) by leveraging its core extensibility mechanisms (APIs, Controllers, Webhooks) while replacing its default data store with a dynamic, database-agnostic data layer.

### Core Architecture and Components

- Control Plane: K8S is the central nervous system; HariKube is the dynamic data layer middleware that routes CR persistence to specialized databases; ETCD is used only for native K8S resources like Namespace and RBAC.
- Data Layer: These are specialized, user-defined, database backends (e.g., SQLite, MySQL, PgSQL) managed by HariKube's dynamic data layer for their respective custom resources (Tenant, User, Email). For simplicity we use SQLite for the demo, but you have the freedome to use your favorite database per TopologyConfig.
- Webhooks: Validate and/or default the Custom Resources before they are persisted to the data store.
- REST Layer: An Aggregated API Server that sits in front of the Kubernetes API, providing the external entry point for the user registration requests. Registration is a synchronious process, doesn't fit into Kubernetes' event driven approach. In this layer you have full control on implementing your business needs.
- Operators: Implement the complex business logic (reconciliation loops) by watching for CR changes, creating new resources (like Tenant and User), and handling setup tasks (Namespace, RBAC).
- Serverless: The Event-Driven layer that reacts to a change in a Custom Resource (Email) to execute lightweight business logic (sending an email) without requiring a full operator.

### Core Functionality: User Registration

The goal of the User Registration workflow is to transform a single initial request (RegistrationRequest) into a series of persistent, interdependent resources —specifically, a Tenant, a User, and an Email record— while simultaneously setting up the necessary infrastructure (Namespace, RBAC). This is all orchestrated declaratively through the Kubernetes API. The codebase focuses on easy understanding, not on perfection.

### User Registration Flow

The flow is initiated by the user and then orchestrated entirely by the HariKube platform components through a series of declarative actions and automated reconciliation loops:

- Request Initiation: The User sends a Request across Kubernetes API which is routed through the Aggregated API Server to create a RegistrationRequest Custom Resource in Kubernetes.
- API Validation and Persistence (RegistrationRequest):
 - The RegistrationRequest passes through a Webhook for Defaulting and Validation.
 - HariKube intercepts the persistence request and directs the RegistrationRequest CR to be stored in the default ETCD data store.
- Operator Reconciliation (RegistrationRequestController):
 - The creation of the RegistrationRequest CR triggers the RegistrationRequestController.
 - The Controller's core job is to create the Tenant, User and Email CRs, waiting for the TenantController to create the Namespace first.
 - Data Persistence: HariKube routes the Tenant CR to the dedicated TenantDatabase, the User CR to the UserDatabase, and the Email CR to the EmailDatabase.
 - Finally, deletes the original RegistrationRequest.
- Dependent Operators and Infrastructure Setup:
 - The creation of the Tenant CR triggers the TenantController, which creates a Namespace (persisted in ETCD).
 - The creation of the User CR triggers the UserController, which sets up RBAC (Role-Based Access Control) and the user's Password Secret (both persisted in ETCD).
- Serverless Trigger and Completion (Email):
 - The creation of the Email CR triggers the EmailTrigger.
 - The trigger invokes the EmailFunction, which sends the welcome email to the User.
 - The function updates the Email CR's status, which is persisted by HariKube to the EmailDatabase.
- Final Response: Once the RegistrationRequest is deleted, the API sends a 201 [OK] response back to the User, signifying that the registration process has been initiated and successfully orchestrated.

## Back to the demo

We are creating a RegistrationRequest via kubectl. Yes you see well, all the services are Kubernetes citizens, so you can manage your users with standard Kubernetes tools.

## The CNPaaS Future is Now

We've reached the end of our HariKube demonstration, and the core takeaway is clear: the future of Cloud-Native service development is here, and it is a unified, intelligent platform built on the evolution of Kubernetes.

HariKube fundamentally delivers on the promise of a true Cloud-Native-Platform-as-a-Service (CNPaaS) by solving the two most critical architectural challenges facing modern enterprises:

- Hyperscale and Data Freedom: We demonstrated how the Dynamic Data Layer effectively replaces the single-instance ETCD bottleneck, allowing you to route CRD persistence to specialized, high-performance databases (like the TenantDatabase, UserDatabase, and EmailDatabase). This unlocks true hyperscale for every Kubernetes deployment.
- Unified Service Design: We showed how functions (Serverless), operators (Microservices), and aggregated APIs (REST) can be seamlessly integrated into a single, cohesive architecture. The entire User Registration Flow was orchestrated solely by the Kubernetes API and its controllers, using CRDs as the single source of truth.

By tackling data scalability and design fragmentation, HariKube empowers your teams to achieve ultimate separation of concerns:

- Developers focus 100% on business logic, consuming clean APIs and CRDs.
- Infrastructure Teams manage platform capabilities, where complexity is an API, not a configuration hurdle.

This architectural transformation means faster innovation, reduced costs, and a significant acceleration of your time-to-market. HariKube makes the platform disappear, leaving only pure business value behind.

Thank you for joining us as we move past the traditional container scheduler and welcome the reality of the Cloud-Native-Platform-as-a-Service.