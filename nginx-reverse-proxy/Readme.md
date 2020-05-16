# Nginx Server

## Start The Server

```zsh
# Start Nginx Main Server With Configuration File
docker run -d --name nginx-reverse-proxyserver \
           -p 8080:8080 \
           --network test-network \
           --mount type=bind,source="$(pwd)/nginx-reverse-proxy/nginx.conf",target=/etc/nginx/nginx.conf \
           nginx
```

## Reload Config File

```zsh
docker exec nginx-reverse-proxyserver nginx -s reload
```