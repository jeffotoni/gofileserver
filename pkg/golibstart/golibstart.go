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

package golibstart

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	sfconfig "github.com/jeffotoni/gofileserver/config"
	"github.com/jeffotoni/gofileserver/pkg/fcrypt"
	"github.com/jeffotoni/gofileserver/pkg/gofile"
	"github.com/jeffotoni/gofileserver/pkg/gofslib"
	"github.com/jeffotoni/gofileserver/pkg/postgres/connection"
)

var (
	confServer    *http.Server
	AUTHORIZATION = `bc8c154ebabc6f3da724e9x5fef79238`
)

func StartFileServer() {

	cfg := sfconfig.GetConfig()

	color.Cyan("Testing services")
	color.Cyan("Postgres: " + connection.TestDb())
	color.Cyan("Config: " + sfconfig.TestConfig())
	color.Cyan("Upload: " + sfconfig.TestConfigUpload())
	color.Yellow("successfully...")

	htmlviews := cfg.Section.Schema + "://" + cfg.Section.ServerHost + ":" + cfg.Section.ServerPort + ""

	ping := cfg.Section.Schema + "://" + cfg.Section.ServerHost + ":" + cfg.Section.ServerPort + "/ping"
	registerUrl := cfg.Section.Schema + "://" + cfg.Section.ServerHost + ":" + cfg.Section.ServerPort + "/register"
	tokenUrl := cfg.Section.Schema + "://" + cfg.Section.ServerHost + ":" + cfg.Section.ServerPort + "/token"
	uploadUrl := cfg.Section.Schema + "://" + cfg.Section.ServerHost + ":" + cfg.Section.ServerPort + "/upload"
	downloadUrl := cfg.Section.Schema + "://" + cfg.Section.ServerHost + ":" + cfg.Section.ServerPort + "/download"

	color.Red("[page html] " + htmlviews)
	color.Red("[POST/GET] " + ping)
	color.Red("[POST] " + registerUrl)
	color.Red("[GET] " + tokenUrl)
	color.Red("[POST] " + uploadUrl)
	color.Red("[GET]  " + downloadUrl)

	color.Yellow("Starting service...")
	color.Green("Host: " + cfg.Section.ServerHost)
	color.Green("Schema: " + cfg.Section.Schema)
	color.Green("Port: " + cfg.Section.ServerPort)
	color.Green("Database: " + cfg.Section.Database)
	color.Green("User: " + cfg.Section.User)
	color.White("Loaded service")

	///create route

	router := mux.NewRouter().StrictSlash(true)

	router.
		HandleFunc("/stop/{id}", func(w http.ResponseWriter, r *http.Request) {

			if r.Method == http.MethodGet {

				HeaderAutorization := r.Header.Get("Authorization")

				fmt.Println("HeaderAutorization: ", HeaderAutorization)

				if HeaderAutorization == "" {

					fmt.Fprintln(w, "http ", 500, "Not authorized")

				} else {

					if HeaderAutorization == AUTHORIZATION {

						vars := mux.Vars(r)
						idStopServer := vars["id"]

						fmt.Println("Id Token: ", idStopServer)
						fmt.Println("Id TFile: ", gofile.ReadFile())

						if idStopServer == gofile.ReadFile() {

							fmt.Fprintln(w, "http ", 200, "ok stop")
							StopListenAndServe()

						} else {

							fmt.Fprintln(w, "http ", 500, "Not authorized...")
						}

					} else {

						fmt.Fprintln(w, "http ", 500, "Not authorized")
					}
				}

			} else if r.Method == http.MethodPost {

				fmt.Fprintln(w, "http ", 500, "Not authorized")

			} else {

				fmt.Fprintln(w, "http ", 500, "Not authorized")
			}
		})

	router.
		HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {

			if r.Method == http.MethodPost || r.Method == http.MethodGet {

				fmt.Fprintln(w, "http ", 200, `{"msg":"pong"}`)

			} else {

				fmt.Fprintln(w, "http ", 500, "Not authorized / Allowed method POST")
			}
		})

	router.
		HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {

			if r.Method == http.MethodPost {

				gofslib.RegisterUserJson(w, r)

			} else if r.Method == http.MethodPost {

				fmt.Fprintln(w, "http ", 500, "Not authorized / Allowed method POST")

			} else {

				fmt.Fprintln(w, "http ", 500, "Not authorized / Allowed method POST")
			}
		})

	router.
		HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {

			if r.Method == "POST" {

				gofslib.GetTokenUser(w, r)
				//fmt.Fprintln(w, "http ", 500, "Not authorized / Allowed method GET")

			} else if r.Method == http.MethodGet {

				gofslib.GetTokenUser(w, r)

			} else {

				fmt.Fprintln(w, "http ", 500, "Not authorized / Allowed method POST")
			}
		})

	router.
		HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {

			if r.Method == http.MethodPut {

				gofslib.UploadFileEasy(w, r)

			} else if r.Method == http.MethodPost {

				gofslib.UploadFileEasy(w, r)

			} else if r.Method == http.MethodGet {

				fmt.Fprintln(w, "http ", 500, "Not authorized / Allowed method POST")

			} else {

				fmt.Fprintln(w, "http ", 500, "Not authorized / Allowed method POST")
			}
		})

	router.
		HandleFunc("/download/{name}", func(w http.ResponseWriter, r *http.Request) {

			pathFileLocal := "msg/error-download.txt"

			if r.Method == http.MethodGet {

				gofslib.DownloadFile(w, r)

			} else if r.Method == http.MethodGet {

				http.ServeFile(w, r, pathFileLocal)

				fmt.Fprintln(w, "http ", 500, "Not authorized")

			} else {

				http.ServeFile(w, r, pathFileLocal)
				fmt.Fprintln(w, "http ", 500, "Not authorized")
			}
		})

	//router.Host("Localhost")
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("../views/"))))

	confServer = &http.Server{

		Handler: router,
		Addr:    cfg.Section.Host + ":" + cfg.Section.ServerPort,

		// Good idea, good live!!!
		//WriteTimeout: 10 * time.Second,
		//ReadTimeout:  10 * time.Second,
	}

	PASS_URL_MD5 := fcrypt.CreateTokenStrong()
	gofile.WriteFile(PASS_URL_MD5)

	log.Fatal(confServer.ListenAndServe())
}

