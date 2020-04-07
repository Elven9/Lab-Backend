# /usr/bin/zsh
# $1 -> for container's name
# $2 -> for container's exposed port
# Update Code
git pull

# Stop Container
docker container stop lab-backend
docker container rm lab-backend

# Build Image
docker image rm elven9/lab-backend:latest
yes | docker image prune

docker build -t elven9/lab-backend:latest .
yes | docker image prune

# Run Container
docker run -d --name $1 -p $2:8080 elven9/lab-backend:latest