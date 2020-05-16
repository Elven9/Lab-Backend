# Grafana

## Socker Image Run Configuration

```shell
# Create Container
docker run -d --name grafana \
           -v grafana-storage:/var/lib/grafana \
           -e "GF_SERVER_DOMAIN=localhost" \
           -e "GF_SERVER_ROOT_URL=https://localhost:3000/grafana/" \
           -e "GF_SERVER_SERVE_FROM_SUB_PATH=true" \
           grafana/grafana

# Attach to a Network
docker network connect test-network grafana
```

## Query Backup

PromQL

```
# Jobs Memory Usage
sum by(group)(label_replace(container_memory_usage_bytes{namespace="default"}, "group", "$1", "pod", "(.*)-(.*)-.*"))

# Jobs CPU Usage
sum by(group)(label_replace(rate(container_cpu_usage_seconds_total{namespace="default"}[$__range]), "group", "$1", "pod", "(.*)-(.*)-.*"))

# Job CPU Usage
sum(rate(container_cpu_usage_seconds_total{namespace="default", pod=~"^[[jobName]]-(worker|ps)-\\d"}[$__range]))

# Job Memory Usage
sum(container_memory_usage_bytes{namespace="default", pod=~"^[[jobName]]-(worker|ps)-\\d"})

# Job Worker Count
sum(kube_pod_container_status_ready{namespace=~"default", pod=~"^[[jobName]]-worker-\\d"})

# Parameter Server Count
sum(kube_pod_container_status_ready{namespace=~"default", pod=~"^[[jobName]]-ps-\\d"})
```