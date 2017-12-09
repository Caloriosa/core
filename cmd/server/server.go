package main

import (
	userapi "core/api/user"
	"flag"

	deviceapi "core/api/device"
	"core/pkg/db"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"net/http"
)

var VERSION = "Unknown-build"

func main() {
	flag.Parse()

	glog.Infof("Initializing Caloriosa Core V %s", VERSION)

	var err error

	err = db.NewMongo()
	if err != nil {
		glog.Fatalf("Exitting.")
		return
	}

	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)
	restful.DefaultContainer.EnableContentEncoding(true)
	restful.DefaultContainer.Router(restful.CurlyRouter{})

	userapi.Register(restful.DefaultContainer)
	deviceapi.Register(restful.DefaultContainer)


	err = http.ListenAndServe(":8080", nil)

	glog.Fatal(err)
}
