package httptypes

import (
	"github.com/emicklei/go-restful"
	"net/http"
)

type HttpResponseStatus struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type HttpResponsePack struct {
	Status  HttpResponseStatus `json:"status"`
	Content interface{}        `json:"content"`
}

var EMPTY_CONTENT interface{}

func Send(httpcode int, code, message string, content interface{}, r *restful.Response) error {
	h := HttpResponsePack{
		HttpResponseStatus{code, message},
		content,
	}
	return r.WriteHeaderAndJson(httpcode, h, restful.MIME_JSON)
}

func SendOK(content interface{}, r *restful.Response) error {
	return Send(http.StatusOK, "OK", "ok", content, r)
}

func SendCreated(content interface{}, r *restful.Response) error {
	return Send(http.StatusOK, "CREATED", "created", content, r)
}

func SendDuplicated(content interface{}, r *restful.Response) error {
	return Send(http.StatusConflict, "DUPLICATED", "Duplicated content or resource", content, r)
}

func SendNotFound(content interface{}, r *restful.Response) error {
	return Send(http.StatusNotFound, "NOT_FOUND", "not found", content, r)
}

func SendInvalidData(content interface{}, r *restful.Response) error {
	return Send(http.StatusBadRequest, "INVALID_DATA", "invalid data", content, r)
}

func SendGeneralError(content interface{}, r *restful.Response) error {
	return Send(http.StatusInternalServerError, "DATASOURCE_ERROR", "error fetching data", content, r)
}

func SendBadAuth(r *restful.Response) error {
	return Send(http.StatusUnauthorized, "UNAUTHORIZED", "you're not authoized", nil, r) // TODO correct message
}