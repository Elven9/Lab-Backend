#/bin/zsh

# Pull The Latest Image
docker pull elven9/lab-backend:latest
docker pull elven9/lab-monitor:latest
docker pull grafana/grafana:latest
docker pull nginx:latest
docker pull bitnami/kubectl

# Set Up Environment Variable for Docker Compose to Run
export DOCKER_COMPOSE_PWD=$(pwd)
export FRONT_END_PORT=3000

# Docker Compose
docker-compose up -d