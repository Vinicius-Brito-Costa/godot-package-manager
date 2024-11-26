package util

import (
	"log"
	"strings"
	"github.com/google/uuid"
)
const TRACE_ID = "TRACE_ID"
var ID = strings.ReplaceAll(uuid.New().String(), "-", "")

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
	if len(str) > 0 {
		log.Println("[WARN][" + TRACE_ID + "=" + ID + "] " + str)
	}
}

func Trace(str string){
	if len(str) > 0 {
		log.Println("[TRACE][" + TRACE_ID + "=" + ID + "] " + str)
	}
}