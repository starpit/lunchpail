package queue

import (
	"context"
	"path/filepath"
	"strings"

	"lunchpail.io/pkg/be"
	"lunchpail.io/pkg/ir/queue"
)

func Ls(ctx context.Context, backend be.Backend, runname, path string) (<-chan string, <-chan error, error) {
	c, err := NewS3ClientForRun(ctx, backend, runname)
	if err != nil {
		return nil, nil, err
	}

	files := make(chan string)
	errors := make(chan error)

	run := queue.RunContext{Bucket: c.Paths.Bucket, RunName: runname, Step: 0}
	prefix := filepath.Join(run.ListenPrefix(), path)

	go func() {
		defer c.Stop()
		defer close(files)
		defer close(errors)
		for o := range c.ListObjects(c.Paths.Bucket, prefix, true) {
			if o.Err != nil {
				errors <- o.Err
			} else {
				files <- strings.Replace(o.Key, prefix+"/", "", 1)
			}
		}
	}()

	return files, errors, nil
}
