package main

import (
	"fmt"
	"github.com/ShrewdSpirit/su/schedule"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}

	schedule.Every(1).Seconds().Do(func(task *schedule.Task) {
		fmt.Printf("1 sec after %dms\n", task.Elapsed.Milliseconds())
	}, nil)

	wg.Add(1)
	schedule.Run(schedule.Config{})
	wg.Wait()
}
