package gzap

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	graylog "github.com/Devatoria/go-graylog"
)

const tlsTransport = "tls"

// Config is an interface representing all the logging configurations accessible
// via environment
type Config interface {
	enableJSONFormatter() bool
	getGraylogAppName() string
	getGraylogHandlerType() graylog.Transport
	getGraylogHost() string
	getGraylogPort() uint
	getGraylogLogLevel() uint
	getGraylogTLSTimeout() time.Duration
	getGraylogLogEnvName() string
	getGraylogSkipInsecureSkipVerify() bool
	getIsTestEnv() bool
	useTLS() bool
	useColoredConsolelogs() bool
}

// EnvConfig represents all the logger configurations available
// when instaniating a new Logger.
type EnvConfig struct{}

func (e *EnvConfig) enableJSONFormatter() bool {
	jsonFormatter := os.Getenv("ENABLE_DATADOG_JSON_FORMATTER")
	if jsonFormatter == "true" {
		return true
	}
	return false
}

func (e EnvConfig) getGraylogAppName() string {
	appName := os.Getenv("GRAYLOG_APP_NAME")
	if appName == "" {
		panic("GRAYLOG_APP_NAME env not set")
	}

	return appName
}

func (e *EnvConfig) getGraylogHandlerType() graylog.Transport {
	defaultHandlerType := tlsTransport
	handlerType := os.Getenv("GRAYLOG_HANDLER_TYPE")

	var transportType graylog.Transport
	if handlerType == tlsTransport {
		transportType = graylog.TCP
	}

	if graylog.Transport(handlerType) == graylog.UDP {
		transportType = graylog.UDP
	}

	// If no transport type is set use tls by default.
	if transportType == "" {
		transportType = graylog.Transport(defaultHandlerType)
	}

	return transportType
}

func (e *EnvConfig) getGraylogHost() string {
	graylogHost := os.Getenv("GRAYLOG_HOST")
	return graylogHost
}

func (e *EnvConfig) getGraylogPort() uint {
	portString := "12201"

	if e.getGraylogHandlerType() == graylog.UDP {
		portString = os.Getenv("GRAYLOG_UDP_PORT")
	}

	if e.getGraylogHandlerType() == graylog.TCP {
		portString = os.Getenv("GRAYLOG_TLS_PORT")
	}

	port, err := strconv.ParseUint(portString, 10, 32)
	if err != nil {
		panic(fmt.Errorf("could not properly parse Graylog port: %s", portString))
	}

	return uint(port)
}

func (e *EnvConfig) getGraylogLogLevel() uint {
	s := os.Getenv("GRAYLOG_LOG_LEVEL")

	level, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		panic(fmt.Errorf("could not properly parse Graylog LogLevel: %s", s))
	}

	return uint(level)
}

func (e *EnvConfig) getGraylogTLSTimeout() time.Duration {
	defaultTimeout := time.Second * 3

	timeoutString := os.Getenv("GRAYLOG_TLS_TIMEOUT_SECS")
	if timeoutString == "" {
		return defaultTimeout
	}

	timeoutSeconds, err := strconv.ParseInt(timeoutString, 10, 32)
	if err != nil {
		panic("invalid GRAYLOG_TLS_TIMEOUT_SECS could not parse int")
	}

	return time.Second * time.Duration(timeoutSeconds)
}

func (e *EnvConfig) getGraylogLogEnvName() string {
	envName := os.Getenv("GRAYLOG_ENV")
	if envName == "" {
		panic("GRAYLOG_ENV not set")
	}

	return envName
}

func (e *EnvConfig) getGraylogSkipInsecureSkipVerify() bool {
	skipInsecure := os.Getenv("GRAYLOG_SKIP_TLS_VERIFY")
	if skipInsecure == "true" {
		return true
	}

	return false
}

func (e *EnvConfig) getIsTestEnv() bool {
	// If we're running test return test logger env.
	if flag.Lookup("test.v") != nil {
		return true
	}

	return false
}

