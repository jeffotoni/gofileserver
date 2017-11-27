# gofileserver
The objective of this project is purely didactic, it is an attempt to solve the problem of competition and parallelism that the project will assume with uploading, reading and downloading files on disk and for the cloud.

Simple RESTFUL API demo server written in Go (golang) in order to solve problems of recording, uploading and downloading files to disk and the Amazon cloud using S3 for various regions, Google Cloud or any other cloud service that you want to implement, using postgresql as Database and token authentication.
The services of amazon s3 were first implemented.


## Used libraries:
- https://github.com/lib/pq - Sql Database.
- https://gopkg.in/gcfg.v1 - Text-based configuration files with "name=value" pairs grouped into sections (gcfg files).
- https://github.com/aws/aws-sdk-go/aws - Package sdk Aws Amazon
- https://github.com/gorilla/mux - Implements a request router and dispatcher for matching incoming requests

---
* [A small summary](#summary)
* [Install](#install)
* [Structure](#structure)
* [Run](#runprogram)
* [Examples Client](#examples-client)
* [Register User](#register-user)
* [Upload File](#upload-files)
* [Download Files](#download-files)

## A small summary operating system

* [Gofileserver.go]

This program is responsible for managing and controlling our APIS.

The upload is done on the local server always, there is a configuration file to determine some features of the config.gcfg system.
Everything is recorded in the database.
Authentication for upload upload was done in 2 ways.

The download checks if the file is local or if it has already been sent to the cloud, if it is still on the local server the system will download locally, otherwise it checks on which cloud server it is to download from the cloud.

The register and token is responsible for creating users and receiving the access token, and also with the possibility to view the tokem with username and password.


* [Gofileupload.go]

This program is responsible for downloading the files that are found on the local server to the cloud, the system evaluates the availability and uploads to the cloud.

It is not deleted files locally, they are deleted in a second time for security, they are checked if it really is in the cloud before physically removes them.

* [Gofileremove.go]

This program is responsible for monitoring the files on disk so they can be removed, it checks in the cloud if the object is actually there, and if it does remove it from the disk physically.

* [Goserversite.go]

This program is just an attempt to create an interface to simulate what the client will see how the files on our platform are.


## Install

```sh
git clone https://github.com/jeffotoni/gofileserver
```

With a [correctly configured](https://golang.org/doc/code.html#Workspaces) Go toolchain:

```sh
go get -u github.com/lib/pq
go get -u gopkg.in/gcfg.v1
go get -u github.com/aws/aws-sdk-go/aws
go get -u github.com/gorilla/mux
go get -u github.com/jeffotoni/gofileserver
```

Configuring the environment to run sdk amazon API
[`SDK Examples`] (https://github.com/aws/aws-sdk-go/tree/master/example)

```
mkdir ~/.aws/
vim ~/.aws/config

[default]
region = us-east-1
output = 

vim ~/.aws/credentials

[default]
aws_access_key_id = AKIX1234567890
aws_secret_access_key = MY-SECRET-KEY
```

Creating postgresql database 
[`Installing`] (Postgres http://www.postgresguide.com/setup/install.html)

```sh
createuser ukkobox -U postgres
psql -d template1 -U postgres
psql=> alter user ukkobox password 'pass123'
createdb ukkobox -U postgres -O ukkobox -E UTF-8 -T template0
psql ukkobox -U postgres -f tables/ukkobox.sql
```

Edit the configuration file

```
vim config/config.gcfg
```

Crontab edit the configuration / crontab -e
```
*/5 * * * * cd /pathprojeto/gofileserver/src && go run gofileupload.go >> /pathprojeto/gofileserver/src/gofileupload.log

*/10 * * * * cd /pathprojeto/gofileserver/src && go run gofileremove.go >> /pathprojeto/gofileserver/src/gofileremove.log

```

## Structure of the program

```go
- gofileserver
	- bin 
	- config
		- config.go
		- config.gcfg
	- dirmsg
		welcome.html
	- pkg	
		- fcrypt
		- gcheck
		- gofrlib
		- gofslib
		- gofuplib
		- postgres
			- connection
	- tables
		ukkobox.sql
	- uploads
	- views
		index.html

	gofileserver.go	
	Server to register users, receive and send files to the file server
	

	gofileupload.go
	Scans all local structure and sends the files to servers in the 
	cloud: Cloud Google, Amazon, DigitalOcean
	

	gofileremove.go
	Scrolls the structure checks whether the files are safe in the 
	cloud and removes the location
	
	goserversite.go
	A template under construction so customers can manage, view, download 
	their uploaded as an online bucket
	
```

## Run the program

```go
go run gofileserver.go 

Conect port : 80
Conect database:  ukkobox
Database User:  ukkobox
Instance /register
Instance /token
Instance /upload
Instance /download

OR New gofileserver 

go run gofileserver2.go start

Testing services
Postgres:  ok
Config:  ok
Host: localhost
Schema: http
Server listening port :  4001
Database ukkobox
Database User:  ukkobox
Instance POST  http://localhost:4001/register
Instance GET   http://localhost:4001/token
Instance POST  http://localhost:4001/upload
Instance GET   http://localhost:4001/download
Loaded service

```

Stopping the server

```go
go run gofileserver2.go stop
```

Compiling gofileserver or gofileserver2

```go
go build gofileserver.go 
go build gofileserver2.go 
```

Body of main function

```go

func main() {

	//config global

	cfg := sfconfig.GetConfig()

	fmt.Println("Server listening port : ", cfg.Section.ServerPort)
	fmt.Println("Database", cfg.Section.Database)
	fmt.Println("Database User: ", cfg.Section.User)

	fmt.Println("Instance POST /register")
	fmt.Println("Instance GET /token")
	fmt.Println("Instance POST /upload")
	fmt.Println("Instance GET /download")

	///create route
	router := mux.NewRouter().StrictSlash(true)

	router.Handle("/", http.FileServer(http.Dir("../dirmsg")))

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

	
	confsc := &http.Server{

		Handler: router,
		Addr:    "127.0.0.1:" + cfg.Section.ServerPort,

		// Good idea, good live!!!

		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.Fatal(confsc.ListenAndServe())

}

```

```go
go build gofileserver2.go 
```

Body of main function 

```go

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

```

Body of main function  startFileServer(){}

```go

func startFileServer() {

	cfg := sfconfig.GetConfig()

	fmt.Println("Testing services")
	fmt.Println("Postgres: ", connection.TestDb())
	fmt.Println("Config: ", sfconfig.TestConfig())

	fmt.Println("Host: Localhost")
	fmt.Println("Schema: http")
	fmt.Println("Server listening port : ", cfg.Section.ServerPort)
	fmt.Println("Database", cfg.Section.Database)
	fmt.Println("Database User: ", cfg.Section.User)

	fmt.Println("Instance POST http://localhost:" + cfg.Section.ServerPort + "/register")
	fmt.Println("Instance GET  http://localhost:" + cfg.Section.ServerPort + "/token")
	fmt.Println("Instance POST http://localhost:" + cfg.Section.ServerPort + "/upload")
	fmt.Println("Instance GET  http://localhost:" + cfg.Section.ServerPort + "/download")
	fmt.Println("Loaded service")

	///create route

	router := mux.NewRouter().StrictSlash(true)
	//router.Host("Localhost")

	router.Handle("/", http.FileServer(http.Dir("../dirmsg")))
	
	router.
		HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {

			if r.Method == http.MethodPost || r.Method == http.MethodGet {

				fmt.Fprintln(w, "http ", 200, `{"msg":"pong"}`)

			} else {

				fmt.Fprintln(w, "http ", 500, "Not authorized / Allowed method POST")
			}
		})
		
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

			} else if r.Method == http.MethodPost {

				fmt.Fprintln(w, "http ", 500, "Not authorized")

			} else {

				fmt.Fprintln(w, "http ", 500, "Not authorized")
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

```

Body of main function  HandleFunc /stop/idEncrypted

```go
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

			} else if r.Method == http.MethodPost {

				fmt.Fprintln(w, "http ", 500, "Not authorized")

			} else {

				fmt.Fprintln(w, "http ", 500, "Not authorized")
			}
		})

```

Body of main function  StopListenAndServe

```go
func StopListenAndServe() {

	fmt.Println("stopping Server File...")
	confServer.Close()
	confServer.Shutdown(nil)
	RemoveFile(socketfile)
}
```

Body of main function  [WriteFileCrypt Encrypting the contents of the file]

```go
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

```

Body of main function  [NewRequestGetStop Encrypts the id of the URL]

```go

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

```

## Examples client

Register user and receive access key 

Using Curl - Sending in json format

```
curl -X POST --data '{"name":"jeff","email":"mail@your.com","password":"321"}' -H "Content-Type:application/json" http://localhost:80/register
```

Using Curl - Access token

```
curl -X GET --data '{"email":"jeff1@gmail.com","password":"321"}' -H "Content-Type:application/json" http://localhost:80/token
```

Uploading with Authorization

```
curl -H 'Authorization:bc8ca54ebabc6f3da724e923fef79238' --form fileupload=@nameFile.bz2 http://localhost:80/upload
```

Uploading with acesskey

```
curl -F 'acesskey:bc8ca54ebabc6f3da724e923fef79238' --form fileupload=@nameFile.bz2 http://localhost:80/upload
```

Download only Authorization

```
curl -H 'Authorization:bc8ca54ebabc6f3da724e923fef79238' -O http://localhost:80/download/nameFile.bz2
```
