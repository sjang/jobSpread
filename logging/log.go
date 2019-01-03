package logging

import (
  "github.com/lestrrat-go/file-rotatelogs"
  "github.com/sirupsen/logrus"
  "jobSpread/config"
  "time"
)

var log *logrus.Logger
var cf *config.Config

func InitLog(cf *config.Config) {
  log = logrus.New()
  log.SetFormatter(&logrus.JSONFormatter{})
  logf, err := rotatelogs.New(
    cf.LogConfig.LogDir + cf.LogConfig.LogFilename + cf.LogConfig.LogFilenameFormat,
  rotatelogs.WithLinkName(cf.LogConfig.LogDir + cf.LogConfig.LogFilename),
    rotatelogs.WithMaxAge(24 * time.Hour),
//    rotatelogs.WithRotationTime(time.Hour),
  )
  if err != nil {
    log.Printf("failed to create rotatelogs: %s", err)
    return
  }

  log.Out = logf
}

func Log() *logrus.Logger {
  return log
}
