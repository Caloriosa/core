package main

import (
	userapi "core/api/user"
	"flag"

	deviceapi "core/api/device"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"net/http"
	"core/pkg/db"
)

var VERSION = "Unknown-build"



func main() {
	flag.Parse()

	glog.Infof("Initializing Caloriosa Core V %s", VERSION)

	db.ConnectMongo()

	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)
	restful.DefaultContainer.EnableContentEncoding(true)
	restful.DefaultContainer.Router(restful.CurlyRouter{})

	userapi.Register(restful.DefaultContainer)
	deviceapi.Register(restful.DefaultContainer)

	glog.Fatal(http.ListenAndServe(":8080", nil))
}
