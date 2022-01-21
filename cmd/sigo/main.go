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

package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	over "github.com/Trendyol/overlog"
	"github.com/cgi-fr/sigo/internal/infra"
	"github.com/cgi-fr/sigo/pkg/sigo"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// Provisioned by ldflags
// nolint: gochecknoglobals
var (
	name      string
	version   string
	commit    string
	buildDate string
	builtBy   string

	verbosity string
	debug     bool
	jsonlog   bool
	colormode string
	k         int
	l         int
	qi        []string
	sensitive []string
	method    string
	info      string
)

func main() {
	//nolint: exhaustivestruct
	rootCmd := &cobra.Command{
		Use:   name,
		Short: "Command line to generalize and anonymize the content of a jsonline flow set",
		Version: fmt.Sprintf(`%v (commit=%v date=%v by=%v)
	Copyright (C) 2022 CGI France \n License GPLv3: GNU GPL version 3 <https://gnu.org/licenses/gpl.html>.
	This is free software: you are free to change and redistribute it.
	There is NO WARRANTY, to the extent permitted by law.`, version, commit, buildDate, builtBy),
		Run: func(cmd *cobra.Command, args []string) {
			// nolint: exhaustivestruct
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

			run()
		},
	}

	var entropy bool

	rootCmd.PersistentFlags().
		StringVarP(&verbosity, "verbosity", "v", "info",
			"set level of log verbosity : none (0), error (1), warn (2), info (3), debug (4), trace (5)")
	rootCmd.PersistentFlags().
		BoolVar(&debug, "debug", false, "add debug information to logs (very slow)")
	rootCmd.PersistentFlags().
		BoolVar(&jsonlog, "log-json", false, "output logs in JSON format")
	rootCmd.PersistentFlags().StringVar(&colormode, "color", "auto", "use colors in log outputs : yes, no or auto")
	// nolint: gomnd
	rootCmd.PersistentFlags().
		IntVarP(&k, "k-value", "k", 3, "k-value for k-anonymization")
	rootCmd.PersistentFlags().
		IntVarP(&l, "l-value", "l", 1, "l-value for l-diversity")
	rootCmd.PersistentFlags().
		StringSliceVarP(&qi, "quasi-identifier", "q", []string{}, "list of quasi-identifying attributes")
	rootCmd.PersistentFlags().
		StringSliceVarP(&sensitive, "sensitive", "s", []string{}, "list of sensitive attributes")
	rootCmd.PersistentFlags().
		StringVarP(&method, "anonymizer", "a", "general",
			"anonymization method used. Select one from this list "+
				"['general', 'meanAggregation', 'medianAggragation', 'outlier', 'laplaceNoise', 'gaussianNoise']")
	rootCmd.PersistentFlags().
		StringVarP(&info, "cluster-info", "i", "", "display cluster for each jsonline flow")
	rootCmd.PersistentFlags().BoolVar(&entropy, "entropy", false, "use entropy model for l-diversity")
	over.MDC().Set("entropy", entropy)

	if err := rootCmd.Execute(); err != nil {
		log.Err(err).Msg("Error when executing command")
		os.Exit(1)
	}
}

func run() {
	initLog()

	log.Info().
		Int("k-anonymity", k).
		Int("l-diversity", l).
		Strs("Quasi-Identifiers", qi).
		Strs("Sensitive", sensitive).
		Str("Method", method).
		Str("Cluster-Info", info).
		Msg("Start SIGO")

	source, err := infra.NewJSONLineSource(os.Stdin, qi, sensitive)
	if err != nil {
		log.Err(err).Msg("Cannot load jsonline source")
		log.Warn().Int("return", 1).Msg("End SIGO")
		os.Exit(1)
	}

	sink := infra.NewJSONLineSink(os.Stdout)

	var debugger sigo.Debugger

	if info != "" {
		debugger = sigo.NewSequenceDebugger(info)
	} else {
		debugger = sigo.NewNoDebugger()
	}

	var anonymizer sigo.Anonymizer

	anonymizers := sigo.NewAnonymizers()
	if method != "" {
		anonymizer = anonymizers.Anonymizer("NoAnonymizer")
	} else {
		anonymizer = anonymizers.Anonymizer(method)
	}

	err = sigo.Anonymize(source, sigo.NewKDTreeFactory(), k, l, len(qi), anonymizer, sink, debugger)
	if err != nil {
		panic(err)
	}
}

// nolint: cyclop
func initLog() {
	color := false

	switch strings.ToLower(colormode) {
	case "auto":
		if isatty.IsTerminal(os.Stdout.Fd()) && runtime.GOOS != "windows" {
			color = true
		}
	case "yes", "true", "1", "on", "enable":
		color = true
	}

	var logger zerolog.Logger
	if jsonlog {
		logger = zerolog.New(os.Stderr)
	} else {
		// nolint: exhaustivestruct
		logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: !color})
	}

	if debug {
		logger = logger.With().Caller().Logger()
	}

	over.New(logger)

	switch verbosity {
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

	log.Info().Msgf("%v %v (commit=%v date=%v by=%v)", name, version, commit, buildDate, builtBy)
}
