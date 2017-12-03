package types

import (
	"github.com/emicklei/go-restful"
	"net/http"
)

type HttpResponseStatus struct {
	Code    int
	Message string
}

type HttpResponsePack struct {
	Status  HttpResponseStatus
	content interface{}
}

var EMPTY_CONTENT interface{}

func Send(code int, message string, content interface{}, r *restful.Response) error {
	h := HttpResponsePack{
		HttpResponseStatus{code, message},
		content,
	}
	return r.WriteHeaderAndJson(code, h, restful.MIME_JSON)
}

func SendOK(content interface{}, r *restful.Response) error {
	return Send(http.StatusOK, "ok", content, r)
}

func SendCreated(content interface{}, r *restful.Response) error {
	return Send(http.StatusCreated, "created", content, r)
}

func SendBadRequest(content interface{}, r *restful.Response) error {
	return Send(http.StatusBadRequest, "bad request", content, r)
}
