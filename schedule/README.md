To start a job, start configuring by calling `Every` with zero or one argument that specifies the repeat interval.
After that you can chain calls to configure how often the task will repeat.

Example:
```go
go schedule.Every().Day().At(15, 30).
    From(time.Now().AddDate(0, 0, 1)).
    To(time.Now().AddDate(0, 0, 7)).
    Do(func(task *schedule.Task) {
        fmt.Printf("task %s\n", task.Id())
    }, nil).
    Then(func(task *schedule.Task) {
        fmt.Printf("\ttask %s is done\n", task.Id())
    })
```