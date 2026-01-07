.PHONY: help build-images deploy-secrets deploy-configmap deploy-postgres deploy-backend deploy-frontend deploy-all clean

build-images:
	@echo "Building images in minikube Docker context..."
	eval $$(minikube docker-env) && \
	docker build -t simplewebapp-backend:v1.0.0 -f backend/Dockerfile backend/ && \
	docker build -t simplewebapp-frontend:v1.0.0 -f frontend/Dockerfile frontend/
	@echo "Images built successfully!"

deploy-secrets:
	@echo "Creating secrets from .env file..."
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found!"; \
		exit 1; \
	fi
	kubectl create secret generic postgres-secret --from-env-file=.env --dry-run=client -o yaml | kubectl apply -f -
	@echo "Secrets created successfully!"

deploy-configmap:
	@echo "Creating ConfigMap from database/init.sql..."
	kubectl create configmap postgres-init-script --from-file=init.sql=database/init.sql --dry-run=client -o yaml | kubectl apply -f -
	@echo "ConfigMap created successfully!"

clean:
	@echo "Deleting all resources..."
	kubectl delete -f k8s/frontend-deployment.yaml --ignore-not-found=true
	kubectl delete -f k8s/frontend-service.yaml --ignore-not-found=true
	kubectl delete -f k8s/backend-deployment.yaml --ignore-not-found=true
	kubectl delete -f k8s/backend-service.yaml --ignore-not-found=true
	kubectl delete -f k8s/postgres-deployment.yaml --ignore-not-found=true
	kubectl delete -f k8s/postgres-service.yaml --ignore-not-found=true
	kubectl delete -f k8s/postgres-pvc.yaml --ignore-not-found=true
	kubectl delete secret postgres-secret --ignore-not-found=true
	kubectl delete configmap postgres-init-script --ignore-not-found=true
	@echo "All resources deleted!"
