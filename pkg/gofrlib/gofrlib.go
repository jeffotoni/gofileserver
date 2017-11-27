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
* project ukkobox
*
* @package     /gofileremove
* @author      jeffotoni
* @copyright   Copyright (c) 2017
* @license     --
* @link        --
* @since       Version 0.1
*/

package gofrlib

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	sfconfig "github.com/jeffotoni/gofileserver/config"
	gcheck "github.com/jeffotoni/gofileserver/pkg/gcheck"
	"github.com/jeffotoni/gofileserver/pkg/postgres/connection"
	_ "github.com/lib/pq"
)

const DefaultDownloadConcurrency = 5

const DefaultDownloadPartSize = 1024 * 1024 * 5

func AwsExistObject(pathKeyFile string) bool {

	cfg := sfconfig.GetConfig()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	if err != nil {
		fmt.Println("fail session aws!", err)
	}

	svc := s3.New(sess)

	result, err := svc.GetObject(&s3.GetObjectInput{

		Bucket: aws.String(cfg.Section.Bucket),
		Key:    aws.String(pathKeyFile),
	})

	if err != nil {

		log.Fatal("Failed to get object", err, result)
		return false

	} else {

		return true
		//fmt.Println(result)
	}
}

func AwsGenerateSignedUrl(pathKeyFile string) string {

	cfg := sfconfig.GetConfig()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	if err != nil {
		fmt.Println("fail session aws!", err)
	}

	svc := s3.New(sess)

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{

		Bucket: aws.String(cfg.Section.Bucket),
		Key:    aws.String(pathKeyFile),
		Body:   strings.NewReader("EXPECTED CONTENTS"),
	})

	UrlAssinada, _ := req.Presign(15 * time.Minute)

	if UrlAssinada != "" {

		return UrlAssinada

	} else {

		return ""
	}
}

func FileExistCloudServer(cloud string, pathFileCloud string, pathKeyFile string) bool {

	switch cloud {

	case "Amazon":
		return AmazonExistObject(pathKeyFile)

	case "Google":
		return false

	case "Dropbox":
		return false

	case "DigitalOcean":
		return false

	default:
		return false
	}

	return false
}

func AmazonExistObject(pathKeyFile string) bool {

	if AwsExistObject(pathKeyFile) {

		///file exist
		return true

	} else {

		return false
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
		//afetUp
		//fmt.Println("Update made: ", fmt.Sprintf("%v", afetUp))

	} else {

		fmt.Println("There is no record in db")
	}
}
