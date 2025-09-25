package services

import (
	config "notification/internal/config"
	services "notification/pkg/services"
)

// registry services for modules business logic

type Factory struct {
	NewService func(cfg *config.Config) services.Notifier[any]
	Converter  func([]byte) (any, error)
}

var registry = map[string]Factory{}

func Register(name string, f Factory) {
	registry[name] = f
}

func Get(name string) (Factory, bool) {
	f, ok := registry[name]
	return f, ok
}
