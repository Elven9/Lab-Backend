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

## Requirement

### Static (no record)

- 系統狀態: 系統有多少資源、剩下多少、每個node剩下多少、有沒有什麼alarm等等
- scheduling: 
(overview) waiting job num, running job num, finish job num, jobs average waiting time, jobs average completion time
(per job) job waiting time, running time, completion time, and each of timestamp
- locality aware的部份主要就是看ps/worker的位置在哪裡，分散率：nodes per job
(overview) 全部的job的平均分散率、以node的角度來看上面有哪些jobs worker
(per job) 每個job的分散率、以job的角度來看worker在哪些node
- scaling主要就是要看worker數量的變化(原本設定vs目前有多少)，還有目前的resource utilization
(overview) 每個node的累計資源長條圖
(per job) 目前的scale和initial scale的差距