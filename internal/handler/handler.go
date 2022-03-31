package handler

import (
	"github.com/KaiserWerk/envy/internal/configuration"
	"github.com/KaiserWerk/envy/internal/logging"
)

type Base struct {
	Config *configuration.AppConfig
	Logger *logging.Logger
}
