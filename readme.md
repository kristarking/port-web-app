# Go App — AKS Deployment (Distroless · Multi-Stage · Terraform · TLS Ingress)

[![build-badge](https://img.shields.io/badge/build-passing-brightgreen)](#) ![go-version](https://img.shields.io/badge/go-1.22-blue) ![license](https://img.shields.io/badge/license-MIT-lightgrey)

> Minimal Go web app built with a multi-stage, distroless Docker image, provisioned to **Azure Kubernetes Service (AKS)** via **Terraform**, deployed with **1 replica**, and exposed via a TLS-enabled Ingress.

---

## Table of contents

- [Project overview](#project-overview)  
- [Features](#features)  
- [Prerequisites](#prerequisites)  
- [Repository layout](#repository-layout)  
- [Multi-stage Dockerfile (distroless)](#multi-stage-dockerfile-distroless)  
- [Terraform (provision AKS)](#terraform-provision-aks)  
- [Kubernetes manifests (Deployment, Service, Ingress + TLS)](#kubernetes-manifests-deployment-service-ingress--tls)  
- [Quickstart](#quickstart)  
- [Local development & testing](#local-development--testing)  
- [Troubleshooting & tips](#troubleshooting--tips)  
- [License & author](#license--author)

---

## Project overview

This repository demonstrates:

- Compiling a Go binary in a build stage and packaging only the binary into a **distroless** runtime image (tiny surface/size).  
- Infrastructure-as-code using **Terraform** to create an AKS cluster and necessary Azure resources.  
- Kubernetes manifests for a Deployment (1 replica), Service, and TLS-enabled Ingress to terminate HTTPS.  
- A small, production-minded layout suitable for CI/CD pipelines.

---

## Features

- ✅ Go application (compiled with `CGO_ENABLED=0`)  
- ✅ Multi-stage Docker build (build → distroless runtime)  
- ✅ Minimal final image size and attack surface  
- ✅ AKS deployment with Terraform provisioning  
- ✅ TLS termination at Ingress (Kubernetes Secret or cert-manager)  
- ✅ 1 replica by default (easy to scale)

---

## Prerequisites

- Go (for local dev): `>=1.18` (example uses 1.22)  
- Docker / Buildx  
- Azure CLI + logged in (`az login`)  
- Terraform `>=1.3`  
- `kubectl` configured for target AKS cluster  
- (Optional) `helm` if using cert-manager for TLS

---

## Repository layout

```
.
├── app/
│   └── main.go                # simple Go HTTP server
├── Dockerfile                 # multi-stage + distroless
├── kubernetes/
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── ingress.yaml
│   └── tls-secret.yaml        # optional example if you provide certs manually
└── terraform/
    ├── main.tf
    ├── variables.tf
    └── outputs.tf
```

---

## Multi-stage Dockerfile (distroless)

```dockerfile
# -------------------------
# Build stage
# -------------------------
FROM golang:1.22 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# produce a static binary for linux
 RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./app

# -------------------------
# Runtime stage - distroless
# -------------------------
FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
```

---

## Terraform — provision AKS (example)

**`terraform/main.tf`**

```hcl
provider "azurerm" {
  features = {}
}

resource "azurerm_resource_group" "rg" {
  name     = var.rg_name
  location = var.location
}

resource "azurerm_kubernetes_cluster" "aks" {
  name                = var.aks_name
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  dns_prefix          = var.dns_prefix
  default_node_pool {
    name       = "agentpool"
    node_count = var.node_count
    vm_size    = var.node_vm_size
  }
  identity {
    type = "SystemAssigned"
  }
}

output "kube_config" {
  value     = azurerm_kubernetes_cluster.aks.kube_config_raw
  sensitive = true
}
```

Deploy:

```bash
cd terraform
terraform init
terraform apply
```

---

## Kubernetes manifests

### `kubernetes/deployment.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-app
  labels:
    app: go-app
spec:
  replicas: 1
  selector:
    matchLabels:
      app: go-app
  template:
    metadata:
      labels:
        app: go-app
    spec:
      containers:
        - name: go-app
          image: <registry>/go-app:latest
          ports:
            - containerPort: 8080
```

### `kubernetes/service.yaml`

```yaml
apiVersion: v1
kind: Service
metadata:
  name: go-app-svc
spec:
  selector:
    app: go-app
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: ClusterIP
```

### `kubernetes/ingress.yaml`

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: go-app-ingress
  annotations:
    kubernetes.io/ingress.class: "nginx"
spec:
  tls:
    - hosts:
        - example.yourdomain.com
      secretName: go-app-tls
  rules:
    - host: example.yourdomain.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: go-app-svc
                port:
                  number: 80
```

### TLS secret (optional)

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: go-app-tls
type: kubernetes.io/tls
data:
  tls.crt: <base64 cert>
  tls.key: <base64 key>
```

---

## Quickstart

### 1️⃣ Provision AKS via Terraform

```bash
cd terraform
terraform init
terraform apply
```

### 2️⃣ Build and push the image

```bash
docker build -t <registry>/go-app:latest .
docker push <registry>/go-app:latest
```

### 3️⃣ Deploy to AKS

```bash
kubectl apply -f kubernetes/
```

### 4️⃣ Confirm deployment

```bash
kubectl get deployments,svc,ingress
```

---

## Local development & testing

Run locally:

```bash
cd app
go run main.go
```

Run tests:

```bash
go test ./...
```

Test container locally:

```bash
docker build -t go-app:local .
docker run --rm -p 8080:8080 go-app:local
```

---

## Troubleshooting & tips

- If TLS isn’t working, ensure your domain matches the cert and secret.
- Make sure the Ingress controller is running.
- Use `kubectl describe ingress` for debugging.
- For production, prefer `cert-manager` + Let’s Encrypt.

---

## License & author

MIT License.

**Author:** Chris O

Pull requests are welcome!


