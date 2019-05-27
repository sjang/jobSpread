package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	ServerPort  string    `json:"serverPort,omitempty"`
	RedisServer []string  `json:"redisServer,omitempty"`
	WorkHomeDir string    `json:"workHomeDir,omitempty"`
	TempDir     string    `json:"tempDir,omitempty"`
	LogConfig   logConfig `json:"logConfig,omitempty"`
	BackupDir   string    `json:"backupDir,omitempty"`
	LogDir      string    `json:"logDir,omitempty"`
	DB          *DBConfig `json:"db,omitempty"` // for reservation
}

type logConfig struct {
	LogDir            string `json:"logDir,omitempty"`
	LogFilename       string `json:"logFilename,omitempty"`
	LogFilenameFormat string `json:"logFilenameFormat,omitempty"`
	BackupDir         string `json:"backupDir,omitempty"`
}

type DBConfig struct {
	DBMSName     string `json:"dbmsName,omitempty"`
	ServerHost   string `json:"serverHost,omitempty"`
	ServerPort   string `json:"serverPort,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	DBName       string `json:"dbName,omitempty"`
	MaxIdleConns int    `json:"maxIdleConns,omitempty"`
	MaxOpenConns int    `json:"maxOpenConns,omitempty"`
}

var conf *Config

func LoadConfig(filepath string) error {
	var config Config
	configFile, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	errJson := jsonParser.Decode(&config)
	if errJson != nil {
		return errJson
	}

	// load
	conf = &config

	return nil
}

func Conf() *Config {
	return conf
}
