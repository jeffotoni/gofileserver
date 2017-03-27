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

	"github.com/gorilla/mux"
	sfconfig "github.com/jeffotoni/gofileserver/config"
	"github.com/jeffotoni/gofileserver/pkg/gofslib"
)

func main() {

	//config global

	cfg := sfconfig.GetConfig()

	fmt.Println("Conect port : ", cfg.Section.ServerPort)
	fmt.Println("Conect database: ", cfg.Section.Database)
	fmt.Println("Database User: ", cfg.Section.User)

	fmt.Println("Instance /register")
	fmt.Println("Instance /token")
	fmt.Println("Instance /upload")
	fmt.Println("Instance /download")

	///create route

	router := mux.NewRouter()

	router.Handle("/", http.FileServer(http.Dir("../dirmsg")))

	router.
		HandleFunc("/register", gofslib.RegisterUserJson).
		Methods("POST")

	router.
		HandleFunc("/token", gofslib.GetTokenUser).
		Methods("POST")

	router.
		Path("/upload").
		HandlerFunc(gofslib.UploadFileEasy).
		Methods("POST")

	router.
		Path("/download/{name}").
		HandlerFunc(gofslib.DownloadFile).
		Methods("GET")

	//After 5 minutes synchronize file upload
	//After 10 minutes synchronize file remove

	// port in config.gcfg

	port := cfg.Section.ServerPort

	log.Fatal(http.ListenAndServe(port, router))

}
