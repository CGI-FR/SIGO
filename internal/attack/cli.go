package attack

import (
	"fmt"
	"os"

	"github.com/cgi-fr/sigo/internal/infra"
	"github.com/cgi-fr/sigo/pkg/attack"
	"github.com/cgi-fr/sigo/pkg/sigo"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// nolint: gochecknoglobals
var (
	qi        []string
	sensitive []string
	info      string
)

func Inject(q, s []string, i string) {
	qi, sensitive, info = q, s, i
}

func NewCommand(fullName string, err *os.File, out *os.File, in *os.File) *cobra.Command {
	// nolint: exhaustivestruct
	cmd := &cobra.Command{
		Use:     "attack ",
		Short:   "Push data to a database with a pushing mode (insert by default)",
		Long:    "",
		Example: fmt.Sprintf("  %[1]s attack \n  %[1]s attack", fullName),
		Run: func(cmd *cobra.Command, args []string) {
			source, err := infra.NewJSONLineSource(cmd.InOrStdin(), qi, sensitive)
			if err != nil {
				log.Err(err).Msg("Cannot load jsonline source")
				log.Warn().Int("return", 1).Msg("End SIGO")
				os.Exit(1)
			}
			sink := infra.NewJSONLineSink(cmd.OutOrStdout())

			var debugger sigo.Debugger

			if info != "" {
				debugger = sigo.NewSequenceDebugger(info)
			} else {
				debugger = sigo.NewNoDebugger()
			}

			err = attack.Identify(source, sigo.NewKDTreeFactory(), len(qi), sink, debugger)
			if err != nil {
				panic(err)
			}
		},
	}
	cmd.SetOut(out)
	cmd.SetErr(err)
	cmd.SetIn(in)

	return cmd
}
