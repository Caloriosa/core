package main

import (
	userapi "core/api/user"
	"flag"

	"core/api/auth"
	deviceapi "core/api/device"
	"core/pkg/config"
	"core/pkg/db"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"net/http"
)

var VERSION = "Unknown-build"

func main() {
	configfile := flag.String("config", "config.yaml", "Path to the config yaml file")
	flag.Parse()

	glog.Infof("Initializing Caloriosa Core V %s", VERSION)

	var err error
	err = config.LoadConfig(*configfile)
	if err != nil {
		glog.Fatal("Error loading config file: ", err)
	}

	glog.Info("Loaded config file: ", config.LoadedConfig)

	err = db.NewMongo()
	if err != nil {
		glog.Fatalf("Exitting.")
		return
	}

	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)
	restful.DefaultContainer.EnableContentEncoding(true)
	restful.DefaultContainer.Router(restful.CurlyRouter{})
	restful.Filter(restful.OPTIONSFilter())
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-My-Header"},
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST"},
		CookiesAllowed: false,
		Container:      restful.DefaultContainer}
	restful.DefaultContainer.Filter(cors.Filter)

	userapi.Register(restful.DefaultContainer)
	deviceapi.Register(restful.DefaultContainer)
	auth.Register(restful.DefaultContainer)

	err = http.ListenAndServe(":8080", nil)

	glog.Fatal(err)
}
