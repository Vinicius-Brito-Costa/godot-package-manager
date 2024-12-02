package logger

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

// There's is definetly a better way of doing logging
const YYYY_MM_DD = "2006_01_02"
const TRACE_ID = "TRACE_ID"
const TRACE = "trace"
const WARN = "warn"

var ip, _ = externalIP()
var currentTime = time.Now().Local()
var logFileName = "gpm_log_" + ip + "_" + currentTime.Format(YYYY_MM_DD) + ".log"
var WARN_LIST []string = []string{ TRACE, WARN }
var ID string = strings.ReplaceAll(uuid.New().String(), "-", "")
var level string = "info"
var gpmFolder string = ""

var infoLog log.Logger = *log.New(os.Stdout, "[INFO ][DATE_TIME=", log.Ldate|log.Ltime)
var warnLog log.Logger = *log.New(os.Stdout, "[WARN ][DATE_TIME=", log.Ldate|log.Ltime)
var traceLog log.Logger = *log.New(os.Stdout, "[TRACE][DATE_TIME=", log.Ldate|log.Ltime)
var errorLog log.Logger = *log.New(os.Stderr, "[ERROR][DATE_TIME=", log.Ldate|log.Ltime)

func setLoggerOutputToFile(logger *log.Logger) os.File {
	if len(gpmFolder) == 0 {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("cannot get user home directory", err)
			return os.File{}
		}
		gpmFolder = filepath.Join(homeDir, ".godot-package-manager")
	}
	_ = os.Mkdir(gpmFolder, os.ModePerm)
	var logWriter, err = os.OpenFile(filepath.Join(gpmFolder, logFileName), os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error on openFile: " + err.Error())
		return os.File{}
	}

	logger.SetOutput(logWriter)
	return *logWriter
}

func Info(str string){
	if len(str) > 0 {
		fmt.Println(str)
		writer := setLoggerOutputToFile(&infoLog)
		defer writer.Close()
		infoLog.Println("][" + TRACE_ID + "=" + ID + "] " + str)
	}
}

func Error(str string, err error){
	if len(str) > 0{
		writer := setLoggerOutputToFile(&errorLog)
		defer writer.Close()
		if err != nil {
			fmt.Println(str + " on error: " + err.Error())
			errorLog.Println("][" + TRACE_ID + "=" + ID + "] " + str + " on error: " + err.Error())
		} else {
			fmt.Println(str + " on error: ")
			errorLog.Println("][" + TRACE_ID + "=" + ID + "] " + str + " on error: ")
		}
		
	}
}

func Warn(str string){
	if len(str) > 0 && slices.Contains(WARN_LIST, GetLogLevel()) {
		fmt.Println(str)
	}
	writer := setLoggerOutputToFile(&warnLog)
	defer writer.Close()
	warnLog.Println("][" + TRACE_ID + "=" + ID + "] " + str)
}

func Trace(str string){
	if len(str) > 0 && GetLogLevel() == TRACE{
		fmt.Println(str)
	}
	writer := setLoggerOutputToFile(&traceLog)
	defer writer.Close()
	traceLog.Println("][" + TRACE_ID + "=" + ID + "] " + str)
}

func SetLogLevel(lvl string) {
	if len(lvl) > 0 {
		level = lvl
	}
}

func GetLogLevel() string {
	return level
}

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}
