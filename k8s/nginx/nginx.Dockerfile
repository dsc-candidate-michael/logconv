FROM nginx:1.7.9

# Just for debugging purposes
RUN apt-get update && apt-get install -y vim 

ADD k8s/nginx/nginx.conf /etc/nginx/nginx.conf