package main

import (
	"flag"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"fmt"
)

var addr=flag.String("addr","127.0.0.1:8080","TCP address to listen to for incoming connections")

func createStorage() (path string){
	userHome:=os.Getenv("HOME")
	fileStorage := userHome+"/fasthttpServerStorage/"
	err:=os.MkdirAll(fileStorage,os.ModePerm)
	if err!=nil{
		fmt.Println("error",err)
	}
	return fileStorage
}

func uploadHandlerFunc(ctx *fasthttp.RequestCtx){
	if string(ctx.Method())=="POST"{
		file,err:=ctx.FormFile("file")
		if err!=nil{
			ctx.SetStatusCode(400)
			log.Fatal(err)
		}
		err=fasthttp.SaveMultipartFile(file,createStorage()+file.Filename)
		if err!=nil{
			log.Fatal(err)
		}
		log.Println("Uploaded ",file.Filename," file")
	}
}

func downloadHandlerFunc(ctx *fasthttp.RequestCtx){

}

func handlerFunc(ctx *fasthttp.RequestCtx){

}

func main(){
	flag.Parse()

	m:=func(ctx *fasthttp.RequestCtx){
		switch string(ctx.Path()) {
		case "/upload":
			uploadHandlerFunc(ctx)
		case "download":
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

