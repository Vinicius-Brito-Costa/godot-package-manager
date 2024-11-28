package util

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
var config = GetLoggingConfig()
var level string = ""

func Info(str string){
	getLevel()
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
	if len(str) > 0 && slices.Contains(WARN_LIST, getLevel()) {
		log.Println("[WARN][" + TRACE_ID + "=" + ID + "] " + str)
	}
}

func Trace(str string){
	if len(str) > 0 && getLevel() == TRACE{
		log.Println("[TRACE][" + TRACE_ID + "=" + ID + "] " + str)
	}
}

func getLevel() string {
	if len(level) == 0{
		level = config.Level
	}

	return level
}