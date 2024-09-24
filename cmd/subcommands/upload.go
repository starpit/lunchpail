//go:build full || observe

package subcommands

import (
	"context"

	"github.com/spf13/cobra"

	"lunchpail.io/cmd/options"
	"lunchpail.io/pkg/be"
	"lunchpail.io/pkg/compilation"
	"lunchpail.io/pkg/runtime/queue"
	"lunchpail.io/pkg/runtime/queue/upload"
)

func newUploadCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "upload <srcDir> <bucket>",
		Short: "Copy data into queue",
		Long:  "Copy data into queue",
		Args:  cobra.MatchAll(cobra.ExactArgs(2), cobra.OnlyValidArgs),
	}

	var runname string
	cmd.Flags().StringVarP(&runname, "run", "r", "", "Inspect the given run, defaulting to using the singleton run")

	opts, err := options.RestoreCompilationOptions()
	if err != nil {
		return nil
	}

	options.AddTargetOptionsTo(cmd, &opts)

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		backend, err := be.New(ctx, opts)
		if err != nil {
			return err
		}

		return queue.UploadFiles(ctx, backend, runname, []upload.Upload{upload.Upload{Path: args[0], Bucket: args[1]}})
	}

	return cmd
}

func init() {
	if compilation.IsCompiled() {
		rootCmd.AddCommand(newUploadCmd())
	}
}