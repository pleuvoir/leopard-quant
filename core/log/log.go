package log

import (
	"bytes"
	"fmt"
	"github.com/gookit/color"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"leopard-quant/core/config"
	"leopard-quant/util"
	"os"
	"path"
	"path/filepath"
	"time"
)

var log *logrus.Logger

func IsTraceEnabled() bool {
	return log.GetLevel() <= logrus.TraceLevel
}

func IsDebugEnabled() bool {
	return log.GetLevel() <= logrus.DebugLevel
}

func IsInfoEnabled() bool {
	return log.GetLevel() <= logrus.InfoLevel
}

func IsWarnEnabled() bool {
	return log.GetLevel() <= logrus.WarnLevel
}

func IsErrorEnabled() bool {
	return log.GetLevel() <= logrus.ErrorLevel
}

func IsFatalEnabled() bool {
	return log.GetLevel() <= logrus.FatalLevel
}

func IsPanicEnabled() bool {
	return log.GetLevel() <= logrus.PanicLevel
}

func Print(args ...interface{}) {
	log.Print(args...)
}

func Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func Println(args ...interface{}) {
	log.Println(args...)
}

func Trace(args ...interface{}) {
	log.Trace(args...)
}

func Tracef(format string, args ...interface{}) {
	log.Tracef(format, args...)
}

func Traceln(args ...interface{}) {
	log.Traceln(args...)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Debugln(args ...interface{}) {
	log.Debugln(args...)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Infoln(args ...interface{}) {
	log.Infoln(args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Warnln(args ...interface{}) {
	log.Warnln(args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Errorln(args ...interface{}) {
	log.Errorln(args...)
}

func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func Fatalln(args ...interface{}) {
	log.Fatalln(args...)
}

func Panic(args ...interface{}) {
	log.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

func Panicln(args ...interface{}) {
	log.Panicln(args...)
}

// InitLog 初始化日志，初始化失败时会抛出异常
// param
//
//	config: 系统配置
func InitLog(config *config.ApplicationConfig) {
	if log != nil {
		log.Warn("日志系统已初始化，请不要重复调用InitLog")
		return
	}
	var maxAgeDuration time.Duration
	var rotationTimeDuration time.Duration
	maxAge := config.Logger.MaxAge
	if maxAge == "" {
		maxAgeDuration = time.Hour * 24 * 60 //60天
	} else {
		duration, err := time.ParseDuration(maxAge)
		if err != nil {
			panic(fmt.Sprintf("init log maxAge fail: %s, %s", maxAge, err.Error()))
		}
		maxAgeDuration = duration
	}
	rotationTime := config.Logger.RotationTime
	if rotationTime == "" {
		rotationTimeDuration = time.Hour * 24 //1天
	} else {
		duration, err := time.ParseDuration(rotationTime)
		if err != nil {
			panic(fmt.Sprintf("init log rotationTime fail: %s, %s", rotationTime, err.Error()))
		}
		rotationTimeDuration = duration
	}
	configLocalFilesystemLogger(config.Logger.Level,
		config.Logger.Path,
		config.Logger.Filename,
		maxAgeDuration,
		rotationTimeDuration)
}

// config logrus log to local filesystem, with file rotation
func configLocalFilesystemLogger(level string, logPath string, logFileName string,
	maxAge time.Duration, rotationTime time.Duration) {
	color.Printf("<light_green>ready to init log:</> level:%s, logPath:%s, logFilename:%s, maxAge:%s, rotationTime:%s \n",
		level, logPath, logFileName, maxAge.String(), rotationTime.String())

	if !util.IsExists(logPath) {
		err := os.MkdirAll(logPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	log = logrus.New()

	lv, err := logrus.ParseLevel(level)
	if err != nil {
		log.Panicf("cannot parse log level: %s, %+v", level, errors.WithStack(err))
	}
	log.SetLevel(lv)

	baseLogPath := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		log.Panicf("config local file system logger error. %+v", errors.WithStack(err))
	}

	log.SetReportCaller(false)

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.TraceLevel: writer,
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &ResetFormatter{})
	log.AddHook(lfHook)
}

type ResetFormatter struct{}

func (m *ResetFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var newLog string
	//HasCaller()为true才会有调用信息
	if entry.HasCaller() {
		fName := filepath.Base(entry.Caller.File)
		newLog = fmt.Sprintf("[%s] [%s] [%s:%d %s] %s\n",
			timestamp, entry.Level, fName, entry.Caller.Line, entry.Caller.Function, entry.Message)
	} else {
		newLog = fmt.Sprintf("[%s] [%s] %s\n", timestamp, entry.Level, entry.Message)
	}
	b.WriteString(newLog)
	return b.Bytes(), nil
}
