package log

import (
	"fmt"
	"io"
	golog "log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/ShrewdSpirit/su/errors"
)

var (
	logger      *golog.Logger
	currentFile *os.File
	writer      io.Writer
	rootDir     string
	currentDay  time.Time
	lastFlush   time.Time

	FlushDelay = 100 * time.Millisecond

	queue = []string{}
	lock  = sync.Mutex{}
)

func Load(logsDir string) error {
	rootDir = logsDir

	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		if err := os.MkdirAll(rootDir, os.ModePerm); err != nil {
			return errors.Newi(err, "failed to create logs directory")
		}
	}

	logger = golog.New(os.Stdout, "", golog.LstdFlags)
	if err := setOutputFile(); err != nil {
		return err
	}

	return nil
}

func Close() {
	flush()
	if currentFile != nil {
		currentFile.Close()
	}
}

func Log(msg string) {
	Logs("I", msg)
}

func Logf(format string, params ...interface{}) {
	Logfs("I", format, params...)
}

func Logfs(state, format string, params ...interface{}) {
	addQueue(state, fmt.Sprintf(format, params...))
}

func Logs(state, msg string) {
	addQueue(state, msg)
}

func Error(err error) {
	addQueue("E", err.Error())
}

func setOutputFile() (err error) {
	currentDay = time.Now()
	logFilename := path.Join(rootDir, fmt.Sprintf("log-%v-%v-%v.txt", currentDay.Year(), int(currentDay.Month()), currentDay.Day()))

	currentFile, err = os.OpenFile(logFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		err = errors.Newfi(err, "failed to create log file %v", logFilename)
		return
	}

	writer = io.MultiWriter(os.Stdout, currentFile)
	logger.SetOutput(writer)

	return
}

func checkLogFileDate() {
	now := time.Now()
	if now.Year() != currentDay.Year() || now.Month() != currentDay.Month() || now.Day() != currentDay.Day() {
		Close()
		if err := setOutputFile(); err != nil {
			fmt.Println(err)
			fmt.Println("LOGS WON'T BE WRITTEN IN LOG FILES")
			logger.SetOutput(os.Stdout)
		}
	} else {
		flush()
	}
}

func addQueue(state, message string) {
	text := fmt.Sprintf("[%v] %v\n", state, message)

	lock.Lock()
	queue = append(queue, text)
	lock.Unlock()

	if time.Now().After(lastFlush.Add(FlushDelay)) {
		checkLogFileDate()
	}
}

func flush() {
	lastFlush = time.Now()

	lock.Lock()
	for _, m := range queue {
		logger.Printf(m)
	}
	queue = []string{}
	lock.Unlock()
}
