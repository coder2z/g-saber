package xlog

import (
	"errors"
	"fmt"
	cfg "github.com/coder2z/g-saber/xcfg"
	"github.com/coder2z/g-saber/xcolor"
	"github.com/coder2z/g-saber/xconsole"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type options struct {
	Dir                 string        `mapStructure:"dir"`                   //Dir 日志输出目录
	Name                string        `mapStructure:"name"`                  //Name 日志文件名称
	Level               string        `mapStructure:"level"`                 //Level 日志初始等级
	AddCaller           bool          `mapStructure:"add_caller"`            //是否添加调用者信息
	MaxSize             int           `mapStructure:"max_size"`              //日志输出文件最大长度，超过改值则截断，默认500M
	MaxAge              int           `mapStructure:"max_age"`               //日志存储最大时间，默认最大保存天数为7天
	MaxBackup           int           `mapStructure:"max_backup"`            //日志存储最大数量，默认最大保存文件个数为10个
	Interval            time.Duration `mapStructure:"interval"`              //日志轮转时间，默认1天
	CallerSkip          int           `mapStructure:"caller_skip"`           //调用堆栈
	Async               bool          `mapStructure:"async"`                 //是否异步，默认异步
	ConfigKey           string        `mapStructure:"config_key"`            //config key
	Debug               bool          `mapStructure:"debug"`                 //debug
	FlushBufferSize     int           `mapStructure:"flush_buffer_size"`     //缓冲大小，默认256 * 1024B
	FlushBufferInterval time.Duration `mapStructure:"flush_buffer_interval"` //缓冲时间，默认30秒

	fields        []zap.Field            `mapStructure:"fields"` //日志初始化字段
	encoderConfig *zapcore.EncoderConfig `mapStructure:"encoder_config"`
	core          zapcore.Core           `mapStructure:"core"`
}

// Filename ...
func (o *options) Filename() string {
	return fmt.Sprintf("%s/%s", o.Dir, o.Name)
}

// RawConfig ...
func RawConfig(key string) *options {
	var config = defaultConfig()
	if err := cfg.UnmarshalKey(key, &config); err != nil {
		if errors.Is(err, cfg.ErrInvalidKey) {
			xconsole.Blue("xlog use default config")
		} else {
			panic(err)
		}
	}
	config.ConfigKey = key
	return config
}

// StdConfig xlog
func StdConfig(name ...string) *options {
	if len(name) == 0 {
		return RawConfig("xlog")
	}
	return RawConfig("xlog." + name[0])
}

// DefaultConfig ...
func defaultConfig() *options {
	return &options{
		Dir:                 ".",
		Name:                "log.log",
		Level:               "info",
		AddCaller:           true,
		MaxSize:             500, // 500M
		MaxAge:              1,   // 1 day
		MaxBackup:           10,  // 10 backup
		Interval:            24 * time.Hour,
		CallerSkip:          2,
		Async:               false,
		ConfigKey:           "xlog",
		Debug:               true,
		FlushBufferSize:     defaultBufferSize,
		FlushBufferInterval: defaultFlushInterval,
		fields:              nil,
		encoderConfig:       DefaultZapConfig(),
		core:                nil,
	}
}

// Build ...
func (o *options) Build(options ...Option) *Logger {
	for _, option := range options {
		option(o)
	}
	if o.encoderConfig == nil {
		o.encoderConfig = DefaultZapConfig()
	}
	if o.Debug {
		o.encoderConfig.EncodeLevel = DebugEncodeLevel
		o.encoderConfig.EncodeTime = timeDebugEncoder
	}
	logger := newLogger(o)
	if o.ConfigKey != "" {
		logger.AutoLevel(o.ConfigKey + ".level")
	}
	return logger
}

func DefaultZapConfig() *zapcore.EncoderConfig {
	return &zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "lv",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func DebugEncodeLevel(lv zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var colorize = xcolor.Red
	switch lv {
	case zapcore.DebugLevel:
		colorize = xcolor.Blue
	case zapcore.InfoLevel:
		colorize = xcolor.Green
	case zapcore.WarnLevel:
		colorize = xcolor.Yellow
	case zapcore.ErrorLevel, zap.PanicLevel, zap.DPanicLevel, zap.FatalLevel:
		colorize = xcolor.Red
	default:
	}
	enc.AppendString(colorize(lv.CapitalString()))
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(t.Unix())
}

func timeDebugEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// Option 可选项
type Option func(c *options)

// WithFileName 设置文件名
func WithFileName(name string) Option {
	return func(c *options) {
		c.Name = name
	}
}

// WithDebug 设置在命令行显示
func WithDebug(debug bool) Option {
	return func(c *options) {
		c.Debug = debug
	}
}

// WithLevel 设置级别
func WithLevel(level string) Option {
	return func(c *options) {
		c.Level = level
	}
}

// WithEnableAsync 是否异步执行，默认异步
func WithEnableAsync(enableAsync bool) Option {
	return func(c *options) {
		c.Async = enableAsync
	}
}

// WithEnableAddCaller 是否添加行号，默认不添加行号
func WithEnableAddCaller(enableAddCaller bool) Option {
	return func(c *options) {
		c.AddCaller = enableAddCaller
	}
}

func WithFields(f []Field) Option {
	return func(c *options) {
		c.fields = f
	}
}
