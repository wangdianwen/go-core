package logger

import "github.com/astaxie/beego/logs"

// logger levels
const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

const (
	LevelInfo  = LevelInformational
	LevelTrace = LevelDebug
	LevelWarn  = LevelWarning
)

type BLogger struct {
	beeLogger *logs.BeeLogger
	FileName  string
}

type Content struct {
	Level   int
	Format  string
	V       []interface{}
	bLogger *BLogger
}

var ContentChan chan *Content
