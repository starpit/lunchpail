package cpu

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	"lunchpail.io/pkg/be"
	"lunchpail.io/pkg/be/events/utilization"
	"lunchpail.io/pkg/be/runs/util"
	"lunchpail.io/pkg/ir/queue"
)

type CpuOptions struct {
	NoClearScreen   bool
	Verbose         bool
	IntervalSeconds int
}

func UI(ctx context.Context, runnameIn string, backend be.Backend, opts CpuOptions) error {
	runname, err := util.WaitForRun(ctx, runnameIn, true, backend)
	if err != nil {
		return err
	}

	group, sctx := errgroup.WithContext(ctx)

	c := make(chan utilization.Model)
	group.Go(func() error {
		defer close(c)
		return backend.Streamer(sctx, queue.RunContext{RunName: runname}).Utilization(c, opts.IntervalSeconds)
	})

	for model := range c {
		if !opts.Verbose && !opts.NoClearScreen {
			fmt.Print("\033[H\033[2J")
		}

		workers := model.Sorted()
		fmt.Println(cpuline(workers, cpu))
		fmt.Println(cpuline(workers, mem))
	}

	return nil
}
