.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: all
all: build_services deploy_services deploy_workflow ## build & deploy all service

.PHONY: build_services
build_services: ## build all service
	gcloud builds submit --tag gcr.io/$(PROJECT_ID)/order-service services/order
	gcloud builds submit --tag gcr.io/$(PROJECT_ID)/stock-service services/stock
	gcloud builds submit --tag gcr.io/$(PROJECT_ID)/payment-service services/payment

.PHONY: deploy_services
deploy_services: ## deploy all service
	gcloud run deploy order-service \
	--image gcr.io/$(PROJECT_ID)/order-service \
	--platform=managed --region=us-central1 \
	--set-env-vars "PROJECT_ID=$(PROJECT_ID)" \
	--no-allow-unauthenticated
	gcloud run deploy stock-service \
	--image gcr.io/$(PROJECT_ID)/stock-service \
	--platform=managed --region=us-central1 \
	--no-allow-unauthenticated
	gcloud run deploy payment-service \
	--image gcr.io/$(PROJECT_ID)/payment-service \
	--platform=managed --region=us-central1 \
	--no-allow-unauthenticated

.PHONY: deploy_workflow
deploy_workflow: ## deploy workflow
	gcloud workflows deploy workflow-test \
	--source=workflow.yml \
	--service-account=$(SERVICE_ACCOUNT)
