package main

//go:generate go run cmd/ast.go

import (
	"example/compilers/lex"
	"os"

	"github.com/rs/zerolog/log"
)

func main() {
	if len(os.Args) != 2 {
		log.Error().Int("num_args", len(os.Args)).Msg("Usage: cmd <file_path>")
		return
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Err(err).Msg("Failed to read file")
		return
	}

	// log.Trace().Msg(string(data))

	scanner := lex.NewScanner(string(data))
	tokens, err := scanner.Scan()

	if err != nil {
		log.Err(err).Msg("Failed to scan.")
	}

	for _, t := range tokens {
		log.Info().Msg(t.String())
	}
}
