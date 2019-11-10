package server

type Server interface {
	Init(...Option) error
	Start() error
	Stop() error
}
