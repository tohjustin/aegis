package config

import (
	"flag"
	"time"
)

const (
	portCfg                       = "port"
	readTimeoutCfg                = "read-timeout"
	writeTimeoutCfg               = "write-timeout"
	excludeCacheControlHeadersCfg = "exclude-cache-control-headers"
)

var (
	port                       *uint
	readTimeout                *uint
	writeTimeout               *uint
	excludeCacheControlHeaders *bool
)

// Flags adds flags related to the application to the given flagset.
func Flags(flags *flag.FlagSet) {
	port = flags.Uint(portCfg, 8080, "Port exposing badget service.")
	readTimeout = flags.Uint(readTimeoutCfg, 2000, "Maximum duration in milliseconds for reading the entire request, including the body.")
	writeTimeout = flags.Uint(writeTimeoutCfg, 2000, "Maximum duration in milliseconds before timing out writes of the response.")
	excludeCacheControlHeaders = flags.Bool(excludeCacheControlHeadersCfg, false, "Flag to exclude HTTP Cache-Control headers from responses.")
}

// Port returns the port number the service is listening to
func Port() uint {
	return *port
}

// ReadTimeout returns the service maximum duration for reading the entire request, including the body
func ReadTimeout() time.Duration {
	return time.Duration(*readTimeout) * time.Millisecond
}

// WriteTimeout returns the service maximum duration before timing out writes of the response
func WriteTimeout() time.Duration {
	return time.Duration(*writeTimeout) * time.Millisecond
}

// ExcludeCacheControlHeaders returns whether to include HTTP Cache-Control headers in responses or not
func ExcludeCacheControlHeaders() bool {
	return *excludeCacheControlHeaders
}
