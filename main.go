package main

import (
	"flag"
	"github.com/valyala/fasthttp"
	"log"
)

var addr=flag.String("addr","127.0.0.1:8080","TCP address to listen to for incoming connections")


func uploadHandlerFunc(ctx *fasthttp.RequestCtx){

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

