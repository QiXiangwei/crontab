package master

import (
	"net"
	"net/http"
)

type ApiServer struct {
	HttpServer *http.Server
}

var (
	GApiServer *ApiServer
)

func handleJobSave(w http.ResponseWriter, r *http.Request) {

}

func InitApiServer() (err error) {
	var (
		httpServer *http.Server
		mux        *http.ServeMux
		listen     net.Listener
	)

	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)

	if listen, err = net.Listen("tcp", ":8090"); err != nil {
		return
	}

	GApiServer = &ApiServer{HttpServer:httpServer}

	go httpServer.Serve(listen)

	return
}