func (e *EnvConfig) useTLS() bool {
	handlerType := os.Getenv("GRAYLOG_HANDLER_TYPE")
	if handlerType == "" {
		panic("GRAYLOG_HANDLER_TYPE env not set")
	}

	if handlerType == tlsTransport {
		return true
	}

	return false
}

func (e *EnvConfig) useColoredConsolelogs() bool {
	envLevel := os.Getenv("THEMUSE_ENV_LEVEL")
	// If the env level is not set use colored logs.
	if envLevel == "0" {
		return true
	}

	return false
}

// CfgConfig implement Config interface from config struct
type CfgConfig struct {
	EnableJSONFormatter bool   `json:"enable_json_formatter" yaml:"enable_json_formatter"`
	AppName             string `json:"app_name" yaml:"app_name"`
	EnvName             string `json:"env_name" yaml:"env_name"`
	HanlderType         string `json:"handler_name" yaml:"handler_name"`
	Host                string `json:"host" yaml:"host"`
	UDPPort             uint   `json:"udp_port" yaml:"udp_port"`
	TLSPort             uint   `json:"tls_port" yaml:"tls_port"`
	TLSTimeoutSeconds   string `json:"tls_timeout_seconds" yaml:"tls_timeout_seconds"`
	LogLevel            uint   `json:"log_level" yaml:"log_level"`
}

func NewDefaultCfgConfig() *CfgConfig {
	cfg := &CfgConfig{
		EnableJSONFormatter: true,
		AppName:             "pracing",
		EnvName:             "pracing",
		HanlderType:         "udp",
		Host:                "127.0.0.1",
		UDPPort:             12001,
		TLSPort:             12001,
		TLSTimeoutSeconds:   "3",
		LogLevel:            4,
	}
	return cfg
}

func (e *CfgConfig) enableJSONFormatter() bool {
	return e.EnableJSONFormatter
}

func (e *CfgConfig) getGraylogAppName() string {
	return e.AppName
}

func (e *CfgConfig) getGraylogHandlerType() graylog.Transport {
	defaultHandlerType := tlsTransport
	handlerType := e.HanlderType

	var transportType graylog.Transport
	if handlerType == tlsTransport {
		transportType = graylog.TCP
	}

	if graylog.Transport(handlerType) == graylog.UDP {
		transportType = graylog.UDP
	}

	// If no transport type is set use tls by default.
	if transportType == "" {
		transportType = graylog.Transport(defaultHandlerType)
	}

	return transportType
}

func (e *CfgConfig) getGraylogHost() string {
	return e.Host
}

func (e *CfgConfig) getGraylogPort() uint {
	if e.getGraylogHandlerType() == graylog.UDP {
		return e.UDPPort
	}

	if e.getGraylogHandlerType() == graylog.TCP {
		return e.TLSPort
	}

	return 12001
}

func (e *CfgConfig) getGraylogLogLevel() uint {
	return e.LogLevel
}

func (e *CfgConfig) getGraylogTLSTimeout() time.Duration {
	defaultTimeout := time.Second * 3

	timeoutString := e.TLSTimeoutSeconds
	if timeoutString == "" {
		return defaultTimeout
	}

	timeoutSeconds, err := strconv.ParseInt(timeoutString, 10, 32)
	if err != nil {
		panic("invalid GRAYLOG_TLS_TIMEOUT_SECS could not parse int")
	}

	return time.Second * time.Duration(timeoutSeconds)
}

func (e *CfgConfig) getGraylogLogEnvName() string {
	return e.EnvName
}

func (e *CfgConfig) getGraylogSkipInsecureSkipVerify() bool {
	return false
}

func (e *CfgConfig) getIsTestEnv() bool {
	// If we're running test return test logger env.
	if flag.Lookup("test.v") != nil {
		return true
	}

	return false
}

func (e *CfgConfig) useTLS() bool {
	handlerType := e.HanlderType
	if handlerType == "" {
		panic("GRAYLOG_HANDLER_TYPE env not set")
	}

	if handlerType == tlsTransport {
		return true
	}

	return false
}

func (e *CfgConfig) useColoredConsolelogs() bool {
	return false
}
