package ext

import (
  "fmt"
  "github.com/mediocregopher/radix/v3"
  "time"
)

var Cluster *radix.Cluster

func InitRedis(addrs []string, nodeConnCount int) {
  // this is a ConnFunc which will set up a connection which is authenticated
  // and has a 1 minute timeout on all operations
  customConnFunc := func(network, addr string) (radix.Conn, error) {
    return radix.Dial(network, addr,
      radix.DialTimeout(5*time.Second),
      radix.DialAuthPass("fpeltm123!"),
      radix.DialReadTimeout(10*time.Second),
    )
  }

  // this cluster will use the ClientFunc to create a pool to each node in the
  // cluster. The pools also use our customConnFunc, but have more connections
  poolFunc := func(network, addr string) (radix.Client, error) {
    return radix.NewPool(network, addr, nodeConnCount, radix.PoolConnFunc(customConnFunc))
  }

  cluster, err := radix.NewCluster(addrs, radix.ClusterPoolFunc(poolFunc))
  if err != nil {
    fmt.Println("redis error")
    Cluster = nil
  }

  Cluster = cluster
}

func RedisCluster() *radix.Cluster {
  return Cluster
}
