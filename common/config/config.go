package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ApplicationEnv        ApplicationEnv `envconfig:"APPLICATION_ENV" default:"development"`
	ApplicationVersion    string         `envconfig:"APPLICATION_VERSION" default:"latest"`
	Host                  string         `envconfig:"HOST" default:"localhost"`
	Port                  int            `envconfig:"PORT" default:"3000"`
	Path                  string         `envconfig:"PATH" default:""`
	VideoPath             string         `envconfig:"VIDEOPATH" default:""`
	Nobrowser             bool           `envconfig:"NOBROWSER" default:""`
	PrintVersion          bool           `envconfig:"PRINTVERSION" default:""`
	Enablep2p             bool           `envconfig:"ENABLEP2P" default:""`
	P2PPort               int            `envconfig:"P2PPORT" default:""`
	P2PEnable             int            `envconfig:"P2P_ENABLE" default:""`
	PrintHelp             bool           `envconfig:"PRINTHELP" default:""`
	SessionSecret         string         `envconfig:"SESSION_SECRET" default:""`
	SQLitePath            string         `envconfig:"SQLITE_PATH" default:""`
	UploadPath            string         `envconfig:"UPLOAD_PATH" default:""`
	LogDir                string         `envconfig:"LOG_DIR" default:""`
	GinMode               string         `envconfig:"GIN_MODE" default:""`
	RedisConnectionString string         `envconfig:"REDIS_CONN_STRING" default:""`
	CloudinaryURL         string         `envconfig:"CLOUDINARY_URL" default:""`
}

func FromEnv(conf *Config) error {
	return envconfig.Process("", conf)
}

func (c *Config) ApplyConfig(procs ...func(c *Config) error) error {
	for _, p := range procs {
		if err := p(c); err != nil {
			return fmt.Errorf("failed to modify config: %w", err)
		}
	}
	return nil
}
