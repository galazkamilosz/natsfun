package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	model "github.com/galazkamilosz/natsfun/models"
	"github.com/nats-io/nats.go"
)

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	file, err := os.Open("./data/input.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	var inData []model.Data
	err = json.Unmarshal(data, &inData)
	if err != nil {
		log.Fatal(err)
	}

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go producer(ctx, i, &wg, &inData)
	}

	<-sigs
	fmt.Println()
	cancel()
	wg.Wait()
}

func builder(group, name string) string {
	return fmt.Sprintf("edgex.reading.%s.%s", group, name)
}

func producer(ctx context.Context, id int, wg *sync.WaitGroup, data *[]model.Data) {
	defer wg.Done()
	counter := 0
	fmt.Printf("Spawning producer [ %d ].\n", id)
	nc, _ := nats.Connect(nats.DefaultURL)
	local := *data
	// ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Producer [ %d ] sent [ %d ] messages.\n", id, counter)
			return
		// case <-ticker.C:
		default:
			for _, in := range local {
				err := nc.Publish(builder(in.Group, in.Name), []byte(fmt.Sprintf("%v", in.Value)))
				counter++
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}
