package log

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

// ViperConfiguration is the configuration to use viper to load Configuration.
type ViperConfiguration struct {
	EnvPrefix  string `default:"GAP_LOG"`
	ConfigFile string
	ConfigName string   `default:"gaplog"`
	ConfigPath []string `default:"[\".\"]"`
	FileSystem afero.Fs `default:"-"`
}

// LoadFromViper loads the Configuration using viper. A default configuration is returned if it
// fails to load from viper.
func ConfigurationFromViper(v ViperConfiguration) (Configuration, error) {
	logViper := viper.New()
	if v.EnvPrefix != "" {
		logViper.AutomaticEnv()
		logViper.SetEnvPrefix(v.EnvPrefix)
	}
	if v.ConfigFile != "" {
		logViper.SetConfigFile(v.ConfigFile)
	}
	if v.ConfigName != "" {
		logViper.SetConfigName(v.ConfigName)
	}
	if v.ConfigPath != nil {
		for _, path := range v.ConfigPath {
			logViper.AddConfigPath(path)
		}
	}
	if v.FileSystem != nil {
		logViper.SetFs(v.FileSystem)
	}

	if err := logViper.ReadInConfig(); err != nil {
		return defaultConfig, err
	}

	config := defaultConfig
	if err := logViper.Unmarshal(&config); err != nil {
		return defaultConfig, err
	}
	if err := validator.New().Struct(config); err != nil {
		return defaultConfig, err
	}
	return config, nil
}
