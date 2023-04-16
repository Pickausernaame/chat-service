package config

import (
	"fmt"
)

type Config struct {
	Global  GlobalConfig  `toml:"global"`
	Log     LogConfig     `toml:"log"`
	Servers ServersConfig `toml:"servers"`
}

func (c Config) Validate() error {
	if err := c.Global.Validate(); err != nil {
		return fmt.Errorf("global config validation error: %v", err)
	}

	if err := c.Log.Validate(); err != nil {
		return fmt.Errorf("log config validation error: %v", err)
	}

	if err := c.Servers.Debug.Validate(); err != nil {
		return fmt.Errorf("debug config validation error: %v", err)
	}

	return nil
}

//go:generate options-gen -out-filename=global_config_gen.go -from-struct=GlobalConfig
type GlobalConfig struct {
	Env string `toml:"env" validate:"required,oneof=dev stage prod"`
}

//go:generate options-gen -out-filename=log_config_gen.go -from-struct=LogConfig
type LogConfig struct {
	Level string `toml:"level" validate:"required,oneof=debug info warn error"`
}

type ServersConfig struct {
	Debug DebugServerConfig `toml:"debug"`
}

//go:generate options-gen -out-filename=debug_server_config_gen.go -from-struct=DebugServerConfig
type DebugServerConfig struct {
	Addr string `toml:"addr" validate:"required,hostname_port"`
}
