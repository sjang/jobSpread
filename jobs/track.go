package jobs

import (
  "github.com/sirupsen/logrus"
  "jobSpread/logging"
  "jobSpread/work"
  "sync"
)
type MetaRes []TrackInfo

func TrackMeta(worker *work.Pool, track_list []int64) []TrackInfo {
  log := logging.Log()
  log.Info("trackmeta start...")

  track_info_res_chan := make(chan *TrackInfo)

  metaRes := InitComposeResponse()

  jobCount := len(track_list)

  var wg sync.WaitGroup
  go func() {
    wg.Add(len(track_list))
    for i := range track_list {
      //fmt.Printf("t id:%d\n", track_list[i])
      track_req := &TrackReq{
        Id:               track_list[i],
        TrackInfoResChan: track_info_res_chan,
      }
      worker.Run(track_req)
      wg.Done()
    }
    wg.Wait()
  }()


  /////////////////////////////////////////////
  // aggregation from track_info_res_chan
  res := func(resCount int) MetaRes {
    defer func() {
      close(track_info_res_chan) // close res channel
    }()

    count := 0
    for track_info := range track_info_res_chan {
      //fmt.Printf("%v\n", track_info)
      metaRes = ComposeResponse(metaRes, track_info)
      count++
      if resCount == count {
        break
      }
    }
    log.WithFields(logrus.Fields{
      "test": "test...",
    }).Info("trackmeta finish...")
    return metaRes
  }(jobCount)
  /////////////////////////////////////////////

  //fmt.Println(res)

  return res
}

func InitComposeResponse() MetaRes {
  var metaRes = make(MetaRes, 0)
  return metaRes
}

func ComposeResponse(metaRes MetaRes, ti *TrackInfo) MetaRes {
  metaRes = append(metaRes, *ti)
  return metaRes
}


