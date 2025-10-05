package resp

import (
	"encoding/json"
	"fmt"
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

/* standard error responce structure for api */

func ErrorSender(w http.ResponseWriter, pLog *loger.LogStruct, pErrCode string, pErr error) {
	pLog.Err(pErr)
	w.WriteHeader(http.StatusInternalServerError)
	if _, lErr := fmt.Fprintf(w, "Error: << %s >>. Please refer to this code for developer fast support: (%s).", pErr.Error(), pErrCode); lErr != nil {
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
