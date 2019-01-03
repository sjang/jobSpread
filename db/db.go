package db

import (
  "fmt"
  _ "github.com/go-sql-driver/mysql"
  "github.com/jmoiron/sqlx"
  "jobSpread/config"
  "log"
)

var db *sqlx.DB

func InitDB(cf *config.Config) {

  dataSourceName := fmt.Sprintf("%s:%s@(%s:%s)/%s", cf.DB.Username, cf.DB.Password, cf.DB.ServerHost, cf.DB.ServerPort, cf.DB.DBName)
  fmt.Println(dataSourceName)
  dbHandle, err := sqlx.Connect(cf.DB.DBMSName, dataSourceName)
  if err != nil {
    log.Fatalln(err)
  }

  db.SetMaxIdleConns(cf.DB.MaxIdleConns)
  db.SetMaxOpenConns(cf.DB.MaxOpenConns)

  db = dbHandle
}

func DB() *sqlx.DB {
 return db
}

func releaseDB() {
  db.Close()
}

