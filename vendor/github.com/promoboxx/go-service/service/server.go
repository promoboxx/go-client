package service

import (
	"github.com/husobee/vestigo"
)

// Server can GetRouter
type Server interface {
	GetRouter() *vestigo.Router
	GetEnvironment() string
	GetServiceName() string
}

type ListenerConfig struct {
	Port int
	Addr string
}
