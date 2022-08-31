package global

import (
	"github.com/rs/zerolog"
	"github.com/zdz1715/webhook/pkg/engine"
)

var (
	AccessLogger  zerolog.Logger
	WebhookLogger zerolog.Logger

	Engine = engine.New()
)
