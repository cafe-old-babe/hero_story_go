package log

import (
	"fmt"
	"log"
	"sync"
)

var writer *dailyFileWriter
var infoLogger, errorLogger *log.Logger

func Config(outputFileName string) {
	writer = &dailyFileWriter{
		fileName:    outputFileName,
		lastYearDay: -1,
		switchLock:  &sync.Mutex{},
	}
	infoLogger = log.New(writer, "[INFO] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix)
	errorLogger = log.New(writer, "[ERROR] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix)

}

func Info(format string, arr ...interface{}) {
	_ = infoLogger.Output(2, fmt.Sprintf(format, arr...))
}
func Error(format string, arr ...interface{}) {
	_ = errorLogger.Output(2, fmt.Sprintf(format, arr...))
}
