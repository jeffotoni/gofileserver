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

package gofile

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jeffotoni/gofileserver/pkg/fcrypt"
)

var (
	socketfileTmp = `gofileserver.red`
	Socketfile    = `gofileserver.lock`
	keyCrypt      = `pKv9MQQIDAQABAmEApvlExjvPp0mYs/i`
)

/** [WriteFileCrypt Method that will encrypt the file with the key] */

func WriteFileCrypt(sname string) {

	// get key
	key := []byte(keyCrypt) // 32 bytes

	file, _ := os.Open(sname) // For read access.

	fi, _ := file.Stat()
	data := make([]byte, 16*fi.Size())
	count, _ := file.Read(data)

	file_cry, _ := os.Create(Socketfile)
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

	filename := Socketfile
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
