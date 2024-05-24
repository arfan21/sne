package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type config struct {
	HttpPort string `mapstructure:"HTTP_PORT"`
	Env      string `mapstructure:"ENV"`

	Redis     redis     `mapstructure:",squash"`
	Service   service   `mapstructure:",squash"`
	RateLimit rateLimit `mapstructure:",squash"`
	Database  database  `mapstructure:",squash"`
}

type service struct {
	Timeout                      int    `mapstructure:"SERVICE_TIMEOUT"`
	Name                         string `mapstructure:"SERVICE_NAME"`
	Version                      string `mapstructure:"SERVICE_VERSION"`
	ShortedURLCharLength         int    `mapstructure:"SERVICE_URL_SHORTED_CHAR_LENGTH"`
	BaseURL                      string `mapstructure:"SERVICE_BASE_URL"`
	MaxAttemptGenerateShortedURL int    `mapstructure:"SERVICE_MAX_ATTEMPT_GENERATE_SHORTED_URL"`
	ShortedURLTTL                int    `mapstructure:"SERVICE_URL_SHORTED_TTL"`
	ShortedURLTTLCache           int    `mapstructure:"SERVICE_URL_SHORTED_TTL_CACHE"`
}

type rateLimit struct {
	Max        int `mapstructure:"RATE_LIMIT_MAX"`
	Expiration int `mapstructure:"RATE_LIMIT_EXPIRATION"`
}

type redis struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     int    `mapstructure:"REDIS_PORT"`
	Username string `mapstructure:"REDIS_USERNAME"`
	Password string `mapstructure:"REDIS_PASSWORD"`
}

type database struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     int    `mapstructure:"DB_PORT"`
	Username string `mapstructure:"DB_USERNAME"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
}

func (c database) GetDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", c.Username, c.Password, c.Host, c.Port, c.Name)
}

var configInstance *config
var viperInstance *viper.Viper

func LoadConfig(filenames ...string) (*viper.Viper, error) {
	if viperInstance != nil {
		return viperInstance, nil
	}
	v := viper.New()
	if len(filenames) > 0 {
		// v.SetConfigName("app")
		v.SetConfigFile(filenames[0])
	} else {
		// check .env file exist
		if _, err := os.Stat(".env"); err == nil {
			v.SetConfigFile(".env")
		}
	}

	initDefaultValue(v)
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil && !strings.Contains(err.Error(), "Not Found in") {
		err = fmt.Errorf("error read config file: %s", err)
		return nil, err
	}

	viperInstance = v
	return viperInstance, nil
}

func ParseConfig(v *viper.Viper) (*config, error) {
	if configInstance != nil {
		return configInstance, nil
	}
	var c config
	var out map[string]interface{}
	err := mapstructure.Decode(&c, &out)
	if err != nil {
		err = fmt.Errorf("error decode config: %s", err)
		return nil, err
	}

	for key := range out {
		vKey := strings.ToLower(strings.ReplaceAll(key, ".", "_"))
		err = v.BindEnv(vKey, key)
		if err != nil {
			err = fmt.Errorf("error bind env: %s", err)
			return nil, err
		}
	}

	err = v.Unmarshal(&c)
	if err != nil {
		err = fmt.Errorf("error unmarshal config: %s", err)
		return nil, err
	}

	configInstance = &c
	return configInstance, nil
}

func GetConfig(filenames ...string) *config {
	if configInstance == nil {
		LoadConfig(filenames...)
		ParseConfig(viperInstance)
	}
	return configInstance
}

func GetViper(filenames ...string) *viper.Viper {
	if viperInstance == nil {
		LoadConfig(filenames...)
		ParseConfig(viperInstance)
	}
	return viperInstance
}

func initDefaultValue(v *viper.Viper) {
	v.SetDefault("HTTP_PORT", "8080")
	v.SetDefault("ENV", "dev")
	v.SetDefault("SERVICE_NAME", "url-shortener")
	v.SetDefault("SERVICE_TIMEOUT", 30)
	v.SetDefault("SERVICE_BASE_URL", "http://localhost:8080")
	v.SetDefault("SERVICE_URL_SHORTED_CHAR_LENGTH", 6)
	v.SetDefault("SERVICE_MAX_ATTEMPT_GENERATE_SHORTED_URL", 5)
	v.SetDefault("REDIS_URL_SHORTED_TTL", 5*365*24)

	v.SetDefault("RATE_LIMIT_MAX", 15)
	v.SetDefault("RATE_LIMIT_EXPIRATION", 30)
}
