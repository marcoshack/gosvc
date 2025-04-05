package config

type ServiceConfig interface {
	Validate() error
	IsValid() bool
	GetLogLevel() string
	GetLogFileName() string
}

type DefaultServiceConfig struct {
}

func (c DefaultServiceConfig) Validate() error {
	return nil
}

func (c DefaultServiceConfig) IsValid() bool {
	return true
}

func (c DefaultServiceConfig) GetLogLevel() string {
	return "debug"
}

func (c DefaultServiceConfig) GetLogFileName() string {
	return ""
}
