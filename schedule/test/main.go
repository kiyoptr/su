package main

import (
	"fmt"
	"github.com/ShrewdSpirit/su/schedule"
	"time"
)

func main() {
	go schedule.Every().Second().
		From(time.Now().Add(3*time.Second)).
		To(time.Now().Add(6*time.Second)).
		Do(func(task *schedule.Task) {
			fmt.Printf("1 Second after %dms\n", task.Elapsed.Milliseconds())
		}, nil)

	time.Sleep(1 * time.Second)
	schedule.Wait()

	//time.Sleep(3*time.Second)
	//fmt.Printf("stopping %d tasks\n", schedule.NumTasks())
	//schedule.Stop()
	//fmt.Println("done")
}
