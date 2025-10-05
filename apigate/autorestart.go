package apigate

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func AutoRestart(lRestartTime, lRestartMinute int) {

	go func() {

		lCurrentTime := time.Now()
		if !(lCurrentTime.Hour() == lRestartTime && lCurrentTime.Minute() == lRestartMinute) {
			lNextExecutionTime := time.Date(lCurrentTime.Year(), lCurrentTime.Month(), lCurrentTime.Day(), lRestartTime, lRestartMinute, 0, 0, lCurrentTime.Location())
			if lNextExecutionTime.Before(lCurrentTime) {
				lNextExecutionTime = lNextExecutionTime.Add(time.Duration(24 * time.Hour))
			}
			fmt.Println("Current Time:", lCurrentTime)
			fmt.Println("Next Execution Time:", lNextExecutionTime, lNextExecutionTime.Sub(lCurrentTime))
			durationUntilNextExecution := lNextExecutionTime.Sub(lCurrentTime)
			time.Sleep(durationUntilNextExecution)
		}

		log.Println("program going to restart within a minute...")
		time.Sleep(1 * time.Minute)

		lErr := restart()

		if lErr != nil {
			log.Println("Error restarting program:", lErr)
		}
		os.Exit(0)
	}()
}

func restart() (lErr error) {
	execPath, lErr := os.Executable()
	if lErr != nil {
		return lErr
	}

	cmd := exec.Command(execPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	lErr = cmd.Start()
	if lErr != nil {
		return lErr
	}
	return
}