func StopListenAndServe() {

	fmt.Println("stopping Server File...")
	confServer.Close()
	confServer.Shutdown(nil)
	gofile.RemoveFile(gofile.Socketfile)
}

/** [NewRequestGetStop This method will stop the server, it makes a Request Get for itself telling the server to stop] */

func NewRequestGetStop() {

	cfg := sfconfig.GetConfig()

	//Read generated key to generate the url

	PASS_URL_MD5 := gofile.ReadFile()

	if PASS_URL_MD5 != "" {

		// Create url to trigger a GET for us restful

		URL_STOP := cfg.Section.Schema + "://" + cfg.Section.ServerHost + ":" + cfg.Section.ServerPort + "/stop/" + PASS_URL_MD5

		// fmt.Println(URL_STOP)
		// Starting our instance to send a NewRequest to our restful

		client := &http.Client{}
		r, err := http.NewRequest("GET", URL_STOP, nil)
		if err != nil {

			fmt.Sprint(err)
			fmt.Println(err)
			os.Exit(1)
		}

		r.Header.Add("Authorization", AUTHORIZATION)
		r.Header.Add("Accept", "application/text-plain")
		_, errx := client.Do(r)
		if errx != nil {

			log.Print(errx)
			fmt.Println(errx)
			os.Exit(1)
		}

	} else {

		fmt.Println("NewRequestGetStop Error URL ")
		os.Exit(1)
	}
}
