# Local build and tests targets

all: build

build:
	go build ./cmd...

clean: 
	rm -f logconv	
	rm -f test-artifacts/*

test: unit-test

test-container: 
	docker run logconv:v1

unit-test:
	go test -v ./...

# Docker and Kubernetes build and deploy targets

image: 
	docker build -t logconv:v1 .
	docker build -t logconv-nginx:v1 . -f k8s/nginx/nginx.Dockerfile

configure: 
	kubectl apply -f k8s/config.yaml

deploy:
	kubectl apply -f k8s/ingress.yaml
	kubectl apply -f k8s/nginx/deployment.yaml
	kubectl expose deployment nginx --type=LoadBalancer --name=nginx

destroy:
	kubectl delete cm logconv-config || true
	kubectl delete deployment nginx || true
	kubectl delete svc nginx || true

create-cluster:
	minikube start
	minikube addons enable ingress
	
destroy-cluster:
	minikube delete