package util

/*
 * Copyright © 2024, "DEADLINE TEAM" LLC
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are not permitted.
 *
 * THIS SOFTWARE IS PROVIDED BY "DEADLINE TEAM" LLC "AS IS" AND ANY
 * EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL "DEADLINE TEAM" LLC BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 * No reproductions or distributions of this code is permitted without
 * written permission from "DEADLINE TEAM" LLC.
 * Do not reverse engineer or modify this code.
 *
 * © "DEADLINE TEAM" LLC, All rights reserved.
 */

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
	"sync"
)

func InitLogger() {
	log.SetReportCaller(true)
	log.SetFormatter(NewTextFormatter(true))
	log.AddHook(newFileLoggingHook())
}

type textFormatter struct {
	loggerName string
	needColor  bool
}

func NewTextFormatter(needColor bool) log.Formatter {
	return &textFormatter{needColor: needColor}
}

func (formatter *textFormatter) Format(entry *log.Entry) ([]byte, error) {
	var logFormat []string
	var args []interface{}

	// Append time
	logFormat = append(logFormat, "%s")
	args = append(args, entry.Time.Format("2006-01-02 15:04:05"))

	// Append logging level
	level := strings.ToUpper(entry.Level.String())
	logFormat = append(logFormat, formatByLoggingLevel(level, formatter.needColor))
	args = append(args, level)

	// Append microservice name
	logFormat = append(logFormat, "[%s, %16s]")
	args = append(args, ApplicationName)
	if entry.Context != nil {
		if traceId, ok := entry.Context.Value(traceIdKey).(string); ok {
			args = append(args, traceId[16:])
		} else {
			args = append(args, "background")
		}
	} else {
		args = append(args, "background")
	}

	// Append execution point
	logFormat = append(logFormat, "%s")
	_, fileName := SubstringLast(entry.Caller.File, "/")
	args = append(args, fmt.Sprintf("%25s:%d", fileName, entry.Caller.Line))

	// Append message
	logFormat = append(logFormat, ": %s")
	args = append(args, entry.Message)

	// Append error if exists
	if err := entry.Data["error"]; err != nil {
		printStack := false
		if needStack, ok := entry.Data["needStack"]; ok {
			if val, ok := needStack.(bool); ok {
				printStack = val
			}
		}

		if err, ok := err.(error); ok {
			if printStack {
				logFormat = append(logFormat, ": %s\n%s")
				args = append(args, err.Error(), getStack())
			} else {
				logFormat = append(logFormat, ": %s\n")
				args = append(args, err.Error())
			}
		}
	} else if !strings.HasSuffix(entry.Message, "\n") {
		logFormat = append(logFormat, "\n")
	}

	return []byte(fmt.Sprintf(strings.Join(logFormat, "\t"), args...)), nil
}

func formatByLoggingLevel(level string, needColor bool) string {
	format := "%s"
	if !needColor {
		return format
	}

	switch level {
	case "TRACE":
		format = "\033[1;34m%s\033[0m" // Purple
	case "DEBUG":
		format = "\033[1;36m%s\033[0m" // Teal
	case "INFO":
		format = "\033[1;32m%s\033[0m" // Green
	case "WARNING":
		format = "\033[1;33m%s\033[0m" // Yellow
	case "ERROR":
		format = "\033[1;31m%s\033[0m" // Red
	}

	return format
}

var fileLoggingMutex sync.Mutex

type fileLoggingHook struct {
	formatter log.Formatter
	filename  string
}

func newFileLoggingHook() log.Hook {
	_ = MkDir("log/server.log")
	return &fileLoggingHook{
		formatter: NewTextFormatter(false),
		filename:  "log/server.log",
	}
}

func (hook *fileLoggingHook) Levels() []log.Level {
	return []log.Level{log.PanicLevel, log.FatalLevel, log.ErrorLevel, log.WarnLevel, log.InfoLevel, log.DebugLevel, log.TraceLevel}
}

func (hook *fileLoggingHook) Fire(entry *log.Entry) error {
	go func() {
		logString, err := hook.formatter.Format(entry)
		if err != nil {
			log.WithError(err).Error("Couldn't format entry to file line")
			return
		}

		if err = writeLineToLogFile(hook.filename, logString); err != nil {
			log.WithError(err).Error("Couldn't write line to server.log")
			return
		}
	}()
	return nil
}

func writeLineToLogFile(filename string, line []byte) error {
	fileLoggingMutex.Lock()
	defer fileLoggingMutex.Unlock()

	directory, err := os.Getwd()
	if err != nil {
		return err
	}

	filePath := directory + "/" + filename
	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.WithError(err).Error("Couldn't open server.log")
		return err
	}
	defer CloseChecker(logFile)

	logFileInfo, err := logFile.Stat()
	if err != nil {
		log.WithError(err).Error("Couldn't get file info from server.log")
		return err
	}
	if logFileInfo.Size()+int64(len(line)) > 1024*1024*100 {
		if err = ArchiveFile(filePath); err != nil {
			log.WithError(err).Error("Couldn't archive file server.log")
			return err
		} else {
			logFile, err = os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
			if err != nil {
				log.WithError(err).Error("Couldn't open server.log")
				return err
			}
			defer CloseChecker(logFile)
		}
	}

	_, err = logFile.Write(line)
	if err != nil {
		return err
	}

	return nil
}

func getStack() string {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			break
		}
		buf = make([]byte, 2*len(buf))
	}
	return string(buf)
}
