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
* @package     /gofileupload
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
	"os"

	sfconfig "github.com/jeffotoni/gofileserver/config"
	gcheck "github.com/jeffotoni/gofileserver/pkg/gcheck"
	"github.com/jeffotoni/gofileserver/pkg/gofuplib"
	"github.com/jeffotoni/gofileserver/pkg/postgres/connection"
	_ "github.com/lib/pq"
)

/** [main description] */
func main() {

	// reactive update ukkobox_uploads set up_sent = 1

	// config global

	cfg := sfconfig.GetConfig()

	//connect

	db := connection.Connect()

	// List only those that were not sent

	sqlC := "select up_id,up_user,up_file,up_cloud from ukkobox_uploads where up_sent = 1"
	rows, _ := db.Query(sqlC)

	contador := 0
	for rows.Next() {

		var up_id string
		var up_user string
		var up_file string
		var up_cloud string

		if err := rows.Scan(&up_id, &up_user, &up_file, &up_cloud); err != nil {
			log.Fatal(err)
		}

		pathFile := cfg.Section.PathLocal + up_user + "/" + up_file

		///Updating to progress

		gofuplib.UpdateUploadSent(up_user, up_file, "4")

		if gcheck.FileExists(pathFile) {

			// Analyzing availability - Upload to cloud

			switch gofuplib.AvailabilityCloud() {

			case "Amazon":
				gofuplib.UploadAmazonAwsBucket(pathFile)

			case "Google":
				gofuplib.UploadAmazonAwsBucket(pathFile)

			case "Dropbox":
				gofuplib.UploadAmazonAwsBucket(pathFile)

			case "DigitalOcean":
				gofuplib.UploadAmazonAwsBucket(pathFile)

			}

			// Updated successfully uploaded

			gofuplib.UpdateUploadSent(up_user, up_file, "2")

		} else {

			/// File does not exist

			fmt.Println("file does not exist")
		}

		contador++
	}

	fmt.Println("Sent end ...")
	os.Exit(1)

}
