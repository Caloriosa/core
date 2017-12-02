package main

import (
	"flag"

	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"io"
	"net/http"
	userapi "core/api/user"
)

var VERSION = "Unknown-build"

func main() {
	flag.Parse()

	glog.Infof("Initializing Caloriosa Core V %s", VERSION)

	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)
	restful.DefaultContainer.EnableContentEncoding(true)
	restful.DefaultContainer.Router(restful.CurlyRouter{})

	userapi.Register(restful.DefaultContainer)

	glog.Fatal(http.ListenAndServe(":8080", nil))
}

func hello(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "world")
}