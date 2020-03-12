# /usr/bin/zsh
# Update Code
git pull

# Build Image
docker image rm elven9/lab-backend:latest
docker image prune

docker build -t elven9/lab-backend:latest .
docker image prune

# Run Container
docker run -d --name lab-backend -p 9000:8080 --rm elven9/lab-backend:latest