package application

import (
	"fmt"
	"postLogger/internal/adapters/driven/saver"
	"postLogger/internal/logger"
	"strings"
	"time"
)

type ApplicationStruct struct {
	stopping bool
	done     chan struct{}
	S        saver.Saver
}

func NewApp(s saver.Saver) (*ApplicationStruct, chan struct{}) {
	done := make(chan struct{})
	return &ApplicationStruct{
		S:    s,
		done: done,
	}, done
}

type Application interface {
	Handle(string, string)
	Stop()
}

func (a *ApplicationStruct) Handle(ts string, log string) {
	pre := log[:strings.Index(log, "in ")+len("in ")]
	logString := fmt.Sprintf("[%s] %s%s", ts, pre, log[len(pre):])
	if strings.Contains(log, "is down") && a.stopping {
		close(a.done)
	}

	err := a.S.Save(logString)
	if err != nil {
		logger.L.Errorf("in application.Handle err: %v\n", err)
	}

}
func (a *ApplicationStruct) Stop() {
	a.stopping = true
	time.Sleep(time.Millisecond * 100)
	close(a.done)
}
