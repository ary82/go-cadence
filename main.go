package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ary82/go-cadence/worker"
)

func main() {
	worker.StartWorker()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT)
	fmt.Println("Cadence worker started, press ctrl+c to terminate...")
	<-done
}
