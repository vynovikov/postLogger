// Assembly point
package main

import (
	"log"
	"os"
	"os/signal"
	"postLogger/internal/adapters/application"
	"postLogger/internal/adapters/driven/saver"
	"postLogger/internal/adapters/driver/rpc"
	"postLogger/internal/logger"
	"syscall"
)

func main() {
	saver, err := saver.NewSaver("logs", int64(100*1024*1024))
	if err != nil {
		log.Fatalf("in postlogger.main.main failed to create saver: %v\n", err)
	}
	app, done := application.NewApp(saver)
	rcv := rpc.NewReceiver(app)

	go rcv.Run()
	go SignalListen(app)
	<-done
	logger.L.Errorln("postLogger is interrupted")

}

// SignalListen listens for Interrupt signal, when receiving one invokes stop function
func SignalListen(app application.Application) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	<-sigChan
	go app.Stop()
}
