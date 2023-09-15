# go-fileserver

Simple HTTP server that serves a folder

## Running locally

```shell
PORT=8080 SERVE_FROM_FOLDER=content go run main.go
```

## Running with Docker

```shell
docker build --tag go-fileserver .
docker run -p 8080:8080 go-fileserver
```

## Usage
```shell
# Uploading a file
curl --data-binary "@avatar.jpg" http://localhost:8080/images/avatar.jpg

# Getting a file
curl http://localhost:8080/images/avatar.jpg -o avatar2.jpg

# Deleting a file
curl -X DELETE http://localhost:8080/images/avatar.jpg
```
