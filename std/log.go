package std

import (
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
	"strings"
	// "syscall"
	"time"

	"github.com/evalphobia/logrus_sentry"
	"github.com/sirupsen/logrus"
	stdnet "github.com/PKUJohnson/solar/std/toolkit/net"
)

// LogFields indicates the log's tags
type LogFields map[string]interface{}

var (
	logFile *os.File
	loc     *time.Location
)

const (
	// TagTopic flags the topic
	TagTopic = "topic"

	// TopicCodeTrace traces the running of code
	TopicCodeTrace = "code_trace"

	// TopicBugReport indicates the bug report topic
	TopicBugReport = "bug_report"

	// TopicCrash indicates the program's panics
	TopicCrash = "crash"

	// TopicUserActivity indicates the user activity like web access, user login/logout
	TopicUserActivity = "user_activity"

	// TagCategory tags the log category
	TagCategory = "category"

	// TagError tags the error category
	TagError = "error"

	// CategoryRPC indicates the rpc category
	CategoryRPC = "rpc"

	// CategoryRedis indicates the redis category
	CategoryRedis = "redis"

	// CategoryMySQL indicates the MySQL category
	CategoryMySQL = "mysql"

	// CategoryElasticsearch indicates the Elasticsearch category
	CategoryElasticsearch = "elasticsearch"
)

// InitLog initializes the logger
func InitLog(conf ConfigLog) {
	logrus.SetLevel(logrus.Level(conf.Level))

	loc, _ = time.LoadLocation("Asia/Shanghai")

	// set sentry
	if conf.SentryDSN != "" {
		hook, err := logrus_sentry.NewSentryHook(
			conf.SentryDSN,
			[]logrus.Level{
				logrus.PanicLevel,
				logrus.FatalLevel,
				logrus.ErrorLevel,
			})

		hook.SetEnvironment(os.Getenv("CONFIGOR_ENV"))

		if err == nil {
			logrus.AddHook(hook)
		}
	}

	switch {
	case os.Getenv("AIRFLOW") == "1":
		// output directly to the stdout & stderr
	case conf.OutputDest == "elasticsearch":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	case conf.OutputDest == "file":
		fullpath,err := filepath.Abs(conf.Path)
		if err != nil {
			logrus.Errorf("Wrong log path, %v", err)
		}

		err = os.MkdirAll(fullpath, os.ModeDir)
		if err != nil {
			logrus.Errorf("Failed to create log path folder, %v", err)
		}

		rotationTime,_ := time.ParseDuration(conf.RotationDuration)
		baseLogPath := path.Join(fullpath, conf.FileName)

		configLocalFilesystemLogger(baseLogPath, rotationTime, conf.RotationCount)
	}
}

func configLocalFilesystemLogger(fullLogPath string,  rotationTime time.Duration, rotationCount uint) {

	writer, err := rotatelogs.New(
		fullLogPath+"-%Y%m%d.log",
		rotatelogs.WithLinkName(fullLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
		rotatelogs.WithRotationCount(rotationCount),
	)
	if err != nil {
		logrus.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.TextFormatter{})
	logrus.AddHook(lfHook)
}

func logName(path string, ip string, time time.Time) string {
	year, month, day := time.Date()

	// replace all dots
	path = strings.Replace(path, ".", "_", -1)
	ip = strings.Replace(ip, ".", "_", -1)
	return fmt.Sprintf("%s.%04d%02d%02d.%s.log", path, year, int(month), day, ip)
}

func logRotate() {
	if logFile == nil {
		return
	}

	parts := strings.Split(logFile.Name(), ".")
	path := parts[0]
	date := parts[1]
	ip := parts[2]

	now := time.Now()
	if loc == nil {
		LogErrorc("timezone", nil, "fail to load Asia/Shanghai")
	} else {
		now = now.In(loc)
	}

	year, month, day := now.Date()
	if date >= fmt.Sprintf("%04d%02d%02d", year, int(month), day) {
		return
	}

	fmt.Println("redirect to: ", logName(path, ip, now))

	// redirect all stdout & stderr to file
	curFile := logFile
	defer curFile.Close()
	redirect(logName(path, ip, now))
}

func redirect(fullPath string) {
	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666|os.ModeSticky)
	if err == nil {
		logFile = file
		logrus.SetOutput(file)
		// syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd()))
		// syscall.Dup2(int(file.Fd()), int(os.Stdout.Fd()))
	} else {
		panic("log file open error: " + err.Error())
	}
}

