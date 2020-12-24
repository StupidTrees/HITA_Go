package logger
import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)
type Level int
var (
	F *os.File
	DefaultPrefix = ""
	DefaultCallerDepth = 2
	logger *log.Logger
	logPrefix = ""
	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)
const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)
func init() {
	filePath := getLogFileFullPath()
	F = openLogFile(filePath)
	logger = log.New(F, DefaultPrefix, log.LstdFlags)
}
func Debugln(v ...interface{}) {
	setPrefix(DEBUG)
	logger.Println(v)
}
func Infoln(v ...interface{}) {
	setPrefix(INFO)
	logger.Println(v)
}
func Println(v ...interface{}) {
	setPrefix(INFO)
	logger.Println(v)
}
func Printf(format string,v ...interface{}) {
	setPrefix(INFO)
	logger.Printf(format,v ...)
}
func Warnln(v ...interface{}) {
	setPrefix(WARNING)
	logger.Println(v)
}
func Errorln(v ...interface{}) {
	setPrefix(ERROR)
	logger.Println(v)
}
func Fatalln(v ...interface{}) {
	setPrefix(FATAL)
	logger.Fatalln(v)
}
func setPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}
	logger.SetPrefix(logPrefix)
}
