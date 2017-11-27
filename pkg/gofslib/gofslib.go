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
* @package     libfileserver
* @author      jeffotoni
* @copyright   Copyright (c) 2017
* @license     --
* @link        --
* @since       Version 0.1
*/

package gofslib

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
	sfconfig "github.com/jeffotoni/gofileserver/config"
	"github.com/jeffotoni/gofileserver/pkg/fcrypt"
	"github.com/jeffotoni/gofileserver/pkg/gcheck"
	"github.com/jeffotoni/gofileserver/pkg/postgres/connection"
)

func DownloadFile(w http.ResponseWriter, r *http.Request) {

	cfg := sfconfig.GetConfig()

	autorization := r.Header.Get("Authorization")

	///check bd acess key

	if autorization == "" {

		fmt.Fprintln(w, "http ", 500, "Not authorized")

	} else {

		if ValidadeAccess(autorization) {

			vars := mux.Vars(r)
			nameFile := vars["name"]

			if nameFile != "" {

				acessekey := autorization

				pathFile := cfg.Section.PathLocal + "/" + acessekey + "/" + nameFile

				if gcheck.FileExists(pathFile) {

					///show msg server

					fmt.Println("file: ", pathFile)

					///out file

					http.ServeFile(w, r, pathFile)

					//send client msg

					fmt.Fprintln(w, "", 200, "OK")

					//Generating logs
					InsertLogDownloads(acessekey, nameFile, "Amazon")

				} else {

					/////////
					// 1 = local, 2 = sent cloud, 3 = Send error, 4 = in progress, 5 = Removed locally, 6 = Lock to run, progress

					GetStatusCloud := GetStatusCloudServer(acessekey)

					fmt.Println("status cloud: ", GetStatusCloud)

					if GetStatusCloud == 1 {

						fmt.Println("File does not exist")
						fmt.Fprintln(w, "", 500, "File does not exist")

					} else if GetStatusCloud == 2 || GetStatusCloud == 5 {

						// Amazon
						// Google

						GetNameCloud := GetNameCloudServer(acessekey)

						if GetNameCloud == "Amazon" {

							//fmt.Println("Cloud download Amazon")
							AmazonDownload(pathFile, w, r)

						} else if GetNameCloud == "Google" {

							fmt.Println("Cloud download Google")

						} else if GetNameCloud == "DigitalOcean" {

							fmt.Println("Cloud download DigitalOcean")

						} else if GetNameCloud == "Dropbox" {

							fmt.Println("Cloud download Dropbox")
						}

					} else if GetStatusCloud == 3 {

						fmt.Println("Error sending to cloud")
						fmt.Fprintln(w, "", 500, "Error sending to cloud")
					}
				}

			} else {

				fmt.Fprintln(w, "", 500, "Required file name")

			}

		} else {

			fmt.Fprintln(w, "", 500, "access denied")
		}
	}
}

func AmazonDownload(pathFileLocal string, w http.ResponseWriter, r *http.Request) {

	cfg := sfconfig.GetConfig()

	file, err := os.Create(pathFileLocal)

	if err != nil {
		log.Fatalln(err)
	}

	defer file.Close()

	config := &aws.Config{
		Region: aws.String("us-east-1"),
	}

	sess := session.New(config)

	downloader := s3manager.NewDownloader(sess)

	fmt.Println("Start Download...", pathFileLocal)

	paramDown := &s3.GetObjectInput{

		Bucket: aws.String(cfg.Section.Bucket),
		Key:    aws.String(pathFileLocal),
	}

	// Perform an upload.

	result, err := downloader.Download(file, paramDown)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Sprint("%s", result)

	// Sends file to client

	http.ServeFile(w, r, pathFileLocal)

	// remove file disk local

	time.Sleep(time.Second * 1)
	fmt.Println("Remove file:", pathFileLocal)

	errx := os.Remove(pathFileLocal)
	gcheck.SError(errx)

}

func InsertLogDownloads(acessekey string, nameFile string, cloud string) {

	//connect
	db := connection.Connect()

	// Save log downloads

	insert := "INSERT INTO ukkobox_download (up_user,up_file,up_cloud) VALUES ($1,$2,$3)"
	stmt, err := db.Prepare(insert)
	afetIns, err := stmt.Exec(acessekey, nameFile, cloud)
	gcheck.SError(err)

	fmt.Println(fmt.Sprintf("%v", afetIns))
}

