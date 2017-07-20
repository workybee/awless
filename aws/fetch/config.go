package awsfetch

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/wallix/awless/logger"
)

type Config struct {
	Sess  *session.Session
	Log   *logger.Logger
	Extra map[string]interface{}
}

func (c *Config) getBoolDefaultTrue(key string) bool {
	if c.Extra == nil {
		return true
	}

	if b, ok := c.Extra[key].(bool); ok {
		return b
	}

	return true
}
