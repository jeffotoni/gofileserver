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
*
* projeto server-upload / donwload / register user
*
* @package     site dinamic
* @author      jeffotoni
* @copyright   Copyright (c) 2017
* @license     --
* @link        --
* @since       Version 0.1
*/

// package main

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	sfconfig "github.com/jeffotoni/gofileserver/config"
)

func main() {

	var dir string

	cfg := sfconfig.GetConfig()

	flag.StringVar(&dir, "views", "../views", "directory to serve files from. Defaults to the current dir")
	flag.Parse()

	// This will serve files under http://localhost:port

	fmt.Println(dir)
	fmt.Println("Start server port:" + cfg.Section.SitePort)

	statics := http.FileServer(http.Dir(dir))
	http.Handle("/", statics)

	sitef := &http.Server{

		Handler: nil,
		Addr:    "127.0.0.1:" + cfg.Section.SitePort,

		// Good practice!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(sitef.ListenAndServe())

}