func UploadFileEasy(w http.ResponseWriter, r *http.Request) {

	cfg := sfconfig.GetConfig()

	autorization := r.Header.Get("Authorization")

	//fmt.Println("Token: ", autorization)

	if autorization == "" {

		acessekey := r.FormValue("acesskey")
		autorization = acessekey
	}

	if autorization == "" {

		fmt.Fprintln(w, "", 500, "Not Authorized")

	} else {

		////check database get id user

		if ValidadeAccess(autorization) {

			///Valid user
			acessekey := autorization

			//fmt.Printf("Acess Key:%s\n", acessekey)

			sizeMaxUpload := r.ContentLength / 1048576 ///Mb

			if sizeMaxUpload > cfg.Section.UploadSize {

				fmt.Println("The maximum upload size: ", cfg.Section.UploadSize, "Mb is large: ", sizeMaxUpload, "Mb", " in bytes: ", r.ContentLength)
				fmt.Fprintln(w, "", 500, "Unsupported file size max: ", cfg.Section.UploadSize, "Mb")

			} else {

				// field upload

				file, handler, errf := r.FormFile("fileupload")

				//fmt.Println("error: ", errf)

				if errf != nil {
					color.Red("Error big file, try again!")
					http.Error(w, "Error parsing uploaded file: "+errf.Error(), http.StatusBadRequest)
					return
				}

				defer file.Close()

				///create dir to key
				pathUpKeyUser := cfg.Section.PathLocal + "/" + acessekey

				existPath, _ := os.Stat(pathUpKeyUser)

				if existPath == nil {

					// create path
					os.MkdirAll(pathUpKeyUser, 0777)
				}

				pathUserAcess := cfg.Section.PathLocal + "/" + acessekey + "/" + handler.Filename

				// copy file and write

				f, _ := os.OpenFile(pathUserAcess, os.O_WRONLY|os.O_CREATE, 0777)
				defer f.Close()
				n, _ := io.Copy(f, file)

				up_size := fmt.Sprintf("%v", r.ContentLength)

				//INSERT / UPDATE DB

				InsertUpdateUploadFileDB(acessekey, handler.Filename, "Amazon", up_size)

				//To display results on server

				name := strings.Split(handler.Filename, ".")
				fmt.Printf("File name: %s\n", name[0])
				fmt.Printf("extension: %s\n", name[1])

				fmt.Println("size file: ", sizeMaxUpload)
				fmt.Println("allowed: ", cfg.Section.UploadSize)

				fmt.Printf("copied: %v bytes\n", n)
				fmt.Printf("copied: %v Kb\n", n/1024)
				fmt.Printf("copied: %v Mb\n", n/1048576)
				fmt.Fprintln(w, "", 200, "OK")

			}

		} else {

			fmt.Fprintln(w, "", 500, "access denied")
		}
	}
}

func InsertUpdateUploadFileDB(up_user string, up_file string, up_cloud string, up_size string) {

	//connect
	db := connection.Connect()

	var up_id int64
	sqlCount := "select up_id as up_id from ukkobox_uploads where up_user = $1 AND up_file = $2"
	rows := db.QueryRow(sqlCount, up_user, up_file)
	rows.Scan(&up_id)

	if up_id > 0 {

		t := time.Now()

		up_update := t.Format("2006-01-02")
		up_updateh := t.Format("15:04:05")

		//Returns the status again to send
		UpdateSql := "UPDATE ukkobox_uploads set up_size = '" + up_size + "', up_update = '" + up_update + "', up_updateh = '" + up_updateh + "', up_sent = 1 WHERE up_id=$1"
		update, err := db.Prepare(UpdateSql)
		resp, err := update.Exec(up_id)
		gcheck.SError(err)

		afetUp, _ := resp.RowsAffected()

		fmt.Println(fmt.Sprintf("%v", afetUp))

	} else {

		///Save to file database

		insert := "INSERT INTO ukkobox_uploads(up_user,up_file,up_cloud,up_size) VALUES ($1,$2,$3,$4)"
		stmt, err := db.Prepare(insert)
		afetIns, err := stmt.Exec(up_user, up_file, up_cloud, up_size)
		gcheck.SError(err)

		fmt.Println(fmt.Sprintf("%v", afetIns))
	}
}

