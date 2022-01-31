package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	fmt.Println("spawning 5 consumers")
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go consumer(ctx, i, &wg)
	}
	<-sigs
	fmt.Println()
	cancel()
	wg.Wait()
}

func consumer(ctx context.Context, id int, wg *sync.WaitGroup) {

	defer wg.Done()
	// defer fmt.Printf("Consumer [ %d ] stopped.\n", id)
	fmt.Printf("Consumer [ %d ] started.\n", id)
	nc, _ := nats.Connect(nats.DefaultURL)
	counter := 0

	nc.Subscribe("edgex.reading.*.*", func(m *nats.Msg) {
		counter++
		//this subscriber performs some strange operations
		time.Sleep(5 * time.Millisecond)
	})

	<-ctx.Done()
	fmt.Printf("Consumer [ %d ] stopping, received [ %d ] messages.\n", id, counter)
}
