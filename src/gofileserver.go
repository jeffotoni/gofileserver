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
* project server-upload / donwload / register user
*
* @package     main
* @author      jeffotoni
* @copyright   Copyright (c) 2017
* @license     --
* @link        --
* @since       Version 0.1
*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	sfconfig "github.com/jeffotoni/gofileserver/config"
	"github.com/jeffotoni/gofileserver/pkg/gofslib"
	"github.com/jeffotoni/gofileserver/pkg/postgres/connection"
)

func main() {

	//config global

	cfg := sfconfig.GetConfig()

	fmt.Println("Testing services")
	fmt.Println("Postgres: ", connection.TestDb())
	fmt.Println("Config: ", sfconfig.TestConfig())

	fmt.Println("Server listening port : ", cfg.Section.ServerPort)
	fmt.Println("Database", cfg.Section.Database)
	fmt.Println("Database User: ", cfg.Section.User)

	serverHost := cfg.Section.serverHost
	schema := cfg.Section.schema

	ping := schema + "://" + serverHost + ":" + cfg.Section.ServerPort + "/ping"
	register := schema + "://" + serverHost + ":" + cfg.Section.ServerPort + "/register"
	token := schema + "://" + serverHost + ":" + cfg.Section.ServerPort + "/token"
	upload := schema + "://" + serverHost + ":" + cfg.Section.ServerPort + "/upload"
	download := schema + "://" + serverHost + ":" + cfg.Section.ServerPort + "/download"

	fmt.Println("[Instance POST] " + ping)
	fmt.Println("[Instance POST] " + register)
	fmt.Println("[Instance GET] " + token)
	fmt.Println("[Instance POST] " + upload)
	fmt.Println("[Instance GET] " + download)

	///create route

	router := mux.NewRouter().StrictSlash(true)

	router.Handle("/", http.FileServer(http.Dir("../dirmsg")))

	router.
		HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {

			if r.Method == http.MethodPost || r.Method == http.MethodGet {

				fmt.Fprintln(w, "http ", 200, `{"msg":"pong"}`)

			} else if r.Method == "GET" {

				fmt.Fprintln(w, "http ", 500, "Not authorized / Allowed method POST")

			} else {

				fmt.Fprintln(w, "http ", 500, "Not authorized / Allowed method POST")
			}
		})

	router.
		HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {

			if r.Method == http.MethodPost {

				gofslib.RegisterUserJson(w, r)

			} else if r.Method == http.MethodGet {

				fmt.Fprintln(w, "http ", 500, "Not authorized / Allowed method POST")

			} else {

				fmt.Fprintln(w, "http ", 500, "Not authorized / Allowed method POST")
			}
		})

	router.
		HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {

			if r.Method == http.MethodPost {

				//gofslib.GetTokenUser(w, r)
				fmt.Fprintln(w, "http ", 500, "Not authorized / Allowed method GET")

			} else if r.Method == http.MethodGet {

				gofslib.GetTokenUser(w, r)

			} else {

				fmt.Fprintln(w, "http ", 500, "Not authorized / Allowed method POST")
			}
		})

	router.
		HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {

			if r.Method == http.MethodPost {

				gofslib.UploadFileEasy(w, r)

			} else if r.Method == http.MethodGet {

				fmt.Fprintln(w, "http ", 500, "Not authorized / Allowed method POST")

			} else {

				fmt.Fprintln(w, "http ", 500, "Not authorized / Allowed method POST")
			}
		})

	router.
		HandleFunc("/download/{name}", func(w http.ResponseWriter, r *http.Request) {

			pathFileLocal := "../msg/error-download.txt"

			if r.Method == http.MethodPost {

				gofslib.DownloadFile(w, r)

			} else if r.Method == http.MethodGet {

				http.ServeFile(w, r, pathFileLocal)

				fmt.Fprintln(w, "http ", 500, "Not authorized")

			} else {

				http.ServeFile(w, r, pathFileLocal)
				fmt.Fprintln(w, "http ", 500, "Not authorized")
			}
		})

	// http server conf

	confsc := &http.Server{

		Handler: router,
		Addr:    "127.0.0.1:" + cfg.Section.ServerPort,

		// Good idea, good live!!!

		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.Fatal(confsc.ListenAndServe())

}
