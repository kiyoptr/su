package schedule

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

/*
Tasks can be done periodically or once.
Each task has a start time and end time so periodical tasks can have a time limit.

ways to implement waiting for next run:
- busy waiting on main goroutine
- each task wait in their goroutine with time.sleep until next run
-
*/

type TaskFunc func(*Task)

const (
	unitSeconds time.Duration = time.Second
	unitMinutes time.Duration = time.Minute
	unitHours   time.Duration = time.Hour
	unitDays    time.Duration = time.Hour * 24
	unitWeeks   time.Duration = time.Hour * 24 * 7
)

type TaskConfig struct {
	id      uuid.UUID
	handler TaskFunc
	oneShot bool
	next    time.Time // is the next time for this task to run
	lastRun time.Time

	unit     time.Duration // unit of interval (hours, days or what)
	interval time.Duration // number of units to repeat (every 3 seconds, the 3 is interval)
	weekDay  time.Weekday
	hour     int
	minute   int

	task *Task
}

type Task struct {
	config  *TaskConfig
	Payload interface{}
	Elapsed time.Duration
}

/*
CRUD funcs
*/

type Config struct {
}

var (
	tasks      = sync.Map{}
	stopped    = make(chan bool, 1)
	tWaitGroup = sync.WaitGroup{}
	local      = time.Local
)

// Run starts processing tasks
func Run(cfg Config) {
	// TODO: storage: load all
	go func() {
		ticker := time.Tick(1 * time.Millisecond)
		for {
			select {
			case <-ticker:
				runTasks()
			case <-stopped:
				// TODO: storage: save all
				return
			}
		}
	}()
}

// Stop stops processing all tasks. It **MUST** be called whenever the program finishes so tasks will be saved.
func Stop() {
	stopped <- true
	tWaitGroup.Wait()
}

func Every(interval int) *TaskConfig {
	now := time.Now()
	return &TaskConfig{
		oneShot:  false,
		lastRun:  now,
		interval: time.Duration(interval),
		weekDay:  now.Weekday(),
		hour:     now.Hour(),
		minute:   now.Minute(),
	}
}

func runTasks() {
	tasks.Range(func(key, value interface{}) bool {
		tc := value.(*TaskConfig)
		if time.Now().After(tc.next) {
			tWaitGroup.Add(1)

			go func(tc *TaskConfig) {
				tc.task.Elapsed = time.Since(tc.lastRun)
				tc.handler(tc.task)
				tc.lastRun = time.Now()
				if tc.oneShot {
					tc.task.Remove()
				} else {
					tc.calculateNextRun()
				}
				tWaitGroup.Done()
			}(tc)
		}

		return true
	})
}

func (t *TaskConfig) calculateNextRun() {
	if t.unit == unitWeeks {
		now := time.Now()
		remainingDays := t.weekDay - now.Weekday()
		if remainingDays <= 0 {
			// schedule for next week
			t.next = now.AddDate(0, 0, 6-int(now.Weekday())+int(t.weekDay)+1)
		} else {
			t.next = now.AddDate(0, 0, int(remainingDays))
		}

		t.next = time.Date(t.next.Year(), t.next.Month(), t.next.Day(), t.hour, t.minute, 0, 0, local)
		t.next = t.next.Add((t.interval - 1) * t.unit)
	} else if t.unit == unitDays {
		t.next = t.next.Add(t.interval * t.unit)
		t.next = time.Date(t.next.Year(), t.next.Month(), t.next.Day(), t.hour, t.minute, 0, 0, local)
	} else {
		t.next = time.Now().Add(t.interval * t.unit)
	}
}

func (t *TaskConfig) Do(f TaskFunc, payload interface{}) *TaskConfig {
	// TODO: storage: new task
	t.handler = f
	t.task = &Task{
		config:  t,
		Payload: payload,
	}

	t.id, _ = uuid.NewRandom()
	tasks.Store(t.id, t)

	t.calculateNextRun()

	return t
}

func (t *Task) Remove() {
	// TODO: storage: delete task
	tasks.Delete(t.config.id)
}

func (t *TaskConfig) At(hour, minute int) *TaskConfig {
	t.hour = hour
	t.minute = minute
	return t
}

func (t *TaskConfig) Once() *TaskConfig {
	t.oneShot = true
	return t
}

func (t *TaskConfig) Second() *TaskConfig { return t.Seconds() }
func (t *TaskConfig) Seconds() *TaskConfig {
	t.unit = unitSeconds
	return t
}

func (t *TaskConfig) Minute() *TaskConfig { return t.Minutes() }
func (t *TaskConfig) Minutes() *TaskConfig {
	t.unit = unitMinutes
	return t
}

func (t *TaskConfig) Hour() *TaskConfig { return t.Hours() }
func (t *TaskConfig) Hours() *TaskConfig {
	t.unit = unitHours
	return t
}

func (t *TaskConfig) Day() *TaskConfig { return t.Days() }
func (t *TaskConfig) Days() *TaskConfig {
	t.unit = unitDays
	return t
}

func (t *TaskConfig) Week() *TaskConfig { return t.Weeks() }
func (t *TaskConfig) Weeks() *TaskConfig {
	t.unit = unitWeeks
	return t
}

func (t *TaskConfig) Saturday() *TaskConfig {
	t.unit = unitWeeks
	t.weekDay = time.Saturday
	return t
}

func (t *TaskConfig) Sunday() *TaskConfig {
	t.unit = unitWeeks
	t.weekDay = time.Sunday
	return t
}

func (t *TaskConfig) Monday() *TaskConfig {
	t.unit = unitWeeks
	t.weekDay = time.Monday
	return t
}

func (t *TaskConfig) Tuesday() *TaskConfig {
	t.unit = unitWeeks
	t.weekDay = time.Tuesday
	return t
}

func (t *TaskConfig) Wednesday() *TaskConfig {
	t.unit = unitWeeks
	t.weekDay = time.Wednesday
	return t
}

func (t *TaskConfig) Thursday() *TaskConfig {
	t.unit = unitWeeks
	t.weekDay = time.Thursday
	return t
}

func (t *TaskConfig) Friday() *TaskConfig {
	t.unit = unitWeeks
	t.weekDay = time.Friday
	return t
}
