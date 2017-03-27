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
* @package     libfileserver
* @author      jeffotoni
* @copyright   Copyright (c) 2017
* @license     --
* @link        --
* @since       Version 0.1
*/

package gofuplib

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	sfconfig "github.com/jeffotoni/gofileserver/config"
	gcheck "github.com/jeffotoni/gofileserver/pkg/gcheck"
	"github.com/jeffotoni/gofileserver/pkg/postgres/connection"
	_ "github.com/lib/pq"
)

//const DefaultDownloadConcurrency = 5
//const DefaultDownloadPartSize = 1024 * 1024 * 5

const MinUploadPartSize int64 = 1024 * 1024 * 5

const DefaultUploadConcurrency = 5

const DefaultUploadPartSize = MinUploadPartSize

func AvailabilityCloud() string {

	/// Amazon 1
	/// Google 2

	clouds := []string{"Amazon", "Google", "DigitalOcean", "Dropbox"}

	///Create an algorithm to make decision and choose the server
	//// random

	var max, min int
	max = 4
	min = 1

	rand.Seed(time.Now().Unix())
	rand.Intn(max - min)
	position := rand.Intn(max-min) + min

	return clouds[position]

}

func UploadAmazonAwsBucket(pathFile string) {

	if pathFile != "" {

		//config global

		cfg := sfconfig.GetConfig()

		//Open file local

		file, err := os.Open(pathFile)

		if err != nil {
			log.Fatalln(err)
		}

		defer file.Close()

		awconfig := &aws.Config{
			Region: aws.String("us-east-1"),
		}

		sess := session.New(awconfig)
		uploader := s3manager.NewUploader(sess)

		bucketName := cfg.Section.Bucket
		keyName := pathFile

		paramUp := &s3manager.UploadInput{
			Bucket: &bucketName,
			Key:    &keyName,
			Body:   file,
		}

		//show log

		fmt.Println("Uploading: " + pathFile)

		// Perform an upload.

		result, err := uploader.Upload(paramUp, func(u *s3manager.Uploader) {

			u.PartSize = 10 * 1024 * 1024 // 10MB part size
			u.LeavePartsOnError = true    // Don't delete the parts if the upload fails.
			fmt.Println("200 ok")
		})

		if err != nil {

			if multierr, ok := err.(s3manager.MultiUploadFailure); ok {
				// Process error and its associated uploadID
				fmt.Println("Error multier:", multierr.Code(), multierr.Message(), multierr.UploadID())

			} else {

				// Process error generically
				fmt.Println("Error generic :", err.Error())
			}
		}

		log.Println(result)

	} else {

		fmt.Println("Error Path is required")
	}
}

func UpdateUploadSent(up_user string, up_file string, up_sent string) {

	// 1 = local, 2 = sent cloud, 3 = Send error, 4 = in progress, 5 = Removed locally, 6 = Lock to run, progress

	//connect

	db := connection.Connect()

	var up_id int64
	sqlCount := "select up_id as up_id from ukkobox_uploads where up_user = $1 AND up_file = $2"
	rows := db.QueryRow(sqlCount, up_user, up_file)
	rows.Scan(&up_id)

	if up_id > 0 {

		//Returns the status again to send
		UpdateSql := "UPDATE ukkobox_uploads set up_sent = " + up_sent + " WHERE up_id = $1"
		update, err := db.Prepare(UpdateSql)
		resp, err := update.Exec(up_id)
		gcheck.SError(err)

		resp.RowsAffected()

		//fmt.Println("Update made: ", fmt.Sprintf("%v", afetUp))

	} else {

		fmt.Println("There is no record in db")
	}
}
