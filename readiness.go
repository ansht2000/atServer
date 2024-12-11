package main

import "net/http"

func handlerReadiness(resWriter http.ResponseWriter, req *http.Request) {
	resWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resWriter.WriteHeader(200)
	resWriter.Write([]byte("OK"))
}