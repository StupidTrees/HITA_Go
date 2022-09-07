package config

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"hita/utils/logger"
	"io/ioutil"
	"log"
	"os"
)

var MysqlIp string
var MysqlUsername string
var MysqlPassword string
var MysqlPort string
var MysqlDbname string
var ServerPrivateKey string
var ServerPublicKey string
var PORT string
var AvatarPath string
var ArticleImagePath string
var SensitivePlaceholderPath string

type ClientInfo struct {
	ClientSecret string
	ClientToken  string
	TokenType    string
}

func init() {
}

//从文件中读取配置信息
func loadFromConfigFile(configFilePath string) error {

	file, err := os.Open(configFilePath)

	if err != nil {
		log.Println(err)
		return err
	}

	data, err := ioutil.ReadAll(file)

	if nil != err {
		return err
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	var jsonConfig jsoniter.Any = json.Get(data)

	MysqlIp = jsonConfig.Get("MYSQL_IP").ToString()
	MysqlUsername = jsonConfig.Get("MYSQL_USERNAME").ToString()
	MysqlPassword = jsonConfig.Get("MYSQL_PASSWORD").ToString()
	MysqlDbname = jsonConfig.Get("MYSQL_DBNAME").ToString()
	MysqlPort = jsonConfig.Get("MYSQL_PORT").ToString()
	PORT = jsonConfig.Get("PORT").ToString()
	ServerPublicKey = jsonConfig.Get("SERVER_PUBLIC_KEY").ToString()
	ServerPrivateKey = jsonConfig.Get("SERVER_PRIVATE_KEY").ToString()
	AvatarPath = jsonConfig.Get("AVATAR_PATH").ToString()
	ArticleImagePath = jsonConfig.Get("ARTICLE_IMAGE_PATH").ToString()
	SensitivePlaceholderPath = jsonConfig.Get("SENSITIVE_PLACEHOLDER_PATH").ToString()
	if MysqlIp == "" || MysqlUsername == "" || MysqlPassword == "" || MysqlPort == "" || PORT == "" || MysqlDbname == "" {
		return errors.New("config is error")
	}

	return nil
}

func LoadConfig(configFile string) error {

	err := loadFromConfigFile(configFile)
	if nil != err {
		logger.Println("Failed to load config,Error:" + err.Error())
		return err
	}

	return nil
}
