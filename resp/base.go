package resp

import (
	"encoding/json"
	"net/http"

	"github.com/saravanan611/loger/loger"
)

const (
	Success = "S"
	Error   = "E"
)

type RespStruct struct {
	Status   string `json:"status,omitempty"`
	ErrCode  string `json:"code,omitempty"`
	Msg      string `json:"msg,omitempty"`
	RespInfo any    `json:"info,omitempty"`
}

func ErrorSender(w http.ResponseWriter, pLog *loger.LogStruct, pErrCode string, pErr error) {
	// w.WriteHeader(http.StatusBadRequest)
	pLog.Err(pErr)
	if lErr := json.NewEncoder(w).Encode(RespStruct{Status: Error, ErrCode: pErrCode, Msg: pErr.Error()}); lErr != nil {
		pLog.Err(lErr)
	}
}

func MsgSender(w http.ResponseWriter, pLog *loger.LogStruct, pInfo any) {
	// w.WriteHeader(http.StatusBadRequest)
	pLog.Info(pInfo)
	if lErr := json.NewEncoder(w).Encode(RespStruct{Status: Success, RespInfo: pInfo}); lErr != nil {
		pLog.Err(lErr)
	}
}
