package config

import (
	"github.com/Pickausernaame/chat-service/internal/validator"
)

type Config struct {
	Global  GlobalConfig  `toml:"global"`
	Log     LogConfig     `toml:"log"`
	Servers ServersConfig `toml:"servers"`
	Sentry  SentryConfig  `toml:"sentry"`
	Clients ClientsConfig `toml:"clients"`
}

func (c Config) Validate() error {
	return validator.Validator.Struct(c)
}

type GlobalConfig struct {
	Env     string `toml:"env" validate:"required,oneof=dev stage prod"`
	Version string `toml:"ver" validate:"semver,omitempty"`
}

type LogConfig struct {
	Level string `toml:"level" validate:"required,oneof=debug info warn error"`
}

type ServersConfig struct {
	Debug  DebugServerConfig   `toml:"debug"`
	Client ServersClientConfig `toml:"client"`
}

type DebugServerConfig struct {
	Addr string `toml:"addr" validate:"required,hostname_port"`
}

type ServersClientConfig struct {
	Addr           string                            `toml:"addr" validate:"required,hostname_port"`
	AllowsOrigins  []string                          `toml:"allow_origins" validate:"required,min=1"`
	RequiredAccess ServersClientRequiredAccessConfig `toml:"required_access"`
}

type ServersClientRequiredAccessConfig struct {
	Resource string `toml:"resource" validate:"required"`
	Role     string `toml:"role" validate:"required"`
}

type SentryConfig struct {
	DSN string `toml:"dsn" validate:"http_url,omitempty"`
}

type ClientsConfig struct {
	Keycloak KeycloakClientConfig `toml:"keycloak"`
}

type KeycloakClientConfig struct {
	BasePath     string `toml:"base_path" validate:"required,http_url"`
	Realm        string `toml:"realm" validate:"required"`
	ClientID     string `toml:"client_id" validate:"required"`
	ClientSecret string `toml:"client_secret" validate:"required"`
	DebugMode    bool   `toml:"debug_mode"`
}
