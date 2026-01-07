# Kubernetes Manifests

## Makefile Helper Commands

The top-level `Makefile` provides convenient commands for common tasks:

### Build Images
Build Docker images in minikube's Docker context:
```bash
make build-images
```

This builds both backend and frontend images with the correct tags (`v1.0.0`) in minikube's Docker daemon.

### Create Secrets
Create Kubernetes secrets from your `.env` file:
```bash
make deploy-secrets
```

This is idempotent. It will create or update the secret.

### Seed Database
Create the ConfigMap from `database/init.sql`:
```bash
make deploy-configmap
```
The ConfigMap is generated from `database/init.sql`, which is also used for local development etc. with docker-compose.

### Clean Up Resources
Delete all Kubernetes resources:
```bash
make clean
```

---

## Manual Setup Instructions

### 1. Build Images in Minikube

```bash
eval $(minikube docker-env)
docker build -t simplewebapp-backend:v1.0.0 -f backend/Dockerfile backend/
docker build -t simplewebapp-frontend:v1.0.0 -f frontend/Dockerfile frontend/
```

### 2. Create Secrets

```bash
kubectl create secret generic postgres-secret --from-env-file=.env
```

### 3. Create ConfigMap for Database Initialization

```bash
kubectl create configmap postgres-init-script --from-file=init.sql=database/init.sql
```

### 4. Deploy Application

Deploy in order:

```bash
# PostgreSQL
kubectl apply -f k8s/postgres-pvc.yaml
kubectl apply -f k8s/postgres-deployment.yaml
kubectl apply -f k8s/postgres-service.yaml

# Backend
kubectl apply -f k8s/backend-deployment.yaml
kubectl apply -f k8s/backend-service.yaml

# Frontend
kubectl apply -f k8s/frontend-deployment.yaml
kubectl apply -f k8s/frontend-service.yaml
```

### 5. Enable Ingress (One-time Setup)

Enable the NGINX Ingress Controller in minikube:

```bash
minikube addons enable ingress
```

Verify the Ingress controller is running:

```bash
kubectl get pods -n ingress-nginx
```

Wait until the `ingress-nginx-controller` pod shows `Running` status.

### 6. Deploy Ingress

```bash
kubectl apply -f k8s/ingress.yaml
```

Check Ingress status:

```bash
kubectl get ingress
```

### 7. Configure Local DNS

Get your minikube IP:

```bash
minikube ip
```

Add the hostname to your `/etc/hosts` file (requires sudo):

```bash
sudo nano /etc/hosts
```

Add this line (replace `192.168.49.2` with your actual minikube IP):

```
192.168.49.2  simplewebapp.local
```

Save and exit (Ctrl+O, Enter, Ctrl+X).

### 8. Verify Deployment

```bash
kubectl get pods
kubectl get services
kubectl get pvc
kubectl get ingress
```

### 9. Access the Application

Open your browser and navigate to:

- **Frontend:** http://simplewebapp.local/
- **API Test:** http://simplewebapp.local/api/hello

The application should load with the Scandinavian fjord background and three interactive buttons.
