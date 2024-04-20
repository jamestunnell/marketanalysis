package commands

import (
	"github.com/rs/zerolog/log"
)

func InitAndRun(c Command) {
	if err := c.Init(); err != nil {
		log.Error().Err(err).Msg("failed to initialize command")

		return
	}

	log.Info().Interface("command", c).Msg("running command")

	if err := c.Run(); err != nil {
		log.Error().Err(err).Msg("command failed")
	}
}
