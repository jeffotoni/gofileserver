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
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	sfconfig "github.com/jeffotoni/gofileserver/config"
	"github.com/jeffotoni/gofileserver/pkg/fcrypt"
	"github.com/jeffotoni/gofileserver/pkg/gofslib"
	"github.com/jeffotoni/gofileserver/pkg/postgres/connection"
)

/** Environment variables and keys */

var (
	confServer    *http.Server
	AUTHORIZATION = `bc8c154ebabc6f3da724e9x5fef79238`
	socketfileTmp = `gofileserver.red`
	socketfile    = `gofileserver.lock`
	keyCrypt      = `pKv9MQQIDAQABAmEApvlExjvPp0mYs/i`
)

func main() {

	// Command line for start and stop server

	if len(os.Args) > 1 {

		command := os.Args[1]

		if command != "" {

			if command == "start" {

				// Start server

				startFileServer()

			} else if command == "stop" {

				// Stop server

				fmt.Println("stop service...")
				NewRequestGetStop()

			} else {

				fmt.Println("Usage: gofileserver {start|stop}")
			}

		} else {

			command = ""
			fmt.Println("No command given")
		}
	} else {

		fmt.Println("Usage: gofileserver {start|stop}")
	}
}

func StopListenAndServe() {

	fmt.Println("stopping Server File...")
	confServer.Close()
	confServer.Shutdown(nil)
	RemoveFile(socketfile)
}

func startFileServer() {

	cfg := sfconfig.GetConfig()

	fmt.Println("Testing services")
	fmt.Println("Postgres: ", connection.TestDb())
	fmt.Println("Config: ", sfconfig.TestConfig())

	fmt.Println("Host: " + cfg.Section.Host)
	fmt.Println("Schema: http")
	fmt.Println("Server listening port : ", cfg.Section.ServerPort)
	fmt.Println("Database", cfg.Section.Database)
	fmt.Println("Database User: ", cfg.Section.User)

	registerUrl := cfg.Section.Schema + "://" + cfg.Section.Host + ":" + cfg.Section.ServerPort + "/register"
	tokenUrl := cfg.Section.Schema + "://" + cfg.Section.Host + ":" + cfg.Section.ServerPort + "/token"
	uploadUrl := cfg.Section.Schema + "://" + cfg.Section.Host + ":" + cfg.Section.ServerPort + "/upload"
	downloadUrl := cfg.Section.Schema + "://" + cfg.Section.Host + ":" + cfg.Section.ServerPort + "/download"

	fmt.Println("Instance POST ", registerUrl)
	fmt.Println("Instance GET  ", tokenUrl)
	fmt.Println("Instance POST ", uploadUrl)
	fmt.Println("Instance GET  ", downloadUrl)
	fmt.Println("Loaded service")

	///create route

	router := mux.NewRouter().StrictSlash(true)
	//router.Host("Localhost")

	router.Handle("/", http.FileServer(http.Dir("../dirmsg")))

	router.
		HandleFunc("/stop/{id}", func(w http.ResponseWriter, r *http.Request) {

			if r.Method == "GET" {

				HeaderAutorization := r.Header.Get("Authorization")

				fmt.Println("HeaderAutorization: ", HeaderAutorization)

				if HeaderAutorization == "" {

					fmt.Fprintln(w, "http ", 500, "Not authorized")

				} else {

					if HeaderAutorization == AUTHORIZATION {

						vars := mux.Vars(r)
						idStopServer := vars["id"]

						fmt.Println("Id Token: ", idStopServer)
						fmt.Println("Id TFile: ", ReadFile())

						if idStopServer == ReadFile() {

							fmt.Fprintln(w, "http ", 200, "ok stop")
							StopListenAndServe()

						} else {

							fmt.Fprintln(w, "http ", 500, "Not authorized")
						}

					} else {

						fmt.Fprintln(w, "http ", 500, "Not authorized")
					}
				}

			} else if r.Method == "POST" {

				fmt.Fprintln(w, "http ", 500, "Not authorized")

			} else {

				fmt.Fprintln(w, "http ", 500, "Not authorized")
			}
		})

	router.
		HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {

			if r.Method == "POST" {

				gofslib.RegisterUserJson(w, r)

			} else if r.Method == "GET" {

				fmt.Fprintln(w, "http ", 500, "Not authorized / Allowed method POST")

			} else {

				fmt.Fprintln(w, "http ", 500, "Not authorized / Allowed method POST")
			}
		})

	router.
		HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {

			if r.Method == "POST" {

				//gofslib.GetTokenUser(w, r)
				fmt.Fprintln(w, "http ", 500, "Not authorized / Allowed method GET")

			} else if r.Method == "GET" {

				gofslib.GetTokenUser(w, r)

			} else {

				fmt.Fprintln(w, "http ", 500, "Not authorized / Allowed method POST")
			}
		})

	router.
		HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {

			if r.Method == "POST" {

				gofslib.UploadFileEasy(w, r)

			} else if r.Method == "GET" {

				fmt.Fprintln(w, "http ", 500, "Not authorized / Allowed method POST")

			} else {

				fmt.Fprintln(w, "http ", 500, "Not authorized / Allowed method POST")
			}
		})

	router.
		HandleFunc("/download/{name}", func(w http.ResponseWriter, r *http.Request) {

			pathFileLocal := "../msg/error-download.txt"

			if r.Method == "GET" {

				gofslib.DownloadFile(w, r)

			} else if r.Method == "GET" {

				http.ServeFile(w, r, pathFileLocal)

				fmt.Fprintln(w, "http ", 500, "Not authorized")

			} else {

				http.ServeFile(w, r, pathFileLocal)
				fmt.Fprintln(w, "http ", 500, "Not authorized")
			}
		})

	confServer = &http.Server{

		Handler: router,
		Addr:    cfg.Section.Host + ":" + cfg.Section.ServerPort,

		// Good idea, good live!!!

		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	PASS_URL_MD5 := fcrypt.CreateTokenStrong()
	WriteFile(PASS_URL_MD5)

	log.Fatal(confServer.ListenAndServe())
}

