package server

// Server is transport server
type Server interface {
	Start() error
	Shutdown() error
}

// HTTP http config
type HTTP struct {
	Addr string `json:"addr" mapstructure:"addr"`
}
