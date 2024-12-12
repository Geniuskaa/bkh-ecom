package config

import (
	"bkh-ecom/internal/logger"
	"context"

	"github.com/spf13/viper"
)

type App struct {
	Port int
	Host string
}

type DB struct {
	Host             string
	Name             string
	User             string
	Pass             string
	Port             int
	PoolSize         int `mapstructure:"pool_size"`
	ConnDialTimeout  int `mapstructure:"conn_dial_timeout"`
	ConnReadTimeout  int `mapstructure:"conn_read_timeout"`
	ConnWriteTimeout int `mapstructure:"conn_write_timeout"`
}

type Entity struct {
	App `mapstructure:"app"`
	DB  DB `mapstructure:"db"`
}

func NewConfig(ctx context.Context) (*Entity, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.AddConfigPath("./configs")
	v.AddConfigPath("/app/configs")

	v.SetConfigType("yaml")
	v.AutomaticEnv()
	err := v.ReadInConfig()
	if err != nil {
		logger.ErrorKV(ctx, logger.Data{
			Msg:   "Config read",
			Panic: err,
		})
	}

	config := Entity{}
	err = v.Unmarshal(&config)
	if err != nil {
		logger.ErrorKV(ctx, logger.Data{
			Msg:   "Config read",
			Panic: err,
		})
	}

	return &config, nil
}
