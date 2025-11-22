# DevOps Portfolio - Azure Kubernetes Deployment

A complete DevOps project demonstrating Infrastructure as Code, containerization, and automated deployment to Azure Kubernetes Service.

---

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Project Structure](#project-structure)
- [Setup Instructions](#setup-instructions)
- [Deployment](#deployment)
- [Accessing the Application](#accessing-the-application)
- [Known Issues & Fixes](#known-issues--fixes)
- [Cleanup](#cleanup)
- [Technologies Used](#technologies-used)

---

## Overview

This project deploys a Go-based portfolio web application to Azure Kubernetes Service (AKS) using:
- **Terraform** for infrastructure provisioning
- **Docker** for containerization with distroless images
- **GitHub Actions** for CI/CD automation
- **Azure Container Registry** for image storage
- **Kubernetes** for orchestration (3 replicas with LoadBalancer)

---

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    GitHub Actions                        │
│  (Build Docker Image → Push to ACR → Deploy to AKS)    │
└──────────────────┬──────────────────────────────────────┘
                   │
                   ▼
         ┌─────────────────────┐
         │  Azure Container    │
         │    Registry (ACR)   │
         │   kingacr.azurecr.io│
         └──────────┬──────────┘
                    │
                    ▼
         ┌─────────────────────┐
         │   Azure Kubernetes  │
         │   Service (AKS)     │
         │  - 3 Replicas       │
         │  - LoadBalancer     │
         └─────────────────────┘
```

---

## Prerequisites

Before you begin, ensure you have:

- [ ] Azure subscription (with contributor access)
- [ ] Azure CLI installed and logged in
- [ ] Terraform v4.3.0+ installed
- [ ] Docker installed
- [ ] kubectl installed
- [ ] Git installed
- [ ] GitHub account

---

## Project Structure

```
devops-portfolio/
│
├── .github/
│   └── workflows/
│       └── ci-cd-pipeline.yml     # GitHub Actions CI/CD workflow
│
├── docker/
│   └── dockerfile                 # Multi-stage Docker build
│
├── k8s/
│   ├── deployment.yml             # Kubernetes deployment (3 replicas)
│   └── service.yml                # LoadBalancer service
│
├── src/
│   └── main.go                    # Go web application
│
├── terraform/
│   ├── main.tf                    # Infrastructure definitions
│   ├── variables.tf               # Terraform variables
│   └── .gitignore                 # Terraform ignore rules
│
└── README.md                      # This file
```

---

## Setup Instructions

### Step 1: Clone Repository

```bash
git clone https://github.com/yourusername/devops-portfolio.git
cd devops-portfolio
```

### Step 2: Set Up Azure Backend for Terraform

Create storage account for Terraform state:

```bash
# Create resource group
az group create --name KingRG --location "South Africa North"

# Create storage account
az storage account create \
  --name kingst \
  --resource-group KingRG \
  --location "South Africa North" \
  --sku Standard_LRS

# Create container
az storage container create \
  --name tfstate \
  --account-name kingst
```

### Step 3: Configure GitHub Secrets

Create Azure service principal:

```bash
az ad sp create-for-rbac \
  --name "github-actions-sp" \
  --role contributor \
  --scopes /subscriptions/70faaaeb-e126-4ca7-95be-9e614709f37b \
  --sdk-auth
```

Add the JSON output as `AZURE_CREDENTIALS` in GitHub repository secrets:
1. Go to repository **Settings** → **Secrets and variables** → **Actions**
2. Click **New repository secret**
3. Name: `AZURE_CREDENTIALS`
4. Value: Paste the JSON output

### Step 4: Provision Infrastructure

```bash
cd terraform

# Initialize Terraform
terraform init

# Review the plan
terraform plan

# Apply the configuration
terraform apply -auto-approve
```

This creates:
- Resource Group: `KingRG` in South Africa North
- AKS Cluster: `DevOpsCluster` (2 x Standard_DS2_v2 nodes)
- Container Registry: `kingacr`

### Step 5: Attach ACR to AKS

**IMPORTANT**: Your AKS cluster needs permission to pull images from ACR:

```bash
# Get ACR ID
ACR_ID=$(az acr show --name kingacr --resource-group KingRG --query id --output tsv)

# Attach ACR to AKS
az aks update \
  --name DevOpsCluster \
  --resource-group KingRG \
  --attach-acr $ACR_ID
```

---

## Deployment

### Manual Deployment (Optional)

If you want to deploy manually before setting up CI/CD:

```bash
# Build and push Docker image
az acr login --name kingacr
docker build -t kingacr.azurecr.io/devops-portfolio:latest -f docker/dockerfile .
docker push kingacr.azurecr.io/devops-portfolio:latest

# Configure kubectl
az aks get-credentials --resource-group KingRG --name DevOpsCluster --admin

# Deploy to Kubernetes
kubectl apply -f k8s/deployment.yml
kubectl apply -f k8s/service.yml
```

### Automated Deployment (CI/CD)

Push code to `main` branch to trigger GitHub Actions:

```bash
git add .
git commit -m "Deploy to AKS"
git push origin main
```

The pipeline will automatically:
1. ✅ Check out code
2. ✅ Authenticate with Azure
3. ✅ Build Docker image
4. ✅ Push to Azure Container Registry
5. ✅ Deploy to AKS cluster

---

## Accessing the Application

### Get LoadBalancer External IP

```bash
kubectl get service app-service

# Wait for EXTERNAL-IP (may take 2-3 minutes)
kubectl get service app-service --watch
```

### Access the Application

Open browser and navigate to:
```
http://<EXTERNAL-IP>
```

### Application Routes

- `/` - Home page
- `/projects` - Projects showcase
- `/contact` - Contact information

---

## Known Issues & Fixes

### Issue 1: Dockerfile COPY Path

**Current issue in `docker/dockerfile`:**
```dockerfile
COPY ../src/ .  # This may fail depending on build context
```

**Fix:** Update to:
```dockerfile
COPY src/ .
```

Or build from root with proper context:
```bash
docker build -t kingacr.azurecr.io/devops-portfolio:latest -f docker/dockerfile .
```

### Issue 2: ACR Authentication

If pods fail to pull images with `ImagePullBackOff` error:

```bash
# Check pod status
kubectl get pods
kubectl describe pod <pod-name>

# Fix by attaching ACR to AKS
az aks update --name DevOpsCluster --resource-group KingRG --attach-acr kingacr
```

### Issue 3: Terraform Resource Group Conflict

The `main.tf` both creates and references `KingRG`. Ensure the resource group is created FIRST by Terraform, or if it exists, import it:

```bash
terraform import azurerm_resource_group.KingRG /subscriptions/70faaaeb-e126-4ca7-95be-9e614709f37b/resourceGroups/KingRG
```

---

## Monitoring & Debugging

### Check Deployment Status

```bash
# View all resources
kubectl get all

# Check pod logs
kubectl logs -l app=devops-portfolio

# Describe deployment
kubectl describe deployment app-deployment

# Check service
kubectl describe service app-service
```

### Scale Deployment

```bash
# Scale to 5 replicas
kubectl scale deployment app-deployment --replicas=5

# Verify
kubectl get pods
```

---

## Cleanup

### Delete Kubernetes Resources

```bash
kubectl delete -f k8s/deployment.yml
kubectl delete -f k8s/service.yml
```

### Destroy Infrastructure

```bash
cd terraform
terraform destroy -auto-approve
```

### Delete GitHub Actions Artifacts (Optional)

Go to repository **Actions** tab and manually delete old workflow runs.

---

## Technologies Used

| Technology | Version | Purpose |
|------------|---------|---------|
| Go | 1.18 | Application backend |
| Docker | Latest | Containerization |
| Kubernetes | Latest | Orchestration |
| Terraform | ~4.3.0 | Infrastructure as Code |
| Azure AKS | Latest | Managed Kubernetes |
| Azure ACR | Basic SKU | Container registry |
| GitHub Actions | Latest | CI/CD automation |
| Distroless | Debian 10 | Minimal container image |

---

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/improvement`)
3. Commit changes (`git commit -am 'Add new feature'`)
4. Push to branch (`git push origin feature/improvement`)
5. Open a Pull Request

---
## Contact

- **Email**: ojedayochristopher@gmail.com
- **LinkedIn**: https://linkedin.com/in/christopherojedayo

**⭐ Star this repository if you found it helpful!**
