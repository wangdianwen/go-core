package logger

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/wangdianwen/go-core/middlewares/email/models"
	"github.com/wangdianwen/go-core/utils"
	"runtime"
	"strings"
)

func init() {
	ContentChan = make(chan *Content, 1e6)
	go func() {
		for {
			select {
			case content := <-ContentChan:
				bErr := content.bLogger.doLogger(content)
				if bErr != nil {
					fmt.Println("bLogger ERROR:" + bErr.String())
				}
				// sending email by critical errors
				if content.Level == LevelCritical {
					_ = models.SysAlert(content.bLogger.FileName, fmt.Sprintf(content.Format, content.V))
				}
			}
		}
	}()
}

func (b *BLogger) write(content *Content) *utils.BError {
	if b.beeLogger == nil {
		err := errors.New("invalid logger, please initialize logger first")
		return &utils.BError{Error: err, Message: err.Error(), Code: utils.ErrorMWLogger}
	}
	switch content.Level {
	case LevelEmergency:
		b.beeLogger.Emergency(content.Format, content.V...)
	case LevelAlert:
		b.beeLogger.Alert(content.Format, content.V...)
	case LevelCritical:
		b.beeLogger.Critical(content.Format, content.V...)
	case LevelError:
		b.beeLogger.Error(content.Format, content.V...)
	case LevelWarn:
		b.beeLogger.Warn(content.Format, content.V...)
	case LevelNotice:
		b.beeLogger.Notice(content.Format, content.V...)
	case LevelInfo:
		b.beeLogger.Info(content.Format, content.V...)
	case LevelDebug:
		b.beeLogger.Debug(content.Format, content.V...)
	default:
		err := errors.New("invalid logger, please initialize logger first")
		return &utils.BError{Error: err, Message: err.Error(), Code: utils.ErrorMWLogger}
	}
	return nil
}

func (b *BLogger) doLogger(content *Content) *utils.BError {
	if b.FileName == "" {
		err := errors.New("invalid logger file name")
		return &utils.BError{Error: err, Message: err.Error(), Code: utils.ErrorMWLogger}
	}
	Log := logs.NewLogger()
	loggerRoot := beego.AppConfig.String("LoggerRoot")
	// default log director
	if loggerRoot == "" {
		loggerRoot = "./logs/"
	}
	filename := loggerRoot + b.FileName
	config := fmt.Sprintf(`{"filename": "%s", "perm": "0777", "maxdays": 3}`, filename)
	err := Log.SetLogger("file", config)
	if err != nil {
		return &utils.BError{Error: err, Message: err.Error(), Code: utils.ErrorService}
	}
	b.beeLogger = Log
	defer b.beeLogger.Close()
	return b.write(content)
}

func (b *BLogger) writeToChan(level int, format string, v ...interface{}) *utils.BError {
	if b.FileName == "" {
		err := errors.New("invalid logger, please initialize logger first")
		return &utils.BError{Error: err, Message: err.Error(), Code: utils.ErrorMWLogger}
	}
	_, file, line, _ := runtime.Caller(2)
	format = fmt.Sprintf("[%s:%d]\t%s", file, line, format)
	tmp := Content{Level: level, Format: format, V: v, bLogger: b}
	ContentChan <- &tmp
	return nil
}

func NewLogger(filename string) (*BLogger, *utils.BError) {
	p := new(BLogger)
	if filename == "" {
		err := errors.New("invalid logger file name")
		return nil, &utils.BError{Error: err, Message: err.Error(), Code: utils.ErrorParameter}
	}
	filename = strings.TrimSpace(filename)
	filename = strings.Trim(filename, "/")
	p.FileName = filename
	return p, nil
}

func (b *BLogger) Emergency(format string, v ...interface{}) *utils.BError {
	return b.writeToChan(LevelEmergency, format, v...)
}

func (b *BLogger) Alert(format string, v ...interface{}) *utils.BError {
	return b.writeToChan(LevelAlert, format, v...)
}

func (b *BLogger) Critical(format string, v ...interface{}) *utils.BError {
	return b.writeToChan(LevelCritical, format, v...)
}

func (b *BLogger) Error(format string, v ...interface{}) *utils.BError {
	return b.writeToChan(LevelError, format, v...)
}

func (b *BLogger) Warn(format string, v ...interface{}) *utils.BError {
	return b.writeToChan(LevelWarn, format, v...)
}

func (b *BLogger) Warning(format string, v ...interface{}) *utils.BError {
	return b.writeToChan(LevelWarn, format, v...)
}

func (b *BLogger) Notice(format string, v ...interface{}) *utils.BError {
	return b.writeToChan(LevelNotice, format, v...)
}

func (b *BLogger) Info(format string, v ...interface{}) *utils.BError {
	return b.writeToChan(LevelInfo, format, v...)
}

func (b *BLogger) Debug(format string, v ...interface{}) *utils.BError {
	return b.writeToChan(LevelDebug, format, v...)
}
