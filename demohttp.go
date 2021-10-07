package main

/*
1. 接收客户端 request，并将 request 中带的 header 写入 response header
2. 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
3. Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
4. 当访问 localhost/healthz 时，应返回200
*/

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func headerHandler(w http.ResponseWriter, req *http.Request) {
	version := os.Getenv("VERSION")
	headers := req.Header
	for k, v := range headers {
		for _, val := range v {
			w.Header().Add(k, val)
		}
	}
	w.Header().Set("VERSION", version)

}

func logHandler(w http.ResponseWriter, req *http.Request) {
	clientIp := ReadUserIP(req)
	if req.URL.Path == "/" || req.URL.Path == "/healthz" {
		log.Println("Client ip -> " + clientIp + " and responseCode -> " + strconv.Itoa(http.StatusOK))
	} else {
		log.Println("Client ip -> " + clientIp + " and responseCode -> 404")
	}
}

func RootHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		headerHandler(w, req)
		fmt.Fprintln(w, "Welcome GO HttpServer")
	} else {
		http.NotFound(w, req)
	}
	logHandler(w, req)
}

func HealthzHandler(w http.ResponseWriter, req *http.Request) {

	if req.URL.Path == "/healthz" {
		headerHandler(w, req)
		fmt.Fprintln(w, "200")
	} else {
		http.NotFound(w, req)
	}
	logHandler(w, req)
}

func main() {
	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/healthz", HealthzHandler)
	log.Fatal(http.ListenAndServe(":80", nil))
}
