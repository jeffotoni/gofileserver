# gofileserver
Simple RESTFUL API demo server written on Go (golang) user Aws S3, postgresql
Authorization and authentication with token
Api designed to upload, download files to your local server and then download 
to cloud servers, the first to be implemented was aws s3
Call authorization required, Upload, Download, register new user

* It implements the `Upload files` interface so it is compatible with the standard `http.ServeMux`.
* It Implements the Download files based on URL host


## Used libraries:
- https://github.com/lib/pq - Sql Database.
- https://gopkg.in/gcfg.v1 - Text-based configuration files with "name=value" pairs grouped into sections (gcfg files).
- https://github.com/aws/aws-sdk-go/aws - Package sdk Aws Amazon
- https://github.com/gorilla/mux - Implements a request router and dispatcher for matching incoming requests

---

* [Install](#install)
* [Examples Client](#examples-client)
* [Register User](#register-user)
* [Upload File](#upload-files)
* [Download Files](#download-files)

## Install

With a [correctly configured](https://golang.org/doc/install#testing) Go toolchain:

```sh
go get -u github.com/jeffotoni/gofileserver
```

## Examples client

Register user and receive access key 

Using Curl - Sending in json format

```sh
curl -X POST --data '{"name":"jeff","email":"mail@your.com","password":"321"}' -H "Content-Type:application/json" http://localhost:4001/register
```

Using Curl - Sending in form format

```sh
curl -X POST --data 'name=jefferson&email=yourmail@yes.com&password=342' http://localhost:4001/registe
```
