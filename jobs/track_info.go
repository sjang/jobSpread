package jobs

import (
  "jobSpread/ext"
  "encoding/json"
  "fmt"
  "github.com/mediocregopher/radix/v3"
  "strconv"
  "strings"
  "time"
)

type TrackInfo struct {
  Id int64      `json:"id"`
  Title string  `json:"title"`
  Artist Artist `json:"artist"`
  Album Album   `json:"album"`
}

type Artist struct {
  Id int64          `json:"id"`
  ArtistName string `json:"artistName"`
}

type Album struct {
  Id int64          `json:"id"`
  AlbumName string `json:"albumtName"`
}

type TrackReq struct {
  Id int64
  TrackInfoResChan chan *TrackInfo
  count int64
}

func (tr *TrackReq) InitTask() {
  return
}

func (tr *TrackReq) DoTask() {
//  fmt.Println("here DoTask:", tr.Id)

  redisCluster := ext.RedisCluster()
  trackKeyId := strings.Join([]string{"track_", strconv.FormatInt(tr.Id, 10)}, "")

  trackRes := &TrackInfo{}

  var metaValue []byte
  start := time.Now()
  err := redisCluster.Do(radix.Cmd(&metaValue, "GET", trackKeyId))
  elapsed := time.Since(start)
  fmt.Printf("redis time:%s\n", elapsed)

  if err != nil {
    // error
    fmt.Printf("redis error(%s):%d\n", err.Error(),tr.Id)
    tr.TrackInfoResChan <- trackRes
    return
  }
  //fmt.Println("redis value:", metaValue)
  if len(metaValue) == 0 { // no exists value
    // print error
  } else {
    // after redis operation /////////////////////
    // deserialzation
    err := json.Unmarshal(metaValue, trackRes)
    if err != nil {
      return
    }
  }
  ///////////////////////////////////////////////////

  tr.TrackInfoResChan <- trackRes
  return
}

func (tr *TrackReq) ErrorTask() {
//  fmt.Println("here ErrorTask:", tr.Id)
  return
}

func (tr *TrackReq) PostTask() {
//  fmt.Println("here PostTask:", tr.Id)
  return
}

//func (ti *TrackInfo) String() {
//  fmt.Printf("id:[%d], Title: [%s], Artists:[%s]", ti.Id, ti.Title, ti.Artists)
//}
