# /usr/bin/zsh
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
docker run -d --name lab-backend -p 9000:8080 elven9/lab-backend:latest