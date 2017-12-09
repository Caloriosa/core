package tools

import (
	"net/url"
	"strconv"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"github.com/golang/glog"
	"reflect"
	"net/http"
)

const DEFAULT_LIMIT = 50

func GetPagination(values url.Values) (int, int, error) {
	var page, limit int
	var err error

	if val, ok := values["page"]; ok  {
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

func GetFilters(values url.Values, mytype interface{}) (bson.M, error) {
	filtersObj := bson.M{}
	var err error
	filterLen := len("filter[")

	for key, val := range values {
		if strings.HasPrefix(key, "filter[") {
			xkey := key[filterLen:len(key)-1]
			xvalue := reflect.ValueOf(mytype)
			var mykind reflect.Kind
			for i := 0; i < reflect.Indirect(xvalue).NumField(); i++ {
				tag := string(reflect.TypeOf(mytype).Elem().Field(i).Tag.Get("json"))
				if strings.Split(tag,",")[0] == xkey {
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
				filtersObj[xkey] = val[0]
			}

			glog.Info("Adding %s = %s to filters", xkey, val)
		}
	}

	return filtersObj, nil
}

func GetToken(req *http.Request) string {
	if tokens, ok := req.Header["Authorization"]; ok {
		if len(tokens) > 0 {
			return strings.TrimPrefix(tokens[0], "Bearer ")
		}
	}

	return ""
}