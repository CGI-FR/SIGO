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
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

	over "github.com/Trendyol/overlog"
	"github.com/cgi-fr/sigo/internal/infra"
	"github.com/cgi-fr/sigo/pkg/reidentification"
	"github.com/cgi-fr/sigo/pkg/sigo"
	"github.com/mattn/go-isatty"
	"github.com/pkg/profile"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type infos struct {
	name      string
	version   string
	commit    string
	buildDate string
	builtBy   string
}

type logs struct {
	verbosity string
	debug     bool
	jsonlog   bool
	colormode string
	info      string
	profiling bool
}

type pdef struct {
	k         int
	l         int
	qi        []string
	sensitive []string
	method    string
	cmdLine   []string
	config    string
}

type reid struct {
	originalFile string
}

func main() {
	var info infos

	var logs logs

	var definition pdef

	//nolint: exhaustivestruct
	rootCmd := &cobra.Command{
		Use:   info.name,
		Short: "Command line to generalize and anonymize the content of a jsonline flow set",
		Version: fmt.Sprintf(`%v (commit=%v date=%v by=%v)
	Copyright (C) 2022 CGI France \n License GPLv3: GNU GPL version 3 <https://gnu.org/licenses/gpl.html>.
	This is free software: you are free to change and redistribute it.
	There is NO WARRANTY, to the extent permitted by law.`, info.version, info.commit, info.buildDate, info.builtBy),
		Run: func(cmd *cobra.Command, args []string) {
			// nolint: exhaustivestruct
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

			definition.flagIsSet(*cmd)

			run(info, definition, logs)
		},
	}

	var entropy bool

	rootCmd.PersistentFlags().
		StringVarP(&logs.verbosity, "verbosity", "v", "info",
			"set level of log verbosity : none (0), error (1), warn (2), info (3), debug (4), trace (5)")
	rootCmd.PersistentFlags().
		BoolVar(&logs.debug, "debug", false, "add debug information to logs (very slow)")
	rootCmd.PersistentFlags().
		BoolVar(&logs.jsonlog, "log-json", false, "output logs in JSON format")
	rootCmd.PersistentFlags().StringVar(&logs.colormode, "color", "auto", "use colors in log outputs : yes, no or auto")
	// nolint: gomnd
	rootCmd.PersistentFlags().
		IntVarP(&definition.k, "k-value", "k", 3, "k-value for k-anonymization")
	rootCmd.PersistentFlags().
		IntVarP(&definition.l, "l-value", "l", 1, "l-value for l-diversity")
	rootCmd.PersistentFlags().
		StringSliceVarP(&definition.qi, "quasi-identifier", "q", []string{}, "list of quasi-identifying attributes")
	rootCmd.PersistentFlags().
		StringSliceVarP(&definition.sensitive, "sensitive", "s", []string{}, "list of sensitive attributes")
	rootCmd.PersistentFlags().
		StringVarP(&definition.method, "anonymizer", "a", "",
			"anonymization method used. Select one from this list "+
				"['general', 'meanAggregation', 'medianAggregation', 'outlier', 'laplaceNoise', 'gaussianNoise', 'swapping']")
	rootCmd.PersistentFlags().
		StringVarP(&logs.info, "cluster-info", "i", "", "display cluster for each jsonline flow")
	rootCmd.PersistentFlags().BoolVarP(&logs.profiling, "profiling", "p", false,
		"start sigo with profiling and generate a cpu.pprof file (debug)")
	rootCmd.PersistentFlags().BoolVar(&entropy, "entropy", false, "use entropy model for l-diversity")
	over.MDC().Set("entropy", entropy)
	rootCmd.PersistentFlags().
		StringVarP(&definition.config, "configuration", "c", "sigo.yml", "name and location of the configuration file")
	rootCmd.PersistentFlags().StringVar(&reid.originalFile, "load-original", "",
		"name and location of the original dataset file")

	if err := rootCmd.Execute(); err != nil {
		log.Err(err).Msg("Error when executing command")
		os.Exit(1)
	}
}

func run(info infos, definition pdef, logs logs) {
	initLog(logs, info)

	// if the configuration file is present in the current directory
	if sigo.Exist(definition.config) {
		if err := definition.initConfig(); err != nil {
			log.Err(err).Msg("Cannot load configuration definition from file")
			log.Warn().Int("return", 1).Msg("End SIGO")
			os.Exit(1)
		}
	}

	log.Info().
		Str("configuration", definition.config).
		Int("k-anonymity", definition.k).
		Int("l-diversity", definition.l).
		Strs("Quasi-Identifiers", definition.qi).
		Strs("Sensitive", definition.sensitive).
		Str("Method", definition.method).
		Str("Cluster-Info", logs.info).
		Bool("Re-identification", reid.originalFile != "").
		Msg("Start SIGO")

	source, err := infra.NewJSONLineSource(os.Stdin, definition.qi, definition.sensitive)
	if err != nil {
		log.Err(err).Msg("Cannot load jsonline source")
		log.Warn().Int("return", 1).Msg("End SIGO")
		os.Exit(1)
	}

	sink := infra.NewJSONLineSink(os.Stdout)

	var debugger sigo.Debugger

	if logs.info != "" {
		debugger = sigo.NewSequenceDebugger(logs.info)
	} else {
		debugger = sigo.NewNoDebugger()
	}

	switch {
	case reid.originalFile != "":
		originalData, err := os.Open(reid.originalFile)
		if err != nil {
			log.Err(err).Msg("Cannot open original dataset")
			log.Warn().Int("return", 1).Msg("End SIGO")
			os.Exit(1)
		}

		original, err := infra.NewJSONLineSource(bufio.NewReader(originalData), definition.qi, definition.sensitive)
		if err != nil {
			log.Err(err).Msg("Cannot load jsonline original dataset")
			log.Warn().Int("return", 1).Msg("End SIGO")
			os.Exit(1)
		}

		err = reidentification.ReIdentify(original, source, reidentification.NewIdentifier("canberra", k), sink)
		if err != nil {
			panic(err)
		}
	default:
		var cpuProfiler interface{ Stop() }

		if logs.profiling {
			cpuProfiler = profile.Start(profile.ProfilePath("."))
		}

		err = sigo.Anonymize(source, sigo.NewKDTreeFactory(), definition.k, definition.l,
			len(definition.qi), newAnonymizer(definition.method), sink, debugger)
		if err != nil {
			panic(err)
		}

		if logs.profiling {
			cpuProfiler.Stop()
		}
	}
}

// nolint: cyclop
func initLog(logs logs, info infos) {
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

	log.Info().Msgf("%v %v (commit=%v date=%v by=%v)", info.name, info.version, info.commit, info.buildDate, info.builtBy)
}

func newAnonymizer(name string) sigo.Anonymizer {
	switch name {
	case "general":
		return sigo.NewGeneralAnonymizer()
	case "meanAggregation":
		return sigo.NewAggregationAnonymizer("mean")
	case "medianAggregation":
		return sigo.NewAggregationAnonymizer("median")
	case "outlier":
		return sigo.NewCodingAnonymizer()
	case "laplaceNoise":
		return sigo.NewNoiseAnonymizer("laplace")
	case "gaussianNoise":
		return sigo.NewNoiseAnonymizer("gaussian")
	case "swapping":
		return sigo.NewSwapAnonymizer()
	default:
		return sigo.NewNoAnonymizer()
	}
}

// initConfig initialize sigo configuration with config file.
func (def *pdef) initConfig() (err error) {
	pdf, err := sigo.LoadConfigurationFromYAML(def.config)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// if cmdLine not contains "k"
	// then we take the value in the configuration file
	// else we take the value put in command line
	if !sigo.Contains(def.cmdLine, "k") {
		def.k = pdf.K
	}

	if !sigo.Contains(def.cmdLine, "l") {
		def.l = pdf.L
	}

	if !sigo.Contains(def.cmdLine, "sensitive") {
		def.sensitive = pdf.Sensitive
	}

	if !sigo.Contains(def.cmdLine, "method") {
		def.method = pdf.Aggregation
	}

	if !sigo.Contains(def.cmdLine, "qi") {
		for _, attributes := range pdf.Rules {
			def.qi = append(def.qi, attributes.Name)
		}
	}

	return nil
}

// flagIsSet adds to cmdLine the flags set on the command line.
func (def *pdef) flagIsSet(cmd cobra.Command) {
	// if k is given as parameter to sigo
	// then k is appended to cmdLine
	if cmd.Root().Flag("k-value").Changed {
		def.cmdLine = append(def.cmdLine, "k")
	}

	if cmd.Root().Flag("l-value").Changed {
		def.cmdLine = append(def.cmdLine, "l")
	}

	if cmd.Root().Flag("quasi-identifier").Changed {
		def.cmdLine = append(def.cmdLine, "qi")
	}

	if cmd.Root().Flag("sensitive").Changed {
		def.cmdLine = append(def.cmdLine, "sensitive")
	}

	if cmd.Root().Flag("anonymizer").Changed {
		def.cmdLine = append(def.cmdLine, "method")
	}
}
