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
	pre := log[:strings.Index(log, "in ")+len("in ")]
	logString := fmt.Sprintf("[%s] %s%s", ts, pre, log[len(pre):])

	err := a.S.Save(logString)
	if err != nil {
		logger.L.Errorf("in application.Handle err: %v\n", err)
	}

}
