# Simple Web Application - Docker + Kubernetes Learning Project

Full-stack application with React frontend, Go backend, and PostgreSQL database, containerized with Docker and ready for Kubernetes deployment. Sole purpose of this app is to showcase and practice kubernetes and docker skills. The app itself doesn't have much functionality. A working version running on a small scale cluster can be visited on: [dev.wtfthiscantbe.art](https://dev.wtfthiscantbe.art)

## Quick Start

1. **Copy environment variables:**

   ```bash
   cp .env.example .env
   ```

   Edit `.env` and set your own database credentials (don't use the example values!):

   ```bash
   POSTGRES_US**=your_username
   POSTGRES_P******D=set_your_***
   POSTGRES_DB=your_database_name
   ```

2. **Fix database init script permissions:**

   ```bash
   chmod 644 database/init.sql
   ```

   **Important:** This step is required for the database initialization to work.

3. **Start the application:**

   ```bash
   docker-compose up --build
   ```

4. **Access the application:**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8081
   - PostgreSQL: localhost:5432

## Project Structure

```
SimpleWebApplication/
├── backend/                 # Go backend service
│   ├── main.go
│   └── Dockerfile
├── frontend/                # React frontend
│   ├── src/
│   ├── Dockerfile          # Production build (with Nginx)
│   └── Dockerfile.dev      # Development build (hot reload)
├── database/
│   └── init.sql            # Database initialization & seed data
├── docker-compose.yml      # Development orchestration
├── .env                    # Environment variables (NOT committed)
└── .env.example            # Environment template (committed)
```

## Development

### Running in Development Mode

The default `docker-compose.yml` runs in development mode:

- **Frontend:** React dev server with hot reload
- **Backend:** Go binary (rebuild required for changes)
- **Database:** PostgreSQL with persistent data

```bash
# Start all services
docker-compose up

# Start in background
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Making Changes

**Frontend changes:**

- Edit files in `frontend/src/`
- Changes reflect automatically (hot reload)

**Backend changes:**

- Edit `backend/main.go`
- Rebuild: `docker-compose up --build backend`

**Database schema changes:**

- Edit `database/init.sql`
- Reset database: `docker-compose down -v && docker-compose up`

## Database

### Access PostgreSQL

```bash
# Connect to database
docker exec -it postgres-db psql -U <username_from_your_env> -d <tablename_from_your_env>

# Useful psql commands:
\dt                 # List tables
\d users            # Describe table
SELECT * FROM users;
\q                  # Quit
```

### Reset Database

```bash
docker-compose down -v
docker-compose up
```

### Services

**Frontend (React + Dev Server)**

- Port: 3000 (configurable)
- Hot reload enabled
- Proxies `/api/*` requests to backend
- Volume mount: `./frontend` → `/app`

**Backend (Go)**

- Port: 8081 (configurable)
- REST API endpoints:
  - `GET /` - Landing page
  - `GET /hello` - Hello endpoint
  - `GET /bye` - Bye endpoint
  - `GET /impressum` - Copyright info

**Database (PostgreSQL 15)**

- Port: 5432 (configurable)
- Auto-initialized with schema and seed data
- Persistent volume: `postgres-data`

### Networking

All services communicate via Docker network `app-network`:

- Frontend can reach backend at `http://backend:8081`
- Backend can reach database at `postgres:5432`
- Host can reach services via exposed ports

## Kubernetes Production Deployment

### Prerequisites

1. **Install doctl (DigitalOcean CLI):**

   ```bash
   # Linux
   cd ~
   wget https://github.com/digitalocean/doctl/releases/download/v1.94.0/doctl-1.94.0-linux-amd64.tar.gz
   tar xf doctl-1.94.0-linux-amd64.tar.gz
   sudo mv doctl /usr/local/bin
   ```

2. **Initialize doctl with your DigitalOcean account:**
   ```bash
   doctl auth init
   ```
   Enter DigitalOcean API token when prompted.

### One-Time Cluster Setup

**1. Create DOKS Cluster** (via DigitalOcean Dashboard or CLI)

**2. Authenticate kubectl with the cluster:**

```bash
# List your clusters to get the cluster ID or name
doctl kubernetes cluster list

# Save cluster credentials to kubectl config
doctl kubernetes cluster kubeconfig save <cluster-id-or-name>

# Verify connection
kubectl get nodes
```

**3. Install NGINX Ingress Controller:**

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.1/deploy/static/provider/cloud/deploy.yaml

# Wait for LoadBalancer to get external IP (takes 1-2 minutes)
kubectl get svc -n ingress-nginx ingress-nginx-controller --watch
```

**4. Install cert-manager (for automatic TLS certificates):**

```bash
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.3/cert-manager.yaml

# Wait for cert-manager to be ready
kubectl wait --for=condition=ready pod -l app.kubernetes.io/instance=cert-manager -n cert-manager --timeout=300s

# Apply Let's Encrypt issuer (update email in the file first!)
kubectl apply -f k8s/letsencrypt-issuer.yaml
```

**5. Create Kubernetes secrets:**

```bash
# PostgreSQL credentials (possible to do it from your local env file, or create them manually with different values)
kubectl create secret generic postgres-secret --from-env-file=.env

# GitHub Container Registry access
kubectl create secret docker-registry ghcr-secret \
  --docker-server=ghcr.io \
  --docker-username=<your-github-username> \
  --docker-password=<github-personal-access-token> \
  --docker-email=<your-email>
```

**6. Create ConfigMap for database initialization:**

```bash
kubectl create configmap postgres-init-script \
  --from-file=init.sql=database/init.sql
```

### TLS/HTTPS Setup (Optional but Recommended)

**Prerequisites:**

- A domain name (e.g., `yourdomain.com`)
- DNS configured to point to your LoadBalancer IP

**Steps:**

1. **Update the Let's Encrypt issuer email:**
   Edit `k8s/letsencrypt-issuer.yaml` and replace `your-email@example.com` with your actual email

2. **Configure DNS:**

   ```bash
   # Get your LoadBalancer IP
   kubectl get svc -n ingress-nginx ingress-nginx-controller

   # Add an A record in your DNS provider:
   # yourdomain.com → <EXTERNAL-IP>
   ```

3. **Update Ingress with TLS:**
   The Ingress manifest will be updated to include TLS configuration and your domain

4. **Certificate issuance:**
   cert-manager will automatically request and renew certificates from Let's Encrypt

### GitHub Actions Setup

Add these secrets to your GitHub repository (Settings → Secrets and variables → Actions):

1. **DIGITALOCEAN_ACCESS_TOKEN**: Your DigitalOcean API token (with read/write permissions)
2. **CLUSTER_ID**: Your DOKS cluster ID (from `doctl kubernetes cluster list`)
3. **LETSENCRYPT_EMAIL**: Your email address for Let's Encrypt certificate notifications

### Manual Deployment (Optional)

If you want to deploy manually without GitHub Actions:

```bash
# Deploy PostgreSQL
kubectl apply -f k8s/postgres-pvc.yaml
kubectl apply -f k8s/postgres-deployment.yaml
kubectl apply -f k8s/postgres-service.yaml

# Deploy backend
kubectl apply -f k8s/backend-deployment.yaml
kubectl apply -f k8s/backend-service.yaml

# Deploy frontend
kubectl apply -f k8s/frontend-deployment.yaml
kubectl apply -f k8s/frontend-service.yaml

# Deploy Ingress
kubectl apply -f k8s/ingress.yaml

# Check deployment status
kubectl get pods
kubectl get svc
kubectl get ingress
```
