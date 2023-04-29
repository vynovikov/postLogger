package main

import (
	"log"
	"postLogger/internal/adapters/application"
	"postLogger/internal/adapters/driven/saver"
	"postLogger/internal/adapters/driver/rpc"
)

func main() {
	saver, err := saver.NewSaver("logs", int64(100*1024*1024))
	if err != nil {
		log.Fatalf("in postlogger.main.main failed to create saver: %v\n", err)
	}
	app := application.NewApp(saver)
	rcv := rpc.NewReceiver(app)

	rcv.Run()

}
