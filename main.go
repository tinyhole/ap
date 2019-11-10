package main

import (
	"github.com/tinyhole/ap/cmd"
	"os"
	"os/signal"
)

func main() {
	cmd.Execute()
	wait()
}

func wait() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
}
