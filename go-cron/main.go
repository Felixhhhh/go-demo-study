package main

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"time"
)

func task() {
	fmt.Println("I am runnning task.", time.Now())
}
func superWang() {
	fmt.Println("I am runnning superWang.", time.Now())
}

func main() {
	s := gocron.NewScheduler()
	s.Every(1).Second().Do(task)
	s.Every(4).Second().Do(superWang)
	sc := s.Start() // keep the channel
	go test(s, sc)  // wait
	<-sc            // it will happens if the channel is closed}
}

func test(s *gocron.Scheduler, sc chan bool) {
	time.Sleep(8 * time.Second)
	s.Remove(task) //remove task
	time.Sleep(8 * time.Second)
	s.Clear()
	fmt.Println("All task removed")
	close(sc) // close the channel
}