type DadaJson struct {
	User     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func RegisterUserJson(w http.ResponseWriter, r *http.Request) {

	objJson := DadaJson{}
	errj := json.NewDecoder(r.Body).Decode(&objJson)

	if errj != nil {

		fmt.Println(errj)
	}

	name := objJson.User
	email := objJson.Email
	password := objJson.Password

	if name != "" && email != "" && password != "" {

		//fmt.Println(name)
		fmt.Println(email)
		//fmt.Println(password)

		///if user exist error
		msg := UserExist(email)

		if msg == "ok" {

			///Write to the database and return a token

			token := UserRegister(name, email, password)

			fmt.Fprintln(w, token)

		} else {

			//To send error message to client

			fmt.Fprintln(w, "Error", 500, msg)

		}

	} else {

		fmt.Fprintln(w, "Error", 500, " Some field is empty!")
	}

}

func UserRegister(name string, email string, password string) string {

	db := connection.Connect()

	var user_name string
	var user_email string
	var user_pass string
	var user_token string
	var passmd5 string

	user_name = strings.ToLower(name)
	user_email = strings.ToLower(email)
	user_pass = password

	///generate token

	time := strconv.FormatInt(time.Now().Unix(), 10)

	tempoExecucao := fmt.Sprintf("%s", time)

	tempoExecucao = tempoExecucao + email + password + name

	//crypt md5 token

	user_token = fcrypt.Md5(tempoExecucao)

	//crypt md5 passwords

	passmd5 = fcrypt.Md5(user_pass)

	var count int64
	sqlCount := "select count(user_id) as count from ukkobox_user where lower(user_email) = $1"
	rows := db.QueryRow(sqlCount, user_email)
	rows.Scan(&count)

	if count > 0 {

		///user exist
		return "user exist!"

	} else {

		insert := "INSERT INTO ukkobox_user(user_name,user_email,user_pass,user_token) VALUES ($1,$2,$3,$4)"

		stmt, err := db.Prepare(insert)

		_, err = stmt.Exec(user_name, user_email, passmd5, user_token)

		gcheck.SError(err)

		return user_token
	}
}

func UserExist(email string) string {

	if email != "" {

		db := connection.Connect()

		emailLower := strings.TrimSpace(email)
		emailLower = strings.Trim(emailLower, " ")
		emailLower = strings.ToLower(emailLower)

		//Verify valid email

		if gcheck.SCheckFormatEmail(emailLower) == nil {

			var count int64
			sqlCount := "select count(user_id) as count from ukkobox_user where lower(user_email) = $1"
			rows := db.QueryRow(sqlCount, emailLower)
			rows.Scan(&count)

			if count > 0 {

				///user exist
				return "user exist"

			} else {

				return "ok"
			}

		} else {

			return "Error invalid email"
		}

	} else {

		return "Error Required email"
	}
}

func GetTokenUser(w http.ResponseWriter, r *http.Request) {

	objJson := DadaJson{}
	errj := json.NewDecoder(r.Body).Decode(&objJson)

	if errj != nil {

		fmt.Println(errj)
	}

	email := objJson.Email
	password := objJson.Password

	emailLower := strings.TrimSpace(email)
	emailLower = strings.Trim(emailLower, " ")
	emailLower = strings.ToLower(emailLower)

	if emailLower != "" && password != "" {

		passmd5 := fcrypt.Md5(password)

		db := connection.Connect()

		var token string
		sqlCount := "select user_token as token from ukkobox_user where lower(user_email) = $1 AND user_pass = $2"

		rows := db.QueryRow(sqlCount, emailLower, passmd5)
		rows.Scan(&token)

		if token != "" {

			fmt.Fprintln(w, "Access token: ", token)

		} else {

			fmt.Fprintln(w, "Access denied")
		}

	} else {

		fmt.Fprintln(w, "Error", 500, " email e password required")
	}

}

func ValidadeAccess(token string) bool {

	if token != "" {

		db := connection.Connect()

		tokenLower := strings.TrimSpace(token)
		tokenLower = strings.Trim(tokenLower, " ")

		var count int64
		sqlCount := "select count(user_id) as count from ukkobox_user where user_token = $1 AND user_active = $2"
		rows := db.QueryRow(sqlCount, tokenLower, 1)
		rows.Scan(&count)

		if count > 0 {

			///user exist
			return true

		} else {

			return false
		}

	} else {

		return false
	}
}

func GetStatusCloudServer(token string) int {

	/// up_sent = {}
	/// 1 = local
	/// 2 = cloud
	/// 3 = error submitting

	if token != "" {

		db := connection.Connect()

		tokenLower := strings.TrimSpace(token)
		tokenLower = strings.Trim(tokenLower, " ")

		var up_sent int
		sqlCount := "select up_sent from ukkobox_uploads where up_user = $1"
		rows := db.QueryRow(sqlCount, token)
		rows.Scan(&up_sent)

		if up_sent > 0 {

			///user exist
			return up_sent

		} else {

			return 0
		}

	} else {

		return 0
	}
}

func GetNameCloudServer(token string) string {

	// Amazon
	// Google

	if token != "" {

		db := connection.Connect()

		tokenLower := strings.TrimSpace(token)
		tokenLower = strings.Trim(tokenLower, " ")

		var up_cloud string
		sqlCount := "select up_cloud from ukkobox_uploads where up_user = $1"
		rows := db.QueryRow(sqlCount, token)
		rows.Scan(&up_cloud)

		if up_cloud != "" {

			///user exist
			return up_cloud

		} else {

			return ""
		}

	} else {

		return ""
	}
}
