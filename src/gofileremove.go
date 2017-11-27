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

package main

import (
	"fmt"
	"log"
	"os"

	sfconfig "github.com/jeffotoni/gofileserver/config"
	gcheck "github.com/jeffotoni/gofileserver/pkg/gcheck"
	"github.com/jeffotoni/gofileserver/pkg/gofrlib"
	"github.com/jeffotoni/gofileserver/pkg/postgres/connection"
	_ "github.com/lib/pq"
)

/** [main description] */
func main() {

	//config global

	cfg := sfconfig.GetConfig()

	// reactive update ukkobox_uploads set up_sent = 1

	//connect
	db := connection.Connect()

	// List only those that have already been sent

	sqlC := "select up_id,up_user,up_file,up_cloud from ukkobox_uploads where up_sent = 2"

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

		// Lock to not occur at the same time removal

		gofrlib.UpdateUploadSent(up_user, up_file, "6")

		//clouds

		clouds := make(map[string]string)

		// Clouds paths

		clouds["Amazon"] = "https://s3.amazonaws.com"
		clouds["Google"] = ""
		clouds["DigitalOcean"] = ""
		clouds["Dropbox"] = ""

		// Cloud path

		pathCloud := clouds[up_cloud]

		pathFileCloud := pathCloud + cfg.Section.PathLocal + up_user + "/" + up_file
		pathFileLocal := cfg.Section.PathLocal + up_user + "/" + up_file

		if gcheck.FileExists(pathFileLocal) {

			// If there is in the cloud delete locally

			if gofrlib.FileExistCloudServer(up_cloud, pathFileCloud, pathFileLocal) {

				// remove local

				fmt.Println("remove file ", pathFileLocal)

				var err = os.Remove(pathFileLocal)
				gcheck.SError(err)

				// File exists in the cloud, and has been removed from the

				gofrlib.UpdateUploadSent(up_user, up_file, "5")

			} else {

				/// File does not exist
				fmt.Println("File does not exist in cloud")
			}
		} else {

			/// File does not exist
			fmt.Println("file does not exist local")
		}

		contador++
	}

	fmt.Println("Files successfully removed ...")

}
