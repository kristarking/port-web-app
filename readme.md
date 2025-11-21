Go Application – AKS Deployment (Distroless, Multi-Stage, Terraform, TLS Ingress)
📌 Project Overview
This project is a production-grade Go application packaged using a multi-stage Docker build and deployed to Azure Kubernetes Service (AKS).
 The goal of this setup is to ensure:
Minimal container image size using distroless base image


Secure and scalable deployment on AKS


Automated infrastructure provisioning using Terraform


Encrypted communication via TLS-enabled Ingress



🚀 Features
🟦 Go Application
Built in Go with performance and simplicity in mind.


Compiled to a minimal binary with no unnecessary runtime dependencies.


🐳 Multi-Stage Docker Build
The app uses a multi-stage Dockerfile:


Stage 1: Builds the Go binary.


Stage 2: Runs it in a distroless image, ensuring:


No package manager


No shell


Very small attack surface


Ultra-lightweight image


☸️ AKS Deployment
Application deployed to Azure Kubernetes Service (AKS).


Deployment configured with:


1 replica


Horizontal scalability support


Kubernetes Deployment, Service, and Ingress files.


🔐 TLS-Enabled Ingress
Configured HTTPS using TLS termination at the Ingress level.


Certificate managed through Kubernetes secrets.


📦 Infrastructure as Code (IaC)
The AKS cluster and required resources were provisioned end-to-end using Terraform:


Resource Group


AKS Cluster


Node Pool


Container Registry


Networking and Ingress



📁 Project Structure
├── app/
│   ├── main.go
├── Dockerfile
├── kubernetes/
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── ingress.yaml
│   └── tls-secret.yaml
└── terraform/
    ├── main.tf
    ├── variables.tf
    ├── outputs.tf


🧱 Multi-Stage Dockerfile
# Build Stage
FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app

# Runtime Stage
FROM gcr.io/distroless/base-debian12
COPY --from=builder /app/app /app
ENTRYPOINT ["/app"]


☁️ Deployment Steps
1️⃣ Provision AKS with Terraform
cd terraform
terraform init
terraform apply

2️⃣ Build & Push Docker Image
docker build -t <registry>/go-app:latest .
docker push <registry>/go-app:latest

3️⃣ Deploy to AKS
kubectl apply -f kubernetes/

4️⃣ Verify TLS Ingress
kubectl get ingress


🌐 Accessing the Application
Once deployed, access the application via the issued HTTPS domain secured by TLS.

🧪 Local Development
Run the app locally:
go run main.go

Run tests:
go test ./...


📎 Technologies Used
Category
Tools
Language
Go
Containerization
Docker (multi-stage, distroless)
Cloud Hosting
Azure Kubernetes Service (AKS)
Infrastructure Automation
Terraform
Security
TLS-enabled Kubernetes Ingress


🙌 Author
Chris O
Feel free to open issues or submit pull requests!

If you'd like, I can also:
✔ Generate Terraform or Kubernetes sample files
 ✔ Add shields/badges to the README
 ✔ Make a more advanced version with diagrams (Mermaid)
Just let me know!

