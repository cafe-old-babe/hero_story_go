package log

import (
	"github.com/pkg/errors"
	"io"
	"os"
	"path"
	"sync"
	"time"
)

type dailyFileWriter struct {
	//日志文件名称
	fileName string
	//上次写入文件日期
	lastYearDay int
	//输出日志文件
	outputFile *os.File
	switchLock *sync.Mutex
}

func (d *dailyFileWriter) Write(byteArr []byte) (n int, err error) {
	if byteArr == nil || len(byteArr) <= 0 {
		return 0, nil
	}
	outputFile, err := d.getOutputFile()
	if err != nil {
		return 0, err
	}

	// write console
	_, _ = os.Stderr.Write(byteArr)
	_, _ = outputFile.Write(byteArr)
	return 0, nil
}

//获取输出文件
//每天创建一个日志文件
func (d *dailyFileWriter) getOutputFile() (io.Writer, error) {
	yearDay := time.Now().YearDay()
	if d.lastYearDay == yearDay && d.outputFile != nil {
		return d.outputFile, nil
	}
	d.switchLock.Lock()
	defer d.switchLock.Unlock()

	if d.lastYearDay == yearDay && d.outputFile != nil {
		return d.outputFile, nil
	}

	d.lastYearDay = yearDay
	err := os.MkdirAll(path.Dir(d.fileName), os.ModePerm)
	if err != nil {
		return nil, err
	}
	newFileName := d.fileName + "." + time.Now().Format("20060102")
	newOutputFile, err := os.OpenFile(newFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil || newOutputFile == nil {
		return nil, errors.Errorf(`create log file: %v, failed: %v`, newFileName, err.Error())
	}
	if d.outputFile != nil {
		_ = d.outputFile.Close()
	}
	d.outputFile = newOutputFile
	return newOutputFile, nil

}
