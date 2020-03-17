package main

import (
	"fmt"
	"github.com/ShrewdSpirit/su/schedule"
	"time"
)

func main() {
	go schedule.Every().Day().At(15, 30).
		From(time.Now().AddDate(0, 0, 1)).
		To(time.Now().AddDate(0, 0, 7)).
		Do(func(task *schedule.Task) {
			fmt.Println("AYYY")
		}, nil).
		Then(func(task *schedule.Task) {
			fmt.Printf("\ttask %s is done\n", task.Id())
		})

	time.Sleep(1 * time.Second)
	schedule.Wait()

	//time.Sleep(3*time.Second)
	//fmt.Printf("stopping %d tasks\n", schedule.NumTasks())
	//schedule.Stop()
	//fmt.Println("done")
}
