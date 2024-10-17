package workstealer

import (
	"context"
	"github.com/spf13/cobra"

	"lunchpail.io/cmd/options"
	"lunchpail.io/pkg/runtime/workstealer"
	"lunchpail.io/pkg/fe/transformer/api"
)

func Run() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "run",
		Short: "Run a work stealer",
		Long:  "Run a work stealer",
		Args:  cobra.MatchAll(cobra.ExactArgs(0), cobra.OnlyValidArgs),
	}

	var pollingInterval int
	cmd.Flags().IntVar(&pollingInterval, "polling-interval", 3, "If polling is employed, the interval between probes")

	var bucket string
	cmd.Flags().StringVar(&bucket, "bucket", "", "Which S3 bucket to use")
	cmd.MarkFlagRequired("bucket")

	var run string
	cmd.Flags().StringVar(&run, "run", "", "Which run are we part of")
	cmd.MarkFlagRequired("run")

	var step int
	cmd.Flags().IntVar(&step, "step", 0, "Which step are we part of")
	cmd.MarkFlagRequired("step")

	lopts := options.AddLogOptions(cmd)

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return workstealer.Run(context.Background(), api.PathArgs{Bucket: bucket, RunName: run, Step: step}, workstealer.Options{PollingInterval: pollingInterval, LogOptions: *lopts})
	}

	return cmd
}
