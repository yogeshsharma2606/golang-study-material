# Lab 16: Docker and Kubernetes

Containerize a small Go HTTP service with a multi-stage Dockerfile, run it with Compose, and apply a sample Kubernetes manifest.

## Objectives

- Build a minimal production-style image (static binary, non-root user).
- Orchestrate locally with Docker Compose including a health check.
- Understand a basic Deployment + Service for Kubernetes.

## Setup

Run locally without Docker:

```bash
cd labs/16-docker/app
go run .
```

Build and run with Docker:

```bash
cd labs/16-docker
docker compose up --build
```

Visit `http://localhost:8080/health`.

## Exercises

1. Inspect image size: `docker images` before and after switching base images.
2. Change `APP_VERSION` in `docker-compose.yml` and confirm `/` reflects it.
3. Apply to a local cluster (e.g. minikube/kind):
   ```bash
   docker build -t docker-lab-app:latest .
   kubectl apply -f k8s/deployment.yaml
   kubectl port-forward svc/docker-lab-app 8080:80
   ```
4. Tighten the Dockerfile: add `.dockerignore` excluding `k8s/` and README.

## What to take away

- Multi-stage builds keep runtime images small and secure.
- Health checks align with Kubernetes probes.
- `CGO_ENABLED=0` produces portable static binaries for Alpine/scratch.

## Cleanup

```bash
docker compose down
kubectl delete -f k8s/deployment.yaml
```

## Related Modules

- Deployment and DevOps modules.
- Lab 09 HTTP patterns (health endpoint).