/** [NewRequestGetStop This method will stop the server, it makes a Request Get for itself telling the server to stop] */

func NewRequestGetStop() {

	cfg := sfconfig.GetConfig()

	//Read generated key to generate the url

	PASS_URL_MD5 := ReadFile()

	if PASS_URL_MD5 != "" {

		// Create url to trigger a GET for us restful

		URL_STOP := cfg.Section.Schema + "://" + cfg.Section.Host + ":" + cfg.Section.ServerPort + "/stop/" + PASS_URL_MD5

		// Starting our instance to send a NewRequest to our restful

		client := &http.Client{}
		r, err := http.NewRequest("GET", URL_STOP, nil)
		if err != nil {

			fmt.Sprint(err)
			os.Exit(1)
		}

		r.Header.Add("Authorization", AUTHORIZATION)
		r.Header.Add("Accept", "application/text-plain")
		_, errx := client.Do(r)
		if errx != nil {

			log.Print(errx)
			os.Exit(1)
		}

	} else {

		fmt.Println("NewRequestGetStop Error URL ")
		os.Exit(1)
	}
}

/** [WriteFileCrypt Method that will encrypt the file with the key] */

func WriteFileCrypt(sname string) {

	// get key
	key := []byte(keyCrypt) // 32 bytes

	file, _ := os.Open(sname) // For read access.

	fi, _ := file.Stat()
	data := make([]byte, 16*fi.Size())
	count, _ := file.Read(data)

	file_cry, _ := os.Create(socketfile)
	defer file_cry.Close()

	ciphertext, _ := fcrypt.Encrypt(key, data[:count])
	file_cry.Write(ciphertext)

}

/** [WriteFile By recording the dynamic key in the file so that it can trigger the NewRequest for our restful] */

func WriteFile(sname string) {

	filename := socketfileTmp

	lcase := []byte(sname)
	perm := os.FileMode(0644)
	err := ioutil.WriteFile(filename, lcase, perm)

	if err != nil {
		fmt.Println("Error Write : ", filename, err)
		os.Exit(1)
	}

	// Encrypting the file with the generated key

	WriteFileCrypt(filename)

	// remove socketfileTmp

	os.Remove(filename)
}

/**  [ReadFile Read the key in the file so it can be used in dynamic URL]*/

func ReadFile() string {

	filename := socketfile
	key := []byte(keyCrypt) // 32 bytes

	_, err := ioutil.ReadFile(filename)

	if err != nil {

		// If the file does not exist create it empty

		WriteFile(" ")
		return ""
	}

	// If you do not return your content

	file_cry, _ := os.Open(filename) // For read access.
	ficry, _ := file_cry.Stat()
	data_cry := make([]byte, ficry.Size())
	count_cry, _ := file_cry.Read(data_cry)

	data_descry, _ := fcrypt.Decrypt(key, data_cry[:count_cry])

	return string(data_descry)
}

func RemoveFile(filename string) {

	os.Remove(filename)
}
