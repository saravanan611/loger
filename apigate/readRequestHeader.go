package apigate

import (
	"net/http"
	"strings"

	"github.com/saravanan611/loger/loger"
)

type RequestorDetails struct {
	RealIP      string
	ForwardedIP string
	Method      string
	Path        string
	Host        string
	RemoteAddr  string
	Header      string
	Body        string
	EndPoint    string
	RequestType string
}

// --------------------------------------------------------------------
// get request header details
// --------------------------------------------------------------------
func GetHeaderDetails(pLog *loger.LogStruct, r *http.Request) string {
	pLog.Info("GetHeaderDetails (+)")
	value1 := ""
	for name, values := range r.Header {
		// Loop over all values for the name.
		for _, value := range values {
			value1 = value1 + " " + name + "-" + value
		}
	}
	pLog.Info("GetHeaderDetails (-)")
	return value1
}

// --------------------------------------------------------------------
// function reads the API requestor details and send return them
// as structure to the caller
// --------------------------------------------------------------------
func GetRequestorDetail(pLog *loger.LogStruct, r *http.Request) RequestorDetails {
	pLog.Info("GetRequestorDetail (+)")

	var reqDtl RequestorDetails
	reqDtl.RealIP = r.Header.Get("Referer")
	reqDtl.ForwardedIP = r.Header.Get("X-Forwarded-For")
	reqDtl.Method = r.Method
	reqDtl.Path = r.URL.Path + "?" + r.URL.RawQuery
	reqDtl.Host = r.Host
	reqDtl.RemoteAddr = r.RemoteAddr
	if strings.Contains(r.URL.Path, "/order/placeorder/") {
		reqDtl.EndPoint = r.URL.Path[:len("/order/placeorder/")]
	} else if strings.Contains(r.URL.Path, "/deals/count/") {
		reqDtl.EndPoint = r.URL.Path[:len("/deals/count/")]
	} else {
		reqDtl.EndPoint = r.URL.Path
	}
	reqDtl.RequestType = r.Header.Get("Content-Type")

	reqDtl.Header = GetHeaderDetails(pLog, r)
	//body, _ := ioutil.ReadAll(r.Body)
	//reqDtl.Body = string(body)
	pLog.Info("GetRequestorDetail (-)")

	return reqDtl
}
