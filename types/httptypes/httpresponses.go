package httptypes

import (
	"github.com/emicklei/go-restful"
	"net/http"
)

type HttpResponseStatus struct {
	HttpCode int    `json:"-"`
	Code     string `json:"code"`
	Message  string `json:"message"`
}

type HttpResponsePack struct {
	Status  HttpResponseStatus `json:"status"`
	Content interface{}        `json:"content"`
}

var EMPTY_CONTENT = map[string]string{}

func Send(httpcode int, code, message string, content interface{}, r *restful.Response) error {
	if content == nil {
		content = EMPTY_CONTENT // apparently nil is bad, let's send empty {} brackets
	}
	h := HttpResponsePack{
		HttpResponseStatus{Code: code, Message: message},
		content,
	}
	return r.WriteHeaderAndJson(httpcode, h, restful.MIME_JSON)
}

func SendResponse(r *restful.Response, resp *HttpResponseStatus, content interface{}) error {
	//return r.WriteHeaderAndJson(resp.HttpCode, HttpResponsePack{Status: *resp, Content: content}, restful.MIME_JSON)
	return Send(resp.HttpCode, resp.Code, resp.Message, content, r)
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

var HTTP_RESPONSE_OK = HttpResponseStatus{HttpCode: http.StatusOK, Code: "OK", Message: "OK"}
var DUPLICATED = HttpResponseStatus{HttpCode: http.StatusConflict, Code: "DUPLICATED", Message: "Duplicated content or resource"}
var NOT_FOUND = HttpResponseStatus{HttpCode: http.StatusNotFound, Code: "NOT_FOUND", Message: "Resource not found"}
var DATASOURCE_ERROR = HttpResponseStatus{HttpCode: http.StatusInternalServerError, Code: "DATASOURCE_ERROR", Message: "An error occurred with datasource (database)"}
var PERMISSION_DENIED = HttpResponseStatus{HttpCode: http.StatusForbidden, Code: "PERMISSION_DENIED", Message: "Permission denied"}
var UNAUTHORIZED = HttpResponseStatus{HttpCode: http.StatusUnauthorized, Code: "UNAUTHORIZED", Message: "Unauthorized request"}
var UNAVAILABLE = HttpResponseStatus{HttpCode: http.StatusGone, Code: "UNAVAILABLE", Message: "Requested resource or content is not available"}
var REMOVE_FAILED = HttpResponseStatus{HttpCode: http.StatusInternalServerError, Code: "REMOVE_FAILED", Message: "Resource or content remove failed"}
var AUTH_FAILED = HttpResponseStatus{HttpCode: http.StatusInternalServerError, Code: "AUTH_FAILED", Message: "Authentication failed"}
var INVALID_DATA = HttpResponseStatus{HttpCode: http.StatusBadRequest, Code: "INVALID_DATA", Message: "Invalid data"}
var INVALID_CREDENTIALS = HttpResponseStatus{HttpCode: http.StatusUnauthorized, Code: "INVALID_CREDENTIALS", Message: "Invalid credentials (login, password)"}
var INVALID_SENSOR = HttpResponseStatus{HttpCode: http.StatusBadRequest, Code: "INVALID_SENSOR", Message: "Sensor(s) is invalid or not exists"}
var INVALID_TOKEN = HttpResponseStatus{HttpCode: http.StatusBadRequest, Code: "INVALID_TOKEN", Message: "Your token is not valid"}
var TOKEN_EXPIRED = HttpResponseStatus{HttpCode: http.StatusUnauthorized, Code: "TOKEN_EXPIRED", Message: "Token expired. Please re-login"}
var USER_EXISTS = HttpResponseStatus{HttpCode: http.StatusConflict, Code: "USER_EXISTS", Message: "User %s exists"}
var WEAK_PASSWORD = HttpResponseStatus{HttpCode: http.StatusBadRequest, Code: "WEAK_PASSWORD", Message: "Chosen password is too weak"}
var INVALID_PASSWORD = HttpResponseStatus{HttpCode: http.StatusBadRequest, Code: "INVALID_PASSWORD", Message: "Chosen password is not valid"}
var INVALID_USERNAME = HttpResponseStatus{HttpCode: http.StatusBadRequest, Code: "INVALID_USERNAME", Message: "Chosen username is not valid"}
var INVALID_EMAIL = HttpResponseStatus{HttpCode: http.StatusBadRequest, Code: "INVALID_EMAIL", Message: "Your email is not valid"}
var PASSWORD_MISMATCH = HttpResponseStatus{HttpCode: http.StatusBadRequest, Code: "PASSWORD_MISMATCH", Message: "Passwords is not match"}
var ACTIVATION_FAILED = HttpResponseStatus{HttpCode: http.StatusInternalServerError, Code: "ACTIVATION_FAILED", Message: "Activation failed"}
var DATA_INCOMPLETE = HttpResponseStatus{HttpCode: http.StatusBadRequest, Code: "DATA_INCOMPLETE", Message: "Recieved content is not complete"}
var METHOD_NOT_ALLOWED = HttpResponseStatus{HttpCode: http.StatusMethodNotAllowed, Code: "METHOD_NOT_ALLOWED", Message: "The method is not allowed."}
var NOT_IMPLEMENTED = HttpResponseStatus{HttpCode: http.StatusNotImplemented, Code: "NOT_IMPLEMENTED", Message: "Not implemented"}
var TIMED_OUT = HttpResponseStatus{HttpCode: http.StatusRequestTimeout, Code: "TIMED_OUT", Message: "Resource request timed out"}
var SERVICE_UNAVAILABLE = HttpResponseStatus{HttpCode: http.StatusServiceUnavailable, Code: "SERVICE_UNAVAILABLE", Message: "Service temporarily unavailable"}
var BUSY = HttpResponseStatus{HttpCode: http.StatusServiceUnavailable, Code: "BUSY", Message: "Service is busy"}
var UNKNOWN = HttpResponseStatus{HttpCode: http.StatusInternalServerError, Code: "UNKNOWN", Message: "Unknown error"}
var INVALID_SIGNATURE = HttpResponseStatus{HttpCode: http.StatusUnauthorized, Code: "INVALID_SIGNATURE", Message: "Invalid application signature"}
