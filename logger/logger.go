package logger

import (
	"log"
	"slices"
	"strings"
	"github.com/google/uuid"
)
const TRACE_ID = "TRACE_ID"
const TRACE = "trace"
const WARN = "warn"
var WARN_LIST []string = []string{ TRACE, WARN }
var ID string = strings.ReplaceAll(uuid.New().String(), "-", "")
var level string = ""

func Info(str string){
	if len(str) > 0 {
		log.Println("[INFO][" + TRACE_ID + "=" + ID + "] " + str)
	}
}

func Error(str string, err error){
	if len(str) > 0{
		if err != nil {
			log.Println("[Error][" + TRACE_ID + "=" + ID + "] " + str + " on error: " + err.Error())
		} else {
			log.Println("[Error][" + TRACE_ID + "=" + ID + "] " + str + " on error: ")
		}
		
	}
}

func Warn(str string){
	if len(str) > 0 && slices.Contains(WARN_LIST, GetLogLevel()) {
		log.Println("[WARN][" + TRACE_ID + "=" + ID + "] " + str)
	}
}

func Trace(str string){
	if len(str) > 0 && GetLogLevel() == TRACE{
		log.Println("[TRACE][" + TRACE_ID + "=" + ID + "] " + str)
	}
}

func SetLogLevel(lvl string) {
	if len(lvl) > 0 {
		level = lvl
	}
}

func GetLogLevel() string {
	return level
}