package schedule

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type TaskFunc func(*Task)

const (
	unitSeconds = time.Second
	unitMinutes = time.Minute
	unitHours   = time.Hour
	unitDays    = time.Hour * 24
	unitWeeks   = time.Hour * 24 * 7
)

type TaskConfig struct {
	id       uuid.UUID
	handler  TaskFunc
	oneShot  bool
	nextStep time.Time // is the nextStep time for this task to run
	lastRun  time.Time

	unit         time.Duration // unit of interval (hours, days or what)
	interval     time.Duration // number of units to repeat (every 3 seconds, the 3 is interval)
	weekDay      time.Weekday
	hour, minute int
	from, to     time.Time

	task *Task
}

type Task struct {
	config  *TaskConfig
	Payload interface{}
	Elapsed time.Duration
}

var Config = struct {
	MaxTasks   int
	TaskWaitUs time.Duration
}{
	MaxTasks:   1024,
	TaskWaitUs: 1,
}

var (
	openTasks  = make(chan int, Config.MaxTasks)
	stopSignal = make(chan int, Config.MaxTasks)
	tWaitGroup = sync.WaitGroup{}
	local      = time.Local
)

// Stop stops processing all tasks. It **MUST** be called whenever the program finishes so tasks will be saved.
func Stop() {
	for i := 0; i < len(openTasks); i++ {
		stopSignal <- 1
	}
	Wait()
}

func Wait() { tWaitGroup.Wait() }

func NumTasks() int { return len(openTasks) }

// Every begins configuring a task. Supply zero or one intervals. No intervals will be counted as 1
func Every(interval ...int) *TaskConfig {
	i := 1
	if len(interval) > 0 {
		i = interval[0]
	}

	now := time.Now()
	return &TaskConfig{
		oneShot:  false,
		lastRun:  now,
		interval: time.Duration(i),
		weekDay:  now.Weekday(),
		hour:     now.Hour(),
		minute:   now.Minute(),
	}
}

func (t *TaskConfig) calculateNextRun() {
	if t.unit == unitWeeks {
		now := time.Now()
		remainingDays := t.weekDay - now.Weekday()
		if remainingDays <= 0 {
			// schedule for nextStep week
			t.nextStep = now.AddDate(0, 0, 6-int(now.Weekday())+int(t.weekDay)+1)
		} else {
			t.nextStep = now.AddDate(0, 0, int(remainingDays))
		}

		t.nextStep = time.Date(t.nextStep.Year(), t.nextStep.Month(), t.nextStep.Day(), t.hour, t.minute, 0, 0, local)
		t.nextStep = t.nextStep.Add((t.interval - 1) * t.unit)
	} else if t.unit == unitDays {
		t.nextStep = t.nextStep.Add(t.interval * t.unit)
		t.nextStep = time.Date(t.nextStep.Year(), t.nextStep.Month(), t.nextStep.Day(), t.hour, t.minute, 0, 0, local)
	} else {
		t.nextStep = time.Now().Add(t.interval * t.unit)
	}
}

func (t *TaskConfig) Remove() {
	// TODO: storage: Remove
}

func (t *TaskConfig) Do(f TaskFunc, payload interface{}) {
	tWaitGroup.Add(1)
	defer tWaitGroup.Done()

	// TODO: storage: new task
	t.handler = f
	t.task = &Task{
		config:  t,
		Payload: payload,
	}

	t.id, _ = uuid.NewRandom()
	t.calculateNextRun()

	openTasks <- 1
	ticker := time.NewTicker(Config.TaskWaitUs * time.Microsecond)
	for {
		select {
		case <-ticker.C:
			if time.Since(t.nextStep) > 0 {

				if time.Now().After(t.from) {
					t.task.Elapsed = time.Since(t.lastRun)
					if t.handler != nil {
						t.handler(t.task)
					}
					t.lastRun = time.Now()
					if t.oneShot {
						goto end
					}
				}
				t.calculateNextRun()

				if t.to.Year() != 1 && t.nextStep.After(t.to) {
					goto end
				}
			}
		case <-stopSignal:
			return
		}
	}

end:
	t.Remove()
	<-openTasks
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

func (t *TaskConfig) From(from time.Time) *TaskConfig {
	t.from = from
	return t
}

func (t *TaskConfig) To(to time.Time) *TaskConfig {
	t.to = to
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
