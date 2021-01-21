package main

import (
	"fmt"
	"time"

	"github.com/kiyoptr/su/ticker"
)

func main() {
	go ticker.Every().Day().At(15, 30).
		From(time.Now().AddDate(0, 0, 1)).
		To(time.Now().AddDate(0, 0, 7)).
		Do(func(task *ticker.Task) {
			fmt.Println("AYYY")
		}, nil).
		Then(func(task *ticker.Task) {
			fmt.Printf("\ttask %s is done\n", task.Id())
		})

	time.Sleep(1 * time.Second)
	ticker.Wait()

	//time.Sleep(3*time.Second)
	//fmt.Printf("stopping %d tasks\n", schedule.NumTasks())
	//schedule.Stop()
	//fmt.Println("done")
}