func getIP(CIDRs []string) string {
	var ip string
	var err error
	if len(CIDRs) > 0 {
		ip, err = stdnet.LocalIPAddrWithin(CIDRs)
		if err != nil {
			return ""
		}
	} else {
		ip, err = stdnet.LocalIPAddr()
		if err != nil {
			return ""
		}
	}
	return ip
}

// LogInfo records Info level information which helps trace the running of program and
// moreover the production infos
func LogInfo(fields LogFields, message string) {
	logrus.WithFields(logrus.Fields{
		TagTopic: TopicCodeTrace,
	}).WithFields(map[string]interface{}(fields)).Info(message)
}

// LogInfoc records the running infos
func LogInfoc(category string, message string) {
	logrus.WithFields(logrus.Fields{
		TagTopic:    TopicCodeTrace,
		TagCategory: category,
	}).Info(message)
}

// LogWarn records the warnings which are expected to be removed, but not influence the
// running of the program
func LogWarn(fields LogFields, message string) {
	logrus.WithFields(logrus.Fields{
		TagTopic: TopicBugReport,
	}).WithFields(map[string]interface{}(fields)).Warn(message)
}

// LogWarnc records the running warnings which are expected to be noticed
func LogWarnc(category string, err error, message string) {
	logrus.WithFields(logrus.Fields{
		TagTopic:    TopicBugReport,
		TagCategory: category,
		TagError:    err,
	}).Warn(message)
}

// LogError records the running errors which are expected to be solved soon
func LogError(fields LogFields, message string) {
	logrus.WithFields(logrus.Fields{
		TagTopic: TopicBugReport,
	}).WithFields(map[string]interface{}(fields)).Error(message)
}

// LogErrorc records the running errors which are expected to be solved soon
func LogErrorc(category string, err error, message string) {
	logrus.WithFields(logrus.Fields{
		TagTopic:    TopicBugReport,
		TagCategory: category,
		TagError:    err,
	}).Error(message)
}

// LogPanic records the running errors which are expected to be severe soon
func LogPanic(fields LogFields, message string) {
	logrus.WithFields(logrus.Fields{
		TagTopic: TopicBugReport,
	}).WithFields(map[string]interface{}(fields)).Panic(message)
}

// LogPanicc records the running errors which are expected to be severe soon
func LogPanicc(category string, err error, message string) {
	logrus.WithFields(logrus.Fields{
		TagTopic:    TopicBugReport,
		TagCategory: category,
		TagError:    err,
	}).Panic(message)
}

// LogInfoLn records Info level information which helps trace the running of program and
// moreover the production infos
func LogInfoLn(args ...interface{}) {
	logrus.Infoln(args)
}

// LogWarnLn records the program warning
func LogWarnLn(args ...interface{}) {
	logrus.WithFields(logrus.Fields{
		TagTopic: TopicCodeTrace,
	}).Warnln(args)
}

// LogErrorLn records the program error, go to fix it!
func LogErrorLn(args ...interface{}) {
	logrus.WithFields(logrus.Fields{
		TagTopic: TopicBugReport,
	}).Errorln(args)
}

// LogFatalLn records the program fatal error, developer should follow immediately
func LogFatalLn(args ...interface{}) {
	logrus.WithFields(logrus.Fields{
		TagTopic: TopicBugReport,
	}).Fatalln(args)
}

// LogPanicLn records the program fatal error, developer should fix otherwise the company dies
func LogPanicLn(args ...interface{}) {
	logrus.WithFields(logrus.Fields{
		TagTopic: TopicBugReport,
	}).Panicln(args)
}

// LogDebugLn records debug information which helps trace the running of program
func LogDebugLn(args ...interface{}) {
	logrus.Debugln(args)
}

// LogDebugc records the running infos
func LogDebugc(category string, message string) {
	logrus.WithFields(logrus.Fields{
		TagTopic:    TopicCodeTrace,
		TagCategory: category,
	}).Debug(message)
}

// LogUserActivity records user activity, like user access page, login/logout
func LogUserActivity(fields LogFields, message string) {
	logrus.WithFields(logrus.Fields{
		TagTopic: TopicUserActivity,
	}).WithFields(map[string]interface{}(fields)).Infoln(message)
}

// LogRecover records when program crashes
func LogRecover(e interface{}) {
	logrus.WithFields(logrus.Fields{
		TagTopic:     TopicCrash,
		"error":      e,
		"stacktrace": string(debug.Stack()),
	}).Errorln("Recovered panic")
}
