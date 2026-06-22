//go:build ignore
// +build ignore

package main

import (
	"fmt"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote" // 导入远程驱动
)

func main() {
	// 连接远程 etcd 服务
	err := viper.AddRemoteProvider(
		"etcd3", "http://127.0.0.1:2379", "/config/app.yaml",
	)
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	err = viper.ReadRemoteConfig() // 读取远程配置
	if err != nil {
		panic(err)
	}

	fmt.Println(viper.GetString("app.name"))
}

/**

$ docker run -d \
  --name etcd \
  -p 2379:2379 \
  quay.io/coreos/etcd:v3.6.7 \
  /usr/local/bin/etcd \
  --name my-etcd \
  --advertise-client-urls http://localhost:2379 \
  --listen-client-urls http://0.0.0.0:237

$ etcdctl put /config/app.yaml "$(cat << 'EOF'
app:
  name: myapp
  mode: dev
  port: 8080
EOF
)"

*/
