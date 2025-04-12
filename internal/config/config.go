package config

type ServiceConfig interface {
	Validate() error
	IsValid() bool
	GetLogLevel() string
	GetLogFileName() string
}

type SampleServiceConfig struct {
	Host string
	Port int
}

func (c SampleServiceConfig) Validate() error {
	return nil
}

func (c SampleServiceConfig) IsValid() bool {
	return c.Validate() == nil
}

func (c SampleServiceConfig) GetLogLevel() string {
	return "debug"
}

func (c SampleServiceConfig) GetLogFileName() string {
	return ""
}
