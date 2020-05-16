# Elven9/Lab-Backend

## Deploy a Whole Server Stack

```zsh
zsh setup-server.sh
```

## Installation & Run

Install Server's Container And Run it.

```zsh
# Pull Image From Docker Hub
docker pull elven9/lab-backend:latest

# Create Container
docker run -d --name api-server --mount type=bind,source=/etc/kubernetes/admin.conf,target=/root/.kube/config elven9/lab-backend:latest

# Test Version
docker run -d --name api-server elven9/lab-backend:latest -escapeCheck=true
```

Or you can build the image yourself on your computer:

```zsh
# Upgrade Script
zsh upgrade-script.sh

# Run The Same Command Mentioned One Section Above
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

### Project Process

- 系統狀態
  - [x] 固定資源有多少
  - [x] 資源目前剩多少
  - [ ] alarm
- Scheduling
  - [x] Waiting Job Number
  - [x] Running Job Number
  - [x] Finish Job Number
  - [x] Jobs Average Waiting Time
  - [x] Jobs Average Complete Time
- Locality
  - [x] ps/worker position on each nodes
  - [ ] 數據：分散率
- Scaling
  - [x] Resource Utilization ( 每個 Node 上 Target 值總和 )
- Job List and Single Job Page
  - [ ] job waiting time, running time, completion time, and each of timestamp
  - [ ] 每個job的分散率、以job的角度來看worker在哪些node
  - [ ] 目前的scale和initial scale的差距

### Project Detail

數據分散率算法：

```
node上有該job的worker / 一個job最少可以用幾個node就能跑
```

### Next Target

[From metrics to insight - prometheus](https://prometheus.io/)

### K8S Server Start / Shutdown

```shell
# 關掉三台機器上的 k8s
sudo kubeadm reset

# 在 master 上打（目前 monitor-1 是 master）
sudo kubeadm init --apiserver-advertise-address=10.8.36.221 --pod-network-cidr=10.244.0.0/16
# --authentication-token-webhook=true --authentication-token-webhook=true

# 附這 admin config
sudo cp /etc/kubernetes/admin.conf ./.kube/config

# 複製檔案在 /tmp/kube-flannel.yml 到家裡資料夾
cp /tmp/kube-flannel.yml .

# Apply
kubectl create -f kube-flannel.yml

# 觀察
kubectl -n kube-system get pod -w

# 接下來要把其他的 Node 加到 k8s cluster 中
# 安裝 Dragon
```