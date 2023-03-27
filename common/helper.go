package common

import (
	"strconv"

	"github.com/go-chassis/openlog"
)

var Errorcodes = make(map[string]interface{})
var Dbname string

type HTTPResponse struct {
	Msg        string      `json:"_msg"`
	Status     int         `json:"_status"`
	Data       interface{} `json:"_data"`
	TotalCount int64       `json:"_totalcount"`
	ErrorCode  int64       `json:"_errorcode"`
}

func ErrorHandler(errcode string, response interface{}, totalcount int64, language string) HTTPResponse {
	if language == "" {
		language = "en"
	}

	ErrDetailes := Errorcodes[errcode].(map[string]interface{})
	msg := ErrDetailes[language].(string)
	Status := ErrDetailes["status"].(string)
	status, err := strconv.ParseInt(Status, 10, 64)
	if err != nil {
		openlog.Error(err.Error())
		return HTTPResponse{Msg: "Internal server error", Status: 500}
	}
	code, err := strconv.ParseInt(errcode, 10, 64)
	if err != nil {
		openlog.Error(err.Error())
		return HTTPResponse{Msg: "Internal server error", Status: 500}
	}

	return HTTPResponse{Msg: msg, Status: int(status), Data: response, TotalCount: totalcount, ErrorCode: code}

}
