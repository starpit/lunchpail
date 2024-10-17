package worker

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"lunchpail.io/cmd/options"
	"lunchpail.io/pkg/runtime/worker"
	"lunchpail.io/pkg/fe/transformer/api"
)

func Run() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "run [-- workerCommand workerCommandArg1 workerCommandArg2 ...]",
		Short: "Run as an application worker",
		Long:  "Run as an application worker",
		Args:  cobra.MatchAll(cobra.OnlyValidArgs),
	}

	var bucket string
	cmd.Flags().StringVar(&bucket, "bucket", "", "Which S3 bucket to use")
	cmd.MarkFlagRequired("bucket")

	var run string
	cmd.Flags().StringVar(&run, "run", "", "Which run are we part of")
	cmd.MarkFlagRequired("run")

	var step int
	cmd.Flags().IntVar(&step, "step", 0, "Which step are we part of")
	cmd.MarkFlagRequired("step")

	var pool string
	cmd.Flags().StringVar(&pool, "pool", "", "Which worker pool are we part of")
	cmd.MarkFlagRequired("pool")

	var pollingInterval int
	cmd.Flags().IntVar(&pollingInterval, "polling-interval", 3, "If polling is employed, the interval between probes")

	var startupDelay int
	cmd.Flags().IntVar(&startupDelay, "delay", 0, "Delay (in seconds) before engaging in any work")

	logOpts := options.AddLogOptions(cmd)

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Nothing to run. Specify the worker command line after a --: %v", args)
		}

		return worker.Run(context.Background(), args, worker.Options{
			StartupDelay:    startupDelay,
			PollingInterval: pollingInterval,
			LogOptions:      *logOpts,
			PathArgs: api.PathArgs{
				Bucket:       bucket,
				RunName: run,
				Step: step,
				PoolName: pool,
				WorkerName: os.Getenv("LUNCHPAIL_POD_NAME"),
			},
		})
	}

	return cmd
}
