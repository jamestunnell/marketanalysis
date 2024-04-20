package commands

import (
	"github.com/rs/zerolog/log"
)

func Run(c Command) {
	log.Info().Interface("command", c).Msg("running command")

	err := c.Run()
	if err != nil {
		log.Error().Err(err).Msg("command failed")
	}
}
