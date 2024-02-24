package log

import (
	"io"
	"io/fs"
	"log/slog"
	"os"
	"time"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/lmittmann/tint"
)

type Configuration struct {
	Type        string                                          `toml:"type" mapstructure:"type" default:"console" validate:"oneof=console file"`
	Level       string                                          `toml:"level" mapstructure:"level" default:"info" validate:"oneof=trace debug info warn error"`
	AddSource   bool                                            `toml:"source" mapstructure:"source" default:"true"`
	ReplaceAttr func(groups []string, attr slog.Attr) slog.Attr `toml:"-" mapstructure:"-"`
	File        FileConfiguration                               `toml:",omitempty,squash" mapstructure:",squash"`
	Console     ConsoleConfiguration                            `toml:",omitempty,squash" mapstructure:",squash"`
}

type FileConfiguration struct {
	File       string      `toml:"file" mapstructure:"file" validate:"isdefault|filepath"`
	Permission fs.FileMode `toml:"permission" mapstructure:"permission" default:"0640"`
	Truncate   bool        `toml:"truncate" mapstructure:"truncate" default:"true"`
}

type ConsoleConfiguration struct {
	Handler    string `toml:"handler" mapstructure:"handler" default:"stderr" validate:"oneof=stderr stdout"`
	TimeFormat string `toml:"time" mapstructure:"time" default:"Kitchen" validate:"oneof=Layout RubyDate RFC822Z RFC1123Z RFC3339 Kitchen DateTime TimeOnly"`
	NoColor    bool   `toml:"nocolor" mapstructure:"nocolor"`
}

var (
	logLevelMaping = map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}

	timefmtMappping = map[string]string{
		"Layout":   time.Layout,
		"RubyDate": time.RubyDate,
		"RFC822Z":  time.RFC822Z,
		"RFC1123Z": time.RFC1123Z,
		"RFC3339":  time.RFC3339,
		"Kitchen":  time.Kitchen,
		"DateTime": time.DateTime,
		"TimeOnly": time.TimeOnly,
	}

	defaultConfig Configuration
)

func init() {
	if err := defaults.Set(&defaultConfig); err != nil {
		panic(err)
	}
	if err := validator.New().Struct(defaultConfig); err != nil {
		panic(err)
	}

	var v ViperConfiguration
	if err := defaults.Set(&v); err != nil {
		panic(err)
	}
	config, _ := ConfigurationFromViper(v)
	slog.SetDefault(Logger(config))
}

func Logger(config Configuration) *slog.Logger {
	var handler slog.Handler
	switch config.Type {
	case "console":
		var timeFormat string
		if config.Console.TimeFormat != "" {
			timeFormat = timefmtMappping[config.Console.TimeFormat]
		}
		var writer io.Writer
		if config.Console.Handler == "stdout" {
			writer = stdout()
		} else {
			writer = stderr()
		}
		handler = tint.NewHandler(writer, &tint.Options{
			AddSource:   config.AddSource,
			Level:       logLevelMaping[config.Level],
			ReplaceAttr: config.ReplaceAttr,
			TimeFormat:  timeFormat,
			NoColor:     config.Console.NoColor,
		})
	case "file":
		flag := os.O_CREATE | os.O_APPEND | os.O_WRONLY
		if config.File.Truncate {
			flag |= os.O_TRUNC
		}
		file, err := os.OpenFile(config.File.File, flag, config.File.Permission)
		if err != nil {
			panic(err)
		}
		handler = slog.NewTextHandler(file, &slog.HandlerOptions{
			AddSource:   config.AddSource,
			Level:       logLevelMaping[config.Level],
			ReplaceAttr: config.ReplaceAttr,
		})
	}
	return slog.New(handler)
}
