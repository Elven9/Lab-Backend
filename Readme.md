# Elven9/Lab-Backend

## Installation & Run

Install Server's Container And Run it.

```zsh
# Pull Image From Docker Hub
docker pull elven9/lab-backend:latest

# Create Container
docker run -d --name lab-backend -p 9000:8080 --rm elven9/lab-backend:latest
```

Or you can build the image yourself on your computer:

```zsh
# Upgrade Script
zsh upgrade-script.sh

# Run The Same Command Mentioned Above
```