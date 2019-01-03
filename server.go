package main

import (
  "encoding/json"
  "flag"
  "fmt"
  "github.com/jmoiron/sqlx"
  "github.com/labstack/echo"
  _  "github.com/labstack/echo/middleware"
  _ "github.com/labstack/gommon/log"
  "github.com/mediocregopher/radix/v3"
  "github.com/sirupsen/logrus"
  "jobSpread/config"
  "jobSpread/db"
  "jobSpread/ext"
  "jobSpread/jobs"
  "jobSpread/work"
  "net/http"
  "os"
  _ "os"
  "jobSpread/logging"
)

var MyEcho *echo.Echo
var log *logrus.Logger
var cf *config.Config
var dbHandle *sqlx.DB

type TrackIds struct {
  TrackList  []int64 `json:"trackList"`
}

func init() {
  // config(json)
  initConfig()

  // for redis
  addrs := make([]string, 0, 0)
  nodeConnCount := 20
  for i := 0; i < len(cf.RedisServer); i++ {
    addrs = append(addrs, cf.RedisServer[i])
  }
  ext.InitRedis(addrs, nodeConnCount)

  // job handle pooling
  maxWorkerPoolSize := 400
  work.CreateWorkerPool(maxWorkerPoolSize)

  // set test data to redis
  setTrackData()

  // init logrus
  logging.InitLog(cf)
  log = logging.Log()

  // init database
  // initDB(cf)
}

func main() {

  log.WithFields(logrus.Fields{
    "what": "meta info using fan out",
  }).Info("start the server")

  // echo init
  MyEcho = echo.New()

  // route
  dispatchers(MyEcho)

  // start echo
  log.Fatal(MyEcho.Start(":" + cf.ServerPort))
}

func dispatchers(e *echo.Echo) {
  // example
  e.GET("/", func(c echo.Context) error {
    return c.String(http.StatusOK, "Hello, World Echo!")
  })

  // Aggregation of track meta info
  e.POST("/track", track)

}

// track api
func track(c echo.Context) error {

  worker := work.MyWorker()
  defer work.ReturnWorker(worker)

  t := new(TrackIds)
  if err := c.Bind(t); err != nil {
    return err
  }
  fmt.Println(*t)
  res := jobs.TrackMeta(worker, t.TrackList)

  return c.JSON(http.StatusOK, res)
}

func initConfig () {

  configFlag := flag.Bool("f", false, "-f config.json")
  flag.Parse()

  //  currentDir := getCurrentDir()
  var configFilepath string
  if *configFlag == true {
    configFilepath = flag.Args()[0]
  } else {
    fmt.Println("Fail to load config...")
    os.Exit(1)
  }
  fmt.Println("Load config... ", configFilepath)

  err := config.LoadConfig(configFilepath)
  if err != nil {
    panic(err)
  }

  cf = config.Conf()
  // conf test
  //fmt.Println(cf.ServerPort)
  //fmt.Println(cf.LogConfig.LogDir)
  //fmt.Println(cf.LogConfig.LogFilename)
  //fmt.Println(cf.DB.ServerHost)
  //fmt.Println(cf.DB.ServerPort)
}

func initDB(cf *config.Config) {
  db.InitDB(cf)
  dbHandle = db.DB()
}



////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func setTrackData() {
  cluster := ext.RedisCluster()

  for i := 100; i < 110; i++ {
    serialized, err := json.Marshal(
      &jobs.TrackInfo{
      Id: int64(i),
      Title: fmt.Sprintf("title-%d", i),
      Artist: jobs.Artist{ Id: int64(i), ArtistName: fmt.Sprintf("artistName-%d", i)},
        })
    if err != nil {
      fmt.Println("marshal error")
    }
    err = cluster.Do(radix.Cmd(nil, "SET", fmt.Sprintf("track_%d", i), string(serialized)))
    if err != nil {
      fmt.Println("redis error")
    }
  }
}
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////


/*
curl -v -X POST -H "Content-Type: application/json" -d '{"trackList":[100,102,101,103,104,105]}' http://localhost:1323/track
*/

