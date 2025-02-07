package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ckshitij/cache/pkg/cache"
)

func multiSignalHandler(sig os.Signal, done chan bool) {
	switch sig {
	case syscall.SIGHUP:
		fmt.Println("Signal: syscall.SIGHUP ", sig.String())
		done <- true
		time.Sleep(1 * time.Second)
		close(done)
		os.Exit(0)
	case syscall.SIGINT:
		fmt.Println("Signal: syscall.SIGINT ", sig.String())
		done <- true
		time.Sleep(1 * time.Second)
		close(done)
		os.Exit(0)
	case syscall.SIGTERM:
		fmt.Println("Signal: syscall.SIGTERM ", sig.String())
		done <- true
		time.Sleep(1 * time.Second)
		close(done)
		os.Exit(0)
	default:
		fmt.Println("Unhandled/unknown signal")
	}
}

/*
Demo to how to consume the inmemory key value datastore
*/
func main() {
	sigchnl := make(chan os.Signal, 1)
	signal.Notify(
		sigchnl,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
	) // we can add more sycalls.SIGQUIT etc.

	done := make(chan bool)
	go func() {
		for {
			s := <-sigchnl
			multiSignalHandler(s, done)
		}
	}()

	ds := cache.NewKeyValueCache[string](1 * time.Second)

	go ds.AutoCleanUp(3*time.Second, done)

	var i int64 = 0
	for {
		key := fmt.Sprintf(" key-%d", i+1)
		value := fmt.Sprintf(" value-%d", i+1)
		ds.Put(key, value)

		time.Sleep(100 * time.Millisecond)
		i++
	}
}
