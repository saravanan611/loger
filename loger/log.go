package loger

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/google/uuid"
)

/* ----------------------- log session------------------- */
type LogStruct struct {
	lUid      string
	lInfoFlag bool
}

/* -----------------------custom error ------------------- */
type ownErr struct {
	lFileInfo string
	lErr      error
}

/* -----------------------err to string-------------------- */
func (pErr *ownErr) Error() string {
	return pErr.lErr.Error()
}

/* --------------------Log Ination ---------------- */
func Init() *LogStruct {
	var lLogRec LogStruct
	lLogRec.lUid = strings.ReplaceAll(uuid.New().String(), "-", "")
	lLogRec.lInfoFlag = strings.EqualFold(os.Getenv("InfoFlog"), "Y")
	return &lLogRec
}

/* --------------------Basic log variable ---------------- */
var (
	err    = log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.Lshortfile)
	info   = log.New(os.Stdout, "INFO: ", log.LstdFlags|log.Lshortfile)
	lOnErr = log.New(os.Stdout, "INFO: ", log.LstdFlags)
)

// ---------- INFO LOGGER ----------
func (lId *LogStruct) Info(pMsg ...any) {
	if !lId.lInfoFlag {
		return
	}
	msg := fmt.Sprint(pMsg...) // cleaner than fmt.Sprintf("%v")
	info.Output(2, fmt.Sprintf("[%s] %s", lId.lUid, msg))
}

// ---------- ERROR LOGGER ----------
func (lId *LogStruct) Err(pErr any) {
	if lVal, lok := pErr.(error); lok {
		var lOWerr *ownErr
		if errors.As(lVal, &lOWerr) {
			lOnErr.Printf("%s [%s] %s", lOWerr.lFileInfo, lId.lUid, lOWerr.lErr.Error())
			return
		}
	}
	// use Output(counter+2, ...) so call depth aligns properly
	err.Output(2, fmt.Sprintf("[%s] %v", lId.lUid, pErr))
}

// ---------- ERROR WRAPPER ----------
func Return(pErr error) error {
	if pErr == nil {
		return nil
	}

	var lOErr *ownErr

	if !errors.As(pErr, &lOErr) {
		_, lFile, lLine, _ := runtime.Caller(1)
		lStrArray := strings.Split(lFile, "/")
		lFilename := lStrArray[len(lStrArray)-2] + "/" + lStrArray[len(lStrArray)-1]
		return &ownErr{lFileInfo: fmt.Sprintf("%s:%d", lFilename, lLine), lErr: pErr}
	}

	return pErr
}
