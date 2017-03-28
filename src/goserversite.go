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

	"github.com/gorilla/mux"
)

func main() {

	// Opening an instance to serve the dynamic site, the project interface

	// fs := http.FileServer(http.Dir("../views"))
	// http.Handle("/", fs)

	// log.Println("Listening 8080 site...")
	// http.ListenAndServe(":8080", nil)
	//

	var dir string

	flag.StringVar(&dir, "../views", ".", "directory to serve files from. Defaults to the current dir")
	flag.Parse()

	r := mux.NewRouter()

	// This will serve files under http://localhost:8080/static

	r.PathPrefix("/views/").Handler(
		http.StripPrefix("/views/",
			http.FileServer(http.Dir(dir))))

	srv := &http.Server{

		Handler: r,
		Addr:    "127.0.0.1:8081",

		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Start server port: 8081")

	log.Fatal(srv.ListenAndServe())

}
