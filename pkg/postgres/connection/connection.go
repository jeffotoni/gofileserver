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

package connection

import (
	"database/sql"
	"fmt"
	"os"

	sfconfig "github.com/jeffotoni/gofileserver/config"
	_ "github.com/lib/pq"
)

var (
	err error
	DBX *sql.DB
)

func Connect() *sql.DB {

	if DBX != nil {

		return DBX

	} else {

		cfg := sfconfig.GetConfig()

		var infoDb = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
			cfg.Section.User, cfg.Section.Password, cfg.Section.Database, "disable")

		DBX, errdb := sql.Open("postgres", infoDb)

		if errdb != nil {
			panic(fmt.Sprintf("Unable to connection to database: %v\n", errdb))
		}

		return DBX
	}
}

func TestDb() string {

	db := Connect()
	var returns string

	if db.Ping() == nil {

		returns = "ok"

	} else {

		returns = "erro"
	}

	if returns != "ok" {

		fmt.Println("erro connection to database Check your settings: ", db.Ping())
		os.Exit(1)
	}

	return returns
}
