package worker

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"lunchpail.io/pkg/fe/transformer/api"
	"lunchpail.io/pkg/runtime/queue"
)

func killfileExists(client queue.S3Client, opts Options) bool {
	path := opts.PathArgs.TemplateP(api.WorkerKillFile)
	return client.Exists(opts.PathArgs.Bucket, filepath.Dir(path), filepath.Base(path))
}

func startWatch(ctx context.Context, handler []string, client queue.S3Client, opts Options) error {
	if err := client.Mkdirp(opts.PathArgs.Bucket); err != nil {
		return err
	}

	alive := opts.PathArgs.TemplateP(api.WorkerAliveMarker)
	if opts.LogOptions.Debug {
		fmt.Fprintf(os.Stderr, "Lunchpail worker touching alive file bucket=%s path=%s\n", opts.PathArgs.Bucket, alive)
	}
	err := client.Touch(opts.PathArgs.Bucket, alive)
	if err != nil {
		return err
	}

	localdir, err := os.MkdirTemp("", "lunchpail_local_queue_")
	if err != nil {
		return err
	}

	tasks, errs := client.Listen(opts.PathArgs.Bucket, opts.PathArgs.ListenPrefix(), "", false)
	for {
		if killfileExists(client, opts) {
			if opts.LogOptions.Verbose {
				fmt.Fprintln(os.Stderr, "Worker got kill file. Shutting down now.")
			}
			break
		}

		select {
		case err := <-errs:
			if opts.LogOptions.Verbose {
				fmt.Fprintln(os.Stderr, err)
			}

			// sleep for a bit
			s := opts.PollingInterval
			if s == 0 {
				s = 3
			}
			time.Sleep(time.Duration(s) * time.Second)

		case <-tasks:
		}

		tasks, err := client.Lsf(opts.PathArgs.Bucket, opts.PathArgs.TemplateP(api.AssignedAndPending))
		if err != nil {
			return err
		}

		for _, task := range tasks {
			if task != "" {
				// TODO: re-check if task still exists in our inbox before starting on it
				a := opts.PathArgs.ForTask(task)
				in := a.TemplateP(api.AssignedAndPending)
				inprogress := a.TemplateP(api.AssignedAndProcessing)
				out := a.TemplateP(api.AssignedAndFinished)

				// capture exit code, stdout and stderr of the handler
				ec := a.TemplateP(api.FinishedWithCode)
				failed := a.TemplateP(api.FinishedWithFailed)
				succeeded := a.TemplateP(api.FinishedWithSucceeded)
				stdout := a.TemplateP(api.FinishedWithStdout)
				stderr := a.TemplateP(api.FinishedWithStderr)

				localinbox := filepath.Join(localdir, "inbox")
				localprocessing := filepath.Join(localdir, "processing", task)
				localoutbox := filepath.Join(localdir, "outbox", task)
				localec := localoutbox + ".code"
				localstdout := localoutbox + ".stdout"
				localstderr := localoutbox + ".stderr"

				err := os.MkdirAll(localinbox, os.ModePerm)
				if err != nil {
					fmt.Println("Internal Error creating local inbox:", err)
					continue
				}
				err = os.MkdirAll(filepath.Dir(localprocessing), os.ModePerm)
				if err != nil {
					fmt.Println("Internal Error creating local processing:", err)
					continue
				}
				err = os.MkdirAll(filepath.Dir(localoutbox), os.ModePerm)
				if err != nil {
					fmt.Println("Internal Error creating local outbox:", err)
					continue
				}

				err = client.Download(opts.PathArgs.Bucket, in, localprocessing)
				if err != nil {
					if !strings.Contains(err.Error(), "key does not exist") {
						// we ignore "key does not exist" errors, as these result from the work
						// we thought we were assigned having been stolen by the workstealer
						fmt.Printf("Internal Error copying task to worker processing %s %s->%s: %v\n", opts.PathArgs.Bucket, in, localprocessing, err)
					}
					continue
				}

				// fmt.Println("sending file to handler: " + in)
				err = os.Remove(localoutbox)
				if err != nil && !os.IsNotExist(err) {
					fmt.Println("Internal Error removing task from local outbox:", err)
					continue
				}

				err = client.Moveto(opts.PathArgs.Bucket, in, inprogress)
				if err != nil {
					fmt.Printf("Internal Error moving task to global processing %s->%s: %v\n", in, inprogress, err)
					continue
				}

				// signify that the process is still going... or prematurely terminated
				os.WriteFile(localec, []byte("-1"), os.ModePerm)

				handlercmd := exec.CommandContext(ctx, handler[0], slices.Concat(handler[1:], []string{localprocessing, localoutbox})...)

				// open stdout/err files for writing
				stdoutfile, err := os.Create(localstdout)
				if err != nil {
					fmt.Println("Internal Error creating stdout file:", err)
					continue
				}
				defer stdoutfile.Close()

				stderrfile, err := os.Create(localstderr)
				if err != nil {
					fmt.Println("Internal Error creating stderr file:", err)
					continue
				}
				defer stderrfile.Close()

				multiout := io.MultiWriter(os.Stdout, stdoutfile)
				multierr := io.MultiWriter(os.Stderr, stderrfile)
				handlercmd.Stdout = multiout
				handlercmd.Stderr = multierr
				err = handlercmd.Run()
				if err != nil {
					fmt.Println("Internal Error running the handler:", err)
					multierr.Write([]byte(err.Error()))
				}
				EC := handlercmd.ProcessState.ExitCode()

				os.WriteFile(localec, []byte(fmt.Sprintf("%d", EC)), os.ModePerm)

				err = client.Upload(opts.PathArgs.Bucket, localec, ec)
				if err != nil {
					fmt.Println("Internal Error moving exitcode to remote:", err)
				}

				err = client.Upload(opts.PathArgs.Bucket, localstdout, stdout)
				if err != nil {
					fmt.Println("Internal Error moving stdout to remote:", err)
				}

				err = client.Upload(opts.PathArgs.Bucket, localstderr, stderr)
				if err != nil {
					fmt.Println("Internal Error moving stderr to remote:", err)
				}

				if EC == 0 {
					if opts.LogOptions.Debug {
						fmt.Printf("Worker succeeded on task %s\n", localprocessing)
					}
					err = client.Touch(opts.PathArgs.Bucket, succeeded)
					if err != nil {
						fmt.Println("Internal Error creating succeeded marker:", err)
					}
				} else {
					err = client.Touch(opts.PathArgs.Bucket, failed)
					if err != nil {
						fmt.Println("Internal Error creating failed marker:", err)
					}
					fmt.Println("Worker error with exit code " + strconv.Itoa(EC) + " while processing " + filepath.Base(in))
				}

				if _, err := os.Stat(localoutbox); err == nil {
					if opts.LogOptions.Debug {
						fmt.Printf("Uploading worker-produced outbox file %s->%s\n", localoutbox, out)
					}
					if err := client.Upload(opts.PathArgs.Bucket, localoutbox, out); err != nil {
						fmt.Printf("Internal Error uploading task to global outbox %s->%s: %v\n", localoutbox, out, err)
					}
				} else if err := client.Moveto(opts.PathArgs.Bucket, inprogress, out); err != nil {
					fmt.Printf("Internal Error moving task to global outbox %s->%s: %v\n", inprogress, out, err)
				}
			}
		}
	}

	if opts.LogOptions.Verbose {
		fmt.Println("Worker exiting normally")
	}

	return nil
}
