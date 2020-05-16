#/bin/zsh

# Set Up Environment Variable for Docker Compose to Run
export DOCKER_COMPOSE_PWD=$(pwd)
export FRONT_END_PORT=3000

# Docker Compose
docker-compose up