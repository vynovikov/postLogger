package application

import (
	"fmt"
	"postLogger/internal/adapters/driven/saver"
	"postLogger/internal/logger"
	"strings"
)

type ApplicationStruct struct {
	S saver.Saver
}

func NewApp(s saver.Saver) *ApplicationStruct {
	return &ApplicationStruct{
		S: s,
	}
}

type Application interface {
	Handle(string, string)
}

func (a *ApplicationStruct) Handle(ts string, log string) {
	//logger.L.Infof("in application.Handle ts: %q, log: %q\n", ts, log)
	pre := log[:strings.Index(log, "in ")+len("in ")]
	logString := fmt.Sprintf("[%s] %s%s%s", ts, pre, "postParser.", log[len(pre):])
	//logger.L.Infof("in application.Handle logString: %q", logString)

	err := a.S.Save(logString)
	if err != nil {
		logger.L.Errorf("in application.Handle err: %v\n", err)
	}

}
