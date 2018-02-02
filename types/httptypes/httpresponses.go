package httptypes

import (
	"github.com/emicklei/go-restful"
	"math"
	"net/http"
)

type HttpResponseStatus struct {
	HttpCode int    `json:"status"`
	Type     string `json:"type"`
	Message  string `json:"message"`
}

type HttpResponsePack struct {
	Status  HttpResponseStatus `json:"status"`
	Content interface{}        `json:"content"`
}

type HttpRestError struct {
	Status *HttpResponseStatus `json:"error"`
}

type HttpPaginatedResponse struct {
	Limit      int         `json:"limit"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	TotalPages int         `json:"totalPages"`
	Data       interface{} `json:"data"`
}

var EMPTY_CONTENT = map[string]string{}

func SendPaginated(r *restful.Response, content interface{}, total, offset, limit int) error {
	response := HttpPaginatedResponse{
		Limit:      limit,
		Total:      total,
		Page:       int(math.Floor(float64(total) / float64(offset))),
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
		Data:       content}

	return r.WriteHeaderAndJson(http.StatusOK, response, restful.MIME_JSON)
}

func SendEntity(r *restful.Response, content interface{}) error {
	if content == nil {
		return r.WriteHeaderAndJson(http.StatusNoContent, nil, restful.MIME_JSON)
	}
	return r.WriteHeaderAndJson(http.StatusOK, content, restful.MIME_JSON)
}

func SendOK(r *restful.Response, data interface{}) error {
	return SendEntity(r, data)
}

func SendError(r *restful.Response, resp *HttpResponseStatus) error {
	httpErr := HttpRestError{Status: resp}
	return r.WriteHeaderAndJson(resp.HttpCode, httpErr, restful.MIME_JSON)
}

func SendDuplicated(r *restful.Response) error {
	return SendError(r, &DUPLICATED)
}

func SendNotFound(r *restful.Response) error {
	return SendError(r, &NOT_FOUND)
}

func SendInvalidData(r *restful.Response) error {
	return SendError(r, &INVALID_DATA)
}

func SendGeneralError(r *restful.Response) error {
	return SendError(r, &DATASOURCE_ERROR)
}

func SendBadAuth(r *restful.Response) error {
	return SendError(r, &INVALID_CREDENTIALS)
}

var HTTP_RESPONSE_OK = HttpResponseStatus{HttpCode: http.StatusOK, Type: "OK", Message: "OK"}
var DUPLICATED = HttpResponseStatus{HttpCode: http.StatusConflict, Type: "DUPLICATED", Message: "Duplicated content or resource"}
var NOT_FOUND = HttpResponseStatus{HttpCode: http.StatusNotFound, Type: "NOT_FOUND", Message: "Resource not found"}
var DATASOURCE_ERROR = HttpResponseStatus{HttpCode: http.StatusInternalServerError, Type: "DATASOURCE_ERROR", Message: "An error occurred with datasource (database)"}
var PERMISSION_DENIED = HttpResponseStatus{HttpCode: http.StatusForbidden, Type: "PERMISSION_DENIED", Message: "Permission denied"}
var UNAUTHORIZED = HttpResponseStatus{HttpCode: http.StatusUnauthorized, Type: "UNAUTHORIZED", Message: "Unauthorized request"}
var UNAVAILABLE = HttpResponseStatus{HttpCode: http.StatusGone, Type: "UNAVAILABLE", Message: "Requested resource or content is not available"}
var REMOVE_FAILED = HttpResponseStatus{HttpCode: http.StatusInternalServerError, Type: "REMOVE_FAILED", Message: "Resource or content remove failed"}
var AUTH_FAILED = HttpResponseStatus{HttpCode: http.StatusInternalServerError, Type: "AUTH_FAILED", Message: "Authentication failed"}
var INVALID_DATA = HttpResponseStatus{HttpCode: http.StatusBadRequest, Type: "INVALID_DATA", Message: "Invalid data"}
var INVALID_CREDENTIALS = HttpResponseStatus{HttpCode: http.StatusUnauthorized, Type: "INVALID_CREDENTIALS", Message: "Invalid credentials (login, password)"}
var INVALID_SENSOR = HttpResponseStatus{HttpCode: http.StatusBadRequest, Type: "INVALID_SENSOR", Message: "Sensor(s) is invalid or not exists"}
var INVALID_TOKEN = HttpResponseStatus{HttpCode: http.StatusBadRequest, Type: "INVALID_TOKEN", Message: "Your token is not valid"}
var TOKEN_EXPIRED = HttpResponseStatus{HttpCode: http.StatusUnauthorized, Type: "TOKEN_EXPIRED", Message: "Token expired. Please re-login"}
var USER_EXISTS = HttpResponseStatus{HttpCode: http.StatusConflict, Type: "USER_EXISTS", Message: "User %s exists"}
var WEAK_PASSWORD = HttpResponseStatus{HttpCode: http.StatusBadRequest, Type: "WEAK_PASSWORD", Message: "Chosen password is too weak"}
var INVALID_PASSWORD = HttpResponseStatus{HttpCode: http.StatusBadRequest, Type: "INVALID_PASSWORD", Message: "Chosen password is not valid"}
var INVALID_USERNAME = HttpResponseStatus{HttpCode: http.StatusBadRequest, Type: "INVALID_USERNAME", Message: "Chosen username is not valid"}
var INVALID_EMAIL = HttpResponseStatus{HttpCode: http.StatusBadRequest, Type: "INVALID_EMAIL", Message: "Your email is not valid"}
var PASSWORD_MISMATCH = HttpResponseStatus{HttpCode: http.StatusBadRequest, Type: "PASSWORD_MISMATCH", Message: "Passwords is not match"}
var ACTIVATION_FAILED = HttpResponseStatus{HttpCode: http.StatusInternalServerError, Type: "ACTIVATION_FAILED", Message: "Activation failed"}
var DATA_INCOMPLETE = HttpResponseStatus{HttpCode: http.StatusBadRequest, Type: "DATA_INCOMPLETE", Message: "Recieved content is not complete"}
var METHOD_NOT_ALLOWED = HttpResponseStatus{HttpCode: http.StatusMethodNotAllowed, Type: "METHOD_NOT_ALLOWED", Message: "The method is not allowed."}
var NOT_IMPLEMENTED = HttpResponseStatus{HttpCode: http.StatusNotImplemented, Type: "NOT_IMPLEMENTED", Message: "Not implemented"}
var TIMED_OUT = HttpResponseStatus{HttpCode: http.StatusRequestTimeout, Type: "TIMED_OUT", Message: "Resource request timed out"}
var SERVICE_UNAVAILABLE = HttpResponseStatus{HttpCode: http.StatusServiceUnavailable, Type: "SERVICE_UNAVAILABLE", Message: "Service temporarily unavailable"}
var BUSY = HttpResponseStatus{HttpCode: http.StatusServiceUnavailable, Type: "BUSY", Message: "Service is busy"}
var UNKNOWN = HttpResponseStatus{HttpCode: http.StatusInternalServerError, Type: "UNKNOWN", Message: "Unknown error"}
var INVALID_SIGNATURE = HttpResponseStatus{HttpCode: http.StatusUnauthorized, Type: "INVALID_SIGNATURE", Message: "Invalid application signature"}
