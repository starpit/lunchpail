package worker

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"lunchpail.io/cmd/options"
	"lunchpail.io/pkg/fe/transformer/api"
	"lunchpail.io/pkg/runtime/worker"
)

func PreStop() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "prestop",
		Short: "Mark this worker as dead",
		Long:  "Mark this worker as dead",
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

	logOpts := options.AddLogOptions(cmd)

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return worker.PreStop(context.Background(), worker.Options{
			LogOptions: *logOpts,
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
