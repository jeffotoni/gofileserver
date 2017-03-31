/***********

 ▄▄▄██▀▀▀▓█████   █████▒ █████▒▒█████  ▄▄▄█████▓ ▒█████   ███▄    █  ██▓
   ▒██   ▓█   ▀ ▓██   ▒▓██   ▒▒██▒  ██▒▓  ██▒ ▓▒▒██▒  ██▒ ██ ▀█   █ ▓██▒
   ░██   ▒███   ▒████ ░▒████ ░▒██░  ██▒▒ ▓██░ ▒░▒██░  ██▒▓██  ▀█ ██▒▒██▒
▓██▄██▓  ▒▓█  ▄ ░▓█▒  ░░▓█▒  ░▒██   ██░░ ▓██▓ ░ ▒██   ██░▓██▒  ▐▌██▒░██░
 ▓███▒   ░▒████▒░▒█░   ░▒█░   ░ ████▓▒░  ▒██▒ ░ ░ ████▓▒░▒██░   ▓██░░██░
 ▒▓▒▒░   ░░ ▒░ ░ ▒ ░    ▒ ░   ░ ▒░▒░▒░   ▒ ░░   ░ ▒░▒░▒░ ░ ▒░   ▒ ▒ ░▓
 ▒ ░▒░    ░ ░  ░ ░      ░       ░ ▒ ▒░     ░      ░ ▒ ▒░ ░ ░░   ░ ▒░ ▒ ░
 ░ ░ ░      ░    ░ ░    ░ ░   ░ ░ ░ ▒    ░      ░ ░ ░ ▒     ░   ░ ░  ▒ ░
 ░   ░      ░  ░                  ░ ░               ░ ░           ░  ░

*
* @about 	project gofileserver / server-upload / donwload / register user
* @autor 	jeffotoni
* @date 	25/03/2017
* @since    Version 0.1
*/

package config

import (
	"fmt"
	"os"

	gcfg "gopkg.in/gcfg.v1"
)

type Config struct {
	Section struct {
		Domain     string
		Port       string
		Database   string
		User       string
		Password   string
		Process    string
		UploadMax  string
		UploadSize int64
		Ping       string
		Bucket     string
		PathLocal  string
		ServerPort string
		SitePort   string
		Host       string
		Schema     string
	}
}

var cfg Config

func GetConfig() *Config {

	//Is at the root of your project

	err := gcfg.ReadFileInto(&cfg, "../config/config.gcfg")

	if cfg.Section.Ping != "ok" {

		fmt.Println("Error reading file config.gcfg ", err)
	}

	return &cfg
}

func TestConfig() string {

	var returns string

	//path, _ := os.Getwd()

	msgerr := "Error reading the config file!"

	if cfg.Section.Ping == "ok" {

		returns = "ok"

	} else {

		returns = "erro"
	}

	if returns != "ok" {

		fmt.Println(msgerr)
		os.Exit(1)
	}

	return returns
}
