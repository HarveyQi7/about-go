package logger

import (
	"about-go/tools"
	"os"
	"time"
)

type DailyRotator struct {
	Level string
	*os.File
}

const (
	INFO  = "info"
	ERROR = "error"
)

const dir = "logs"

func getPath(level string) string {
	if !tools.IsPathExist(dir) {
		os.Mkdir(dir, 0666)
	}
	date := time.Now().Format("2006-01-02")
	return dir + "/" + date + "_" + level + ".log"
}

func getFile(path string) (*os.File, error) {
	if tools.IsPathExist(path) {
		return os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0666)
	} else {
		return os.Create(path)
	}
}

func Writer(level string) *DailyRotator {
	file, err := getFile(getPath(level))
	if err != nil {
		panic(err)
	}
	return &DailyRotator{level, file}
}

func (l *DailyRotator) Write(p []byte) (n int, err error) {
	path := getPath(l.Level)
	if !tools.IsPathExist(path) {
		f, e := getFile(path)
		if e != nil {
			return 0, e
		} else {
			l.File.Close()
			l.File = f
		}
	}
	return l.File.Write(p)
}
