# Local build and tests targets
all: build

build:
	go build ./cmd...

clean: 
	rm -f logconv	
	rm -f test-artifacts/*

test: unit-test

unit-test:
	go test -v ./...

# Docker and Kubernetes build and deploy targets
image: 
	docker build -t logconv:v1 .
	docker build -t logconv-nginx:v1 . -f k8s/nginx/nginx.Dockerfile

configure: 
	kubectl create configmap logconv-config --from-file=k8s/logconv.config

deploy:
	kubectl apply -f k8s/echoserver/deployment.yaml
	kubectl apply -f k8s/echoserver/service.yaml
	kubectl apply -f k8s/nginx/deployment.yaml
	# kubectl apply -f k8s/nginx/service.yaml
	kubectl expose deployment nginx --type=NodePort --name=nginx

destroy:
	kubectl delete cm logconv-config || true
	kubectl delete deployment echoserver || true
	kubectl delete deployment nginx || true
	kubectl delete svc echoserver || true
	kubectl delete svc nginx || true

create-cluster:
	minikube start
	# minikube addons enable ingress
	eval $(minikube docker-env)

destroy-cluster:
	minikube delete