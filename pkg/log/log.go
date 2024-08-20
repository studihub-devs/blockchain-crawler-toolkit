package log

import (
	"fmt"
	"log"
	"new-token/pkg/server"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Log Variable
var logger *logrus.Logger
var crawlFailLooger *logrus.Logger

// Log Level Data Type
type logLevel string

// Log Level Data Type Constant
const (
	LogLevelPanic logLevel = "panic"
	LogLevelFatal logLevel = "fatal"
	LogLevelError logLevel = "error"
	LogLevelWarn  logLevel = "warn"
	LogLevelDebug logLevel = "debug"
	LogLevelTrace logLevel = "trace"
	LogLevelInfo  logLevel = "info"
)

// Initialize Function in Helper Logging
func init() {
	// Initialize Log as New Logrus Logger
	logger = logrus.New()

	// Set Log Format to JSON Format
	logger.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: false,
		TimestampFormat:  time.RFC3339Nano,
	})

	if strings.ToLower(server.Config.GetString("LOG_OUTPUT")) == "console" {
		// Set Log Output to STDOUT
		logger.SetOutput(os.Stdout)
	} else {
		filePath := server.Config.GetString("LOG_OUTPUT")
		file, error := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if error != nil {
			log.Fatalln(error)
		}

		fmt.Println("Log into file : ", filePath)
		logger.SetOutput(file)
	}

	crawlFailLooger = logrus.New()

	crawlFailLooger.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: false,
		TimestampFormat:  time.RFC3339Nano,
	})

	if strings.ToLower(server.Config.GetString("LOG_FILE_PATH_FOR_CRAWFAIL")) == "console" {
		// Set Log Output to STDOUT
		crawlFailLooger.SetOutput(os.Stdout)
	} else {
		filePathCrawlFail := server.Config.GetString("LOG_FILE_PATH_FOR_CRAWFAIL")
		fileCrawlFail, error := os.OpenFile(filePathCrawlFail, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if error != nil {
			log.Fatalln(error)
		}

		fmt.Println("Log into file : ", filePathCrawlFail)
		crawlFailLooger.SetOutput(fileCrawlFail)
	}
	// Set Log Level
	switch strings.ToLower(server.Config.GetString("SERVER_LOG_LEVEL")) {
	case "panic":
		logger.SetLevel(logrus.PanicLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "trace":
		logger.SetLevel(logrus.TraceLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}
}

// Println Function
func Println(level logLevel, label string, message interface{}) {
	// Make Sure Log Is Not Empty Variable
	if logger != nil {
		// Set Service Name Log Information
		service := strings.ToLower(server.Config.GetString("SERVER_NAME"))

		// Print Log Based On Log Level Type
		switch level {
		case "panic":
			logger.WithFields(logrus.Fields{
				"service": service,
				"label":   label,
			}).Panicln(message)
		case "fatal":
			logger.WithFields(logrus.Fields{
				"service": service,
				"label":   label,
			}).Fatalln(message)
		case "error":
			logger.WithFields(logrus.Fields{
				"service": service,
				"label":   label,
			}).Errorln(message)
		case "warn":
			logger.WithFields(logrus.Fields{
				"service": service,
				"label":   label,
			}).Warnln(message)
		case "debug":
			logger.WithFields(logrus.Fields{
				"service": service,
				"label":   label,
			}).Debug(message)
		case "trace":
			logger.WithFields(logrus.Fields{
				"service": service,
				"label":   label,
			}).Traceln(message)
		default:
			logger.WithFields(logrus.Fields{
				"service": service,
				"label":   label,
			}).Infoln(message)
		}
	}
}

func CrawlFailWrite(level logLevel, label string, message interface{}) {
	if crawlFailLooger != nil {
		// Set Service Name Log Information
		service := strings.ToLower(server.Config.GetString("SERVER_NAME"))

		// Print Log Based On Log Level Type
		crawlFailLooger.WithFields(logrus.Fields{
			"service": service,
			"label":   label,
		}).Infoln(message)
	}
}
