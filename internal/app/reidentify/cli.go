// Copyright (C) 2022 CGI France
//
// This file is part of SIGO.
//
// SIGO is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// SIGO is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with SIGO.  If not, see <http://www.gnu.org/licenses/>.

package reidentify

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	over "github.com/Trendyol/overlog"
	"github.com/cgi-fr/sigo/internal/infra"
	"github.com/cgi-fr/sigo/pkg/reidentification"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type logs struct {
	name      string
	version   string
	commit    string
	buildDate string
	builtBy   string
	verbosity string
	debug     bool
	jsonlog   bool
	colormode string
}

type reid struct {
	qi         []string
	sensitive  []string
	original   string
	anonymized string
	threshold  float32
}

// NewCommand implements the cli reidentification command.
func NewCommand(fullName string, err *os.File, out *os.File, in *os.File) *cobra.Command {
	var data reid

	var logs logs

	//nolint: exhaustivestruct
	cmd := &cobra.Command{
		Use:   "reidentification",
		Short: "Re-identify anonymized data from an original dataset",
		Long:  "",
		Example: fmt.Sprintf("  %[1]s reidentification -q x,y -s z --load-original o --load-anonymized a"+
			"--threshold 0.8", fullName),
		PreRun: func(cmd *cobra.Command, args []string) {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

			initLog(logs)

			log.Info().
				Str("openData", data.original).
				Str("anonymizedData", data.anonymized).
				Float32("threshold", data.threshold).
				Strs("Quasi-Identifiers", data.qi).
				Strs("Sensitive", data.sensitive).
				Msg("Reidentification mode")
		},
		Run: func(cmd *cobra.Command, args []string) {
			sink := infra.NewJSONLineSink(out)
			original := reidentification.LoadFile(data.original, data.qi, data.sensitive)
			anonymized := reidentification.LoadFile(data.anonymized, data.qi, data.sensitive)

			// Reidentification
			err := reidentification.ReIdentify(original, anonymized,
				reidentification.NewIdentifier("euclidean", data.threshold), sink)
			if err != nil {
				log.Err(err).Msg("Cannot reidentify data")
				log.Warn().Int("return", 1).Msg("End SIGO")
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&logs.verbosity, "verbosity", "v", "info",
		"set level of log verbosity : none (0), error (1), warn (2), info (3), debug (4), trace (5)")
	cmd.Flags().BoolVar(&logs.debug, "debug", false, "add debug information to logs (very slow)")
	cmd.Flags().BoolVar(&logs.jsonlog, "log-json", false, "output logs in JSON format")
	cmd.Flags().StringVar(&logs.colormode, "color", "auto", "use colors in log outputs : yes, no or auto")
	cmd.Flags().StringSliceVarP(&data.qi, "quasi-identifier", "q", []string{}, "list of quasi-identifying attributes")
	cmd.Flags().StringSliceVarP(&data.sensitive, "sensitive", "s", []string{}, "list of sensitive attributes")
	cmd.Flags().StringVar(&data.original, "load-original", "", "name and location of the original dataset file")
	cmd.Flags().StringVar(&data.anonymized, "load-anonymized", "", "name and location of the anonymized dataset file")
	//nolint: gomnd
	cmd.Flags().Float32VarP(&data.threshold, "threshold", "t", 0.5, "re-identification threshold")
	cmd.SetOut(out)
	cmd.SetErr(err)
	cmd.SetIn(in)

	return cmd
}

// nolint: cyclop
func initLog(logs logs) {
	color := false

	switch strings.ToLower(logs.colormode) {
	case "auto":
		if isatty.IsTerminal(os.Stdout.Fd()) && runtime.GOOS != "windows" {
			color = true
		}
	case "yes", "true", "1", "on", "enable":
		color = true
	}

	var logger zerolog.Logger
	if logs.jsonlog {
		logger = zerolog.New(os.Stderr)
	} else {
		// nolint: exhaustivestruct
		logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: !color})
	}

	if logs.debug {
		logger = logger.With().Caller().Logger()
	}

	over.New(logger)

	switch logs.verbosity {
	case "trace", "5":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
		log.Info().Msg("Logger level set to trace")
	case "debug", "4":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Info().Msg("Logger level set to debug")
	case "info", "3":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		log.Info().Msg("Logger level set to info")
	case "warn", "2":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error", "1":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}

	log.Info().Msgf("%v %v (commit=%v date=%v by=%v)", logs.name, logs.version, logs.commit, logs.buildDate, logs.builtBy)
}
