# CronX 分布式任务系统

## 项目介绍

CronX 是一个简洁优雅的分布式任务系统，里面有一个 jobWatcher 示例业务实现。

jobWatcher 实现了定时同步 MySQL 和 Kubernetes 之间的任务状态。

job 表中会定义启动 Kubernetes Job 所运行的镜像，镜像中可以实现任何业务逻辑。

CronX 启动后，jobWatcher 主要干了两件事：

1. 找到 MySQL 表中状态为 `Normal` 的 job 记录，在 Kubernetes 集群中名为 `cronx` 的 Namespace 下启动并运行对应的 Job

2. 通过 Kubernetes List & Watch 机制实时同步在 Kubernetes 中运行的 Job 状态到 MySQL 表对应的 job 记录中

## 快速开始

这里分别演示如何在本地运行 CronX 项目以及将其部署到 Kubernetes 集群中。

**注意**：此示例均以搭载 Apple 芯片的 Mac 系统为例进行演示，如果你是 Windows 系统且执行过程中遇到任何问题，欢迎提 [issues](https://github.com/jianghushinian/go-programming/issues)
一起交流讨论。

### 本地调试

#### 1. clone 项目到本地

```shell
# clone 项目
$ git clone git@github.com:jianghushinian/go-programming.git
# 进入项目根目录
$ cd go-programming/ch8/cronx
```

接下来的步骤如无特殊说明，一切操作都是在 `cronx/` 目录下完成。

#### 2. 安装依赖

通过 docker-compose 安装部署 CronX 依赖的 MySQL 和 Redis。

```shell
$ docker compose -p cronx -f deployments/docker-compose.yaml up -d
[+] Running 5/5
 ✔ Network cronx_cronx              Created                                                                                                        0.0s 
 ✔ Volume cronx_cronx_redis_data    Created                                                                                                        0.0s 
 ✔ Volume cronx_cronx_mariadb_data  Created                                                                                                        0.0s 
 ✔ Container cronx-mariadb          Started                                                                                                        0.3s 
 ✔ Container cronx-redis            Started                 
```

> 注：这里使用 MariaDB 来替代 MySQL，你也可以使用现有的数据库。

#### 3. 准备 SQL

这里主要是建库、建表、创建测试数据。

```shell
$ docker exec -i cronx-mariadb mariadb -u root -pcronx < ./assets/schema.sql
```

`assets/schema.sql` 中有如下两条测试数据：

```sql
INSERT
IGNORE INTO `job` (`id`, `name`, `namespace`, `info`, `status`) VALUES (1, 'demo-job-1', 'cronx', '{"image":"alpine","command":["sleep"],"args":["60"]}', 'Normal');
INSERT
IGNORE INTO `job` (`id`, `name`, `namespace`, `info`, `status`) VALUES (2, 'demo-job-2', 'cronx', '{"image":"busybox","command":["echo"],"args":["Hello Cronx!"]}', 'Normal');
```

各字段含义如下：

- `name`：任务名称，会作为运行在 Kubernetes 集群中的 Job 名称。
- `namespace`：Job 所在的 Kubernetes 命名空间，CronX 项目将所有 Job 统一运行在名为 `cronx` 的 Namespace 下。
- `info`：Kubernetes Job 所需的 Spec 信息，包含镜像、启动命令以及启动参数。通过指定不同镜像，可以实现任意业务逻辑，这极大的扩展了
  jobWatcher 的业务灵活性。
- `status`：任务状态，`Normal` 状态的任务会被同步到 Kubernetes 中执行，启动后状态变为 Pending，接着实时监听 Kubernetes 中对应的
  Job 状态并同步回 MySQL job 表中，直到任务成功完成或执行失败。

#### 4. 构建可执行文件

```shell
# 构建
$ make build
Building cronx for darwin/arm64...
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 \
                go build -ldflags "-X cronx/pkg/version.GitVersion=4d2a173-dirty -X cronx/pkg/version.GitCommit=4d2a173d53e4e4b0ce9f862411b620c5b7ddebfa -X cronx/pkg/version.GitTreeState=dirty -X cronx/pkg/version.BuildDate=2026-02-28T15:04:05Z" -o bin/cronx ./cmd/main.go

# 查看程序版本
$ ./bin/cronx --version
{
  "gitVersion": "4d2a173-dirty",
  "gitCommit": "4d2a173d53e4e4b0ce9f862411b620c5b7ddebfa",
  "gitTreeState": "dirty",
  "buildDate": "2026-02-28T15:04:05Z",
  "goVersion": "go1.25.0",
  "compiler": "gc",
  "platform": "darwin/arm64"
}
```

#### 5. 修改配置

示例配置放在 `configs/cronx.yaml` 中，内容如下：

```yaml
kubeconfig: /Users/jianghushinian/.kube/config

mysql:
  host: 127.0.0.1:3306
  username: root
  password: cronx
  database: cronx
  max-idle-connections: 100
  max-open-connections: 100
  max-connection-life-time: 10s
  log-level: 4

redis:
  addr: 127.0.0.1:6379
  # username: cronx
  password: cronx
  database: 0

log:
  level: debug
  format: text
  output: [ "stdout" ]
```

其他配置没什么好说的，如果你使用上面的 `docker-compose.yaml` 部署依赖，则无需修改。这里唯一需要注意的配置项是 `kubeconfig`
，这是 CronX 能够连接 Kubernetes 的配置文件路径。

#### 6. 运行程序

```shell
$ ./bin/cronx -c configs/cronx.yaml
time=2026-02-28T15:04:05.673+08:00 level=WARN msg="Failed to create in-cluster config" err="unable to load in-cluster configuration, KUBERNETES_SERVICE_HOST and KUBERNETES_SERVICE_PORT must be defined"
time=2026-02-28T15:04:05.709+08:00 level=INFO msg="Starting cronx..."
time=2026-02-28T15:04:05.717+08:00 level=INFO msg="Successfully acquired lock" lockName=cronx:lock
time=2026-02-28T15:04:05.718+08:00 level=INFO msg="Job watcher init"
time=2026-02-28T15:04:05.721+08:00 level=INFO msg="Job watcher initiated"
time=2026-02-28T15:04:05.721+08:00 level=INFO msg="Successfully registers watchers"
time=2026-02-28T15:04:05.721+08:00 level=INFO msg="Successfully started cronx server"
time=2026-02-28T15:04:05.721+08:00 level=INFO msg=start
...
```

程序启动后，便可以观察 MySQL 表中 job 的状态，以及 Kubernetes 集群中 Job 的运行结果。

**注意**：CronX 中只有一个 Watcher 实现，那便是 jobWatcher。而 jobWatcher 依赖 Kubernetes 集群，如果你没有 Kubernetes 集群，并且对
Kubernetes 也不熟悉，可以参考 jobWatcher 实现自己的业务逻辑。CronX 本身通过 Watcher 接口支持灵活扩展，并不强依赖
Kubernetes。

### 将 CronX 部署到 Kubernetes 集群

按照如下步骤能够在本机 Kubernetes 集群中部署 CronX。

#### 1. 搭建 Kubernetes 集群

这里以 kind 为例在本机安装 Kubernetes 集群。

**注意**：如果你不想使用 kind 创建 Kubernetes 集群，后续操作的 Kubernetes 资源 yaml 文件中的 storageClassName 配置项可能需要修改成当前集群中的
StorageClass 名称。

参考如下文档安装 kind（用于创建 Kubernetes 集群）：

- https://kind.sigs.Kubernetes.io/docs/user/quick-start/#installation

参考如下文档安装 kubectl（用于操作 Kubernetes 集群）：

- https://kubernetes.io/zh-cn/docs/tasks/tools/install-kubectl-macos/

创建一个单节点的 Kubernetes 集群：

```shell
# 创建集群
$ kind create cluster
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.32.2) 🖼
 ✓ Preparing nodes 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Have a nice day! 👋

# 查看集群节点信息
$ kubectl get nodes
NAME                 STATUS   ROLES           AGE   VERSION
kind-control-plane   Ready    control-plane   29s   v1.32.2
```

#### 2. 安装依赖

在 Kubernetes 集群中安装部署 CronX 依赖的 MySQL 和 Redis。

```shell
# 安装资源
$ kubectl apply -f deployments/cronx-infrastructure.yaml 
namespace/cronx created
persistentvolumeclaim/mysql-pvc created
persistentvolumeclaim/redis-pvc created
service/mysql created
service/redis created
deployment.apps/mysql created
deployment.apps/redis created

# 资源准备就绪
$ kubectl -n cronx get all
NAME                         READY   STATUS    RESTARTS   AGE
pod/mysql-847c78b48c-hs4n9   1/1     Running   0          6s
pod/redis-6dd7b68f9-s5lgv    1/1     Running   0          6s

NAME            TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
service/mysql   ClusterIP   10.96.253.190   <none>        3306/TCP   7s
service/redis   ClusterIP   10.96.66.29     <none>        6379/TCP   7s

NAME                    READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/mysql   1/1     1            1           7s
deployment.apps/redis   1/1     1            1           7s

NAME                               DESIRED   CURRENT   READY   AGE
replicaset.apps/mysql-847c78b48c   1         1         1       7s
replicaset.apps/redis-6dd7b68f9    1         1         1       7s
```

#### 3. 准备 SQL

在 MySQL 中执行 `assets/schema.sql` 文件中的 SQL 语句准备测试数据。

```shell
$ kubectl -n cronx exec -i mysql-847c78b48c-hs4n9 -- \
  mysql -uroot -pcronx < assets/schema.sql
```

#### 4. 构建 CronX 镜像

构建当前平台镜像：

```shell
$ make TAG=v0.0.1 GOOS=linux GOARCH=amd64 docker-build
```

通过此命令构建的镜像名称为 `cronx/cronx:v0.0.1`。

> 注：你也可以考虑通过 `make docker-buildx` 命令构建跨平台镜像。

#### 5. 部署 CronX 到 Kubernetes 集群

将构建的本地镜像导入通过 kind 创建的 Kubernetes 集群中：

```shell
$ kind load docker-image cronx/cronx:v0.0.1 --name kind
Image: "cronx/cronx:v0.0.1" with ID "sha256:98f3f5aa31e69a641a33fe3767d80c323052475ac6f001b46f4ea19ab6e8f4c5" not yet present on node "kind-control-plane", loading...
```

> 注：你也可以通过此方式导入任何镜像到 Kubernetes 集群中。

部署 CronX 到 Kubernetes 集群：

```shell
$ make deploy TAG=v0.0.1
Deploying using Kubernetes manifests...
sed -e "s|image: .*|image: cronx/cronx:v0.0.1|" deployments/cronx.yaml | kubectl apply -f -
namespace/cronx unchanged
serviceaccount/cronx-sa created
role.rbac.authorization.k8s.io/cronx-role created
rolebinding.rbac.authorization.k8s.io/cronx-rb created
configmap/cronx-config created
deployment.apps/cronx created
```

等待 CronX Pod 全部启动成功：

```shell
$ kubectl -n cronx get pod
NAME                     READY   STATUS         RESTARTS   AGE
cronx-6dc594788b-ffgj2   1/1     Running        0          33s
cronx-6dc594788b-ghnkp   1/1     Running        0          33s
cronx-6dc594788b-zrl6j   1/1     Running        0          33s
...
```

3 个副本中只会有 1 个副本抢到分布式锁并执行任务调度：

```shell
$ kubectl -n cronx logs cronx-6dc594788b-ghnkp
{"time":"2026-02-28T15:04:05.419512676Z","level":"INFO","msg":"Starting cronx..."}
{"time":"2026-02-28T15:04:05.423082801Z","level":"INFO","msg":"Successfully acquired lock","lockName":"cronx:lock"}
{"time":"2026-02-28T15:04:05.423705384Z","level":"INFO","msg":"Job watcher init"}
{"time":"2026-02-28T15:04:05.426010926Z","level":"INFO","msg":"Job watcher initiated"}
{"time":"2026-02-28T15:04:05.428750218Z","level":"INFO","msg":"Successfully registers watchers"}
{"time":"2026-02-28T15:04:05.429047759Z","level":"INFO","msg":"Successfully started cronx server"}
{"time":"2026-02-28T15:04:05.429335176Z","level":"INFO","msg":"start"}
{"time":"2026-02-28T15:04:05.429676509Z","level":"INFO","msg":"schedule","now":"2026-02-28T15:04:05.429604134Z","entry":1,"next":"2026-02-28T15:04:35Z"}
{"time":"2026-02-28T15:04:05.532283551Z","level":"INFO","msg":"Started workers"}
```

稍等片刻就会有两个 Job 正在执行：

```shell
$ kubectl -n cronx get job                   
NAME         STATUS     COMPLETIONS   DURATION   AGE
demo-job-1   Running    0/1           53s        53s
demo-job-2   Complete   1/1           27s        53s
```

查看任务日志：

```shell
$ kubectl -n cronx logs demo-job-2-mfqz2
Hello Cronx!
```

Job 状态会实时同步至 MySQL 中：

```shell
$  kubectl -n cronx exec -i mysql-847c78b48c-hs4n9 -- \ 
  mysql -uroot -pcronx -e 'use cronx; select * from job\G;'
mysql: [Warning] Using a password on the command line interface can be insecure.
*************************** 1. row ***************************
        id: 1
      name: demo-job-1
 namespace: cronx
      info: {"image":"alpine","command":["sleep"],"args":["60"]}
    status: Running
created_at: 2026-02-28 15:04:05
updated_at: 2026-02-28 15:05:05
*************************** 2. row ***************************
        id: 2
      name: demo-job-2
 namespace: cronx
      info: {"image":"busybox","command":["echo"],"args":["Hello Cronx!"]}
    status: Succeeded
created_at: 2026-02-28 15:04:05
updated_at: 2026-02-28 15:05:07
```

至此，整个调试、部署过程全部结束。

**注意**：你可能已经注意到了，将 CronX 部署到 Kubernetes 集群时并没有修改 `kubeconfig` 配置项，这是因为 jobWatcher 代码已经做了兼容，识别在 Kubernetes 集群中运行 CronX 时，则通过 in-cluster 方式与 Kubernetes 进行连接。

> 注：如果你体验完成，想要清理 Kubernetes 集群，可以执行 `kind delete cluster` 删库跑路，谨慎操作 :)。
