.PHONY: build

build:
	go build -o $(BINARY) .

container:
	@docker build -t my-secure-registry.org/liquibase-lock-release:v1

deploy:
	@echo "==== Deploying CronJob to K8s, make sure secret is set in the cluster ===="
	@kubectl apply -f ./k8s