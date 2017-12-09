package tools

import (
	"core/pkg/db"
	"core/types"
	"fmt"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const DEFAULT_LIMIT = 50

func GetPagination(values url.Values) (int, int, error) {
	var page, limit int
	var err error

	if val, ok := values["page"]; ok {
		page, err = strconv.Atoi(val[0])
		if err != nil {
			return 0, 0, err
		}
	} else {
		page = 0
	}

	if val, ok := values["limit"]; ok {
		limit, err = strconv.Atoi(val[0])
		if err != nil {
			return 0, 0, nil
		}
	} else {
		limit = DEFAULT_LIMIT
	}

	return page, limit, nil
}

func GetFilters(values url.Values, mytype interface{}) (map[string]interface{}, error) {
	filtersObj := bson.M{}
	var err error
	filterLen := len("filter[")

	for key, val := range values {
		if strings.HasPrefix(key, "filter[") {
			xkey := key[filterLen : len(key)-1]
			xvalue := reflect.ValueOf(mytype)
			var mykind reflect.Kind
			for i := 0; i < reflect.Indirect(xvalue).NumField(); i++ {
				tag := string(reflect.TypeOf(mytype).Elem().Field(i).Tag.Get("json"))
				if strings.Split(tag, ",")[0] == xkey {
					mykind = reflect.TypeOf(mytype).Elem().Field(i).Type.Kind()
				}
			}

			if mykind == reflect.Int {
				filtersObj[xkey], err = strconv.Atoi(val[0])
				if err != nil {
					return nil, err
				}
			} else if mykind == reflect.Bool {
				if val[0] == "true" {
					filtersObj[xkey] = true
				} else {
					filtersObj[xkey] = false
				}
			} else if mykind == reflect.String {
				filtersObj[xkey] = bson.M{"$regex": bson.RegEx{Pattern: fmt.Sprintf("^%s$", val[0]), Options: "i"}}
			}

			glog.Info("Adding %s = %s to filters", xkey, val)
		}
	}

	return filtersObj, nil
}

func GetToken(req *http.Request) *types.Token {
	mytoken := []types.Token{}
	glog.Info("Headers: ", req.Header)
	if tokens, ok := req.Header["Authorization"]; ok {
		glog.Info("Found tokens: ", tokens)
		if len(tokens) > 0 {
			glog.Info("X token: ", tokens[0], "token: ", mytoken)
			xtoken := strings.TrimPrefix(tokens[0], "Bearer ")
			if err := db.MONGO.Get(COLLECTION_TOKENS, &mytoken, bson.M{"token": xtoken}, 0, 1); err == nil {
				glog.Info("Extracted token: ", mytoken)
				if mytoken[0].ExpireAt.Before(time.Now().UTC()) {
					return &mytoken[0]
				} else {
					db.MONGO.Connection.Collection(COLLECTION_TOKENS).DeleteDocument(&mytoken[0])
					glog.Info("Deleting expired token.")
				}
			}
		}
	}

	return nil
}
