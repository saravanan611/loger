package apigate

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/saravanan611/loger/loger"
)

type ResponseCaptureWriter struct {
	http.ResponseWriter
	status int
	body   []byte
}

func (rw *ResponseCaptureWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *ResponseCaptureWriter) Write(body []byte) (int, error) {
	rw.body = append(rw.body, body...)
	return rw.ResponseWriter.Write(body)
}

func (rw *ResponseCaptureWriter) Status() int {
	if rw.status == 0 {
		return http.StatusOK
	}
	return rw.status
}

func (rw *ResponseCaptureWriter) Body() []byte {
	return rw.body
}

type SessionResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Url    string `json:"url"`
}

var (
	allowOrigin     = "*"
	allowCredential = false
	allowHeader     = []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "credentials"}
)

func SetHeader(pHeader ...string) {
	if len(pHeader) > 0 {
		allowHeader = append(allowHeader, pHeader...)
	}
}

func SetOrigin(pOrigin string) {
	allowOrigin = pOrigin
}

func EnableCredential() {
	allowCredential = true
}

// Middleware to log requests and route based on API version
func logMiddleware(pNext http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Initialize the logger
		(w).Header().Set("Access-Control-Allow-Origin", allowOrigin)
		(w).Header().Set("Access-Control-Allow-Credentials", fmt.Sprint(allowCredential))
		(w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		(w).Header().Set("Access-Control-Allow-Headers", strings.Join(allowHeader, ","))

		log := loger.Init()
		log.Info("LogMiddleware (+)")
		// Check if it is an OPTIONS request

		requestorDetail := GetRequestorDetail(log, r)

		body, lErr := io.ReadAll(r.Body)
		if lErr != nil {
			log.Err(lErr)
		}

		requestorDetail.Body = string(body)
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		contentLength := r.Header.Get("Content-Length")
		if contentLength != "" {
			length, lErr := strconv.Atoi(contentLength)
			if lErr != nil {
				log.Err(lErr)
			}

			// if length >= 2<<20 { // 2 MB = 2 * 1024 * 1024 bytes
			if length >= 1<<20 { // 1 MB = 1 * 1024 * 1024 bytes
				requestorDetail.Body = "File Data"
			}
		}

		log.Info("Req Info :", requestorDetail)

		// Move the logging of request after setting the context
		captureWriter := &ResponseCaptureWriter{ResponseWriter: w}

		pNext.ServeHTTP(captureWriter, r)

		log.Info("Resp Info :", captureWriter.Body())

		log.Info("LogMiddleware (-)")
	})

}

func SetServer(pRuterFunc func(pRouterInfo *mux.Router), pReadTimeout, pWriteTimeout, pIdleTimeout, pPortAdrs int) error {

	if pPortAdrs < 5000 {
		return loger.Return(fmt.Errorf(" Address must be Greater then or equal to 5000"))
	}

	if pReadTimeout == 0 {
		pReadTimeout = 30
	}
	if pWriteTimeout == 0 {
		pWriteTimeout = 30
	}
	if pIdleTimeout == 0 {
		pIdleTimeout = 120
	}

	lRouter := mux.NewRouter()
	pRuterFunc(lRouter)

	lRouter.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Optional Call Success")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"status":"E","error": "Method %s not allowed on %s"}`, r.Method, r.URL.Path)
	})

	lHandler := logMiddleware(lRouter)
	lSrv := &http.Server{
		ReadTimeout:  time.Duration(pReadTimeout) * time.Second,
		WriteTimeout: time.Duration(pWriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(pIdleTimeout) * time.Second,
		Handler:      lHandler,
		Addr:         fmt.Sprintf(":%d", pPortAdrs),
	}

	log.Printf("server start on :%d ....", pPortAdrs)
	if lErr := lSrv.ListenAndServe(); lErr != nil {
		return loger.Return(lErr)
	}

	return nil
}
