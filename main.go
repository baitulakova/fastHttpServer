package main

import (
	"flag"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"io"
)

var addr=flag.String("addr","127.0.0.1:8080","TCP address to listen to for incoming connections")

func createStorage() string {
	userHome:=os.Getenv("HOME")
	fileStorage := userHome+"/fasthttpServerStorage/"
	err:=os.MkdirAll(fileStorage,os.ModePerm)
	if err!=nil{
		log.Fatal(err)
	}
	return fileStorage
}

var storage string

func uploadHandlerFunc(ctx *fasthttp.RequestCtx){
	if string(ctx.Method())=="POST"{
		file,err:=ctx.FormFile("file")
		if err!=nil{
			ctx.SetStatusCode(400)
			log.Fatal(err)
		}
		err=fasthttp.SaveMultipartFile(file,storage+file.Filename)
		if err!=nil{
			log.Fatal(err)
		}
		ctx.SetStatusCode(200)
		log.Println("Uploaded ",file.Filename," file")
	}
}

func downloadHandlerFunc(ctx *fasthttp.RequestCtx){
	file:=ctx.URI().QueryArgs()
	filename:=file.Peek("filename")
	if string(filename)==""{
		ctx.SetStatusCode(404)
		log.Fatal("Error in upload file. File not found")
	}
	f,err:=os.Open(storage+string(filename))
	if err!=nil{
		ctx.SetStatusCode(404)
		log.Fatal("Error opening file: ",err)
	}
	defer f.Close()
	io.Copy(ctx,f)
	ctx.SetStatusCode(200)
}

func handlerFunc(ctx *fasthttp.RequestCtx){

}

func main(){
	storage=createStorage()

	flag.Parse()

	m:=func(ctx *fasthttp.RequestCtx){
		switch string(ctx.Path()) {
		case "/upload":
			uploadHandlerFunc(ctx)
		case "/download":
			downloadHandlerFunc(ctx)
		case "/":
			handlerFunc(ctx)
		default:
			ctx.Error("not found",fasthttp.StatusNotFound)
		}
	}

	server:=fasthttp.Server{
		Handler:m,
		DisableKeepalive: false,
	}

	log.Println("Server listening at " , *addr)
	err:=server.ListenAndServe(*addr)
	if err!=nil{
		log.Fatal(err)
	}
}

