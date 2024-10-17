package api

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"lunchpail.io/pkg/fe/linker/queue"
)

// i.e. "/run/{{.RunName}}/step/{{.Step}}"
func (q PathArgs) ListenPrefix() string {
	A := strings.Split(Unassigned, "/")
	return filepath.Join(A[0:5]...)
}

const (
	Unassigned = "/run/{{.RunName}}/step/{{.Step}}/unassigned"

	AssignedAndPending = "/run/{{.RunName}}/step/{{.Step}}/inbox/pool/{{.PoolName}}/worker/{{.WorkerName}}/{{.Task}}"
	AssignedAndProcessing = "/run/{{.RunName}}/step/{{.Step}}/processing/pool/{{.PoolName}}/worker/{{.WorkerName}}/{{.Task}}"
	AssignedAndFinished = "/run/{{.RunName}}/step/{{.Step+1}}/unassigned" // i.e. step 1's output is step 2's input

	FinishedWithCode = "/run/{{.RunName}}/step/{{.Step}}/codes/pool/{{.PoolName}}/worker/{{.WorkerName}}/{{.Task}}"
	FinishedWithStdout = "/run/{{.RunName}}/step/{{.Step}}/stdout/pool/{{.PoolName}}/worker/{{.WorkerName}}/{{.Task}}"
	FinishedWithStderr = "/run/{{.RunName}}/step/{{.Step}}/stderr/pool/{{.PoolName}}/worker/{{.WorkerName}}/{{.Task}}"
	FinishedWithSucceeded = "/run/{{.RunName}}/step/{{.Step}}/succeeded/pool/{{.PoolName}}/worker/{{.WorkerName}}/{{.Task}}"
	FinishedWithFailed = "/run/{{.RunName}}/step/{{.Step}}/failed/pool/{{.PoolName}}/worker/{{.WorkerName}}/{{.Task}}"

	WorkerKillFile = "/run/{{.RunName}}/step/{{.Step}}/killfiles/pool/{{.PoolName}}/worker/{{.WorkerName}}"
	
	AllDoneMarker = "/run/{{.RunName}}/step/{{.Step}}/markers/alldone"
	WorkerAliveMarker = "/run/{{.RunName}}/step/{{.Step}}/markers/alive/pool/{{.PoolName}}/worker/{{.WorkerName}}"
	WorkerDeadMarker = "/run/{{.RunName}}/step/{{.Step}}/markers/dead/pool/{{.PoolName}}/worker/{{.WorkerName}}"
)

type PathArgs struct {
	Bucket string
	RunName string
	Step int
	PoolName string
	WorkerName string
	Task string
}

func (q PathArgs) ForPool(name string) PathArgs {
	q.PoolName = name
	return q
}

func (q PathArgs) ForWorker(name string) PathArgs {
	q.WorkerName = name
	return q
}

func (q PathArgs) ForTask(name string) PathArgs {
	q.Task = name
	return q
}

func (q PathArgs) Template(path string) (string, error) {
	tmpl, err := template.New("tmp").Parse(path)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	if err := tmpl.Execute(&b, q); err != nil {
		return "", err
	}

	return b.String(), nil
}

func (q PathArgs) TemplateP(path string) string {
	s, err := q.Template(path)
	if err != nil {
		return ""
	}
	return s
}

// Path in s3 to store the queue for the given run
func QueuePrefixPath0(queueSpec queue.Spec, runname string) string {
	return filepath.Join("lunchpail", "runs", runname)
}

// Path in s3 to store the queue for the given run
func QueuePrefixPath(queueSpec queue.Spec, runname string) string {
	return filepath.Join(queueSpec.Bucket, QueuePrefixPath0(queueSpec, runname))
}

// The queue path for a worker that specifies the pool name and the worker name
func QueueSubPathForWorker(poolName, workerName string) string {
	return filepath.Join(poolName, workerName)
}

// Opposite of `QueueSubPathForWorker`, e.g. test7f-pool1/w96bh -> (test7f-pool1,w96bh)
func ExtractNamesFromSubPathForWorker(combo string) (poolName string, workerName string, err error) {
	if idx := strings.Index(combo, "/"); idx < 0 {
		// TODO error handling here. what do we want to do?
		err = fmt.Errorf("Invalid subpath %s", combo)
	} else {
		poolName = combo[:idx]
		workerName = combo[idx+1:]
	}
	return
}

// Path in s3 to store the queue data for a particular worker in a
// particular pool for a particular run. Note how we need to defer the
// worker name until run time, which we do via a
// $LUNCHPAIL_POD_NAME env var.
func QueuePrefixPathForWorker0(queueSpec queue.Spec, runname, poolName string) string {
	return filepath.Join(
		QueuePrefixPath0(queueSpec, runname),
		"queues",
		QueueSubPathForWorker(poolName, "$LUNCHPAIL_POD_NAME"),
	)
}

// Path in s3 to store the queue data for a particular worker in a
// particular pool for a particular run. Note how we need to defer the
// worker name until run time, which we do via a
// $LUNCHPAIL_POD_NAME env var.
func QueuePrefixPathForWorker(queueSpec queue.Spec, runname, poolName string) string {
	return filepath.Join(
		queueSpec.Bucket,
		QueuePrefixPathForWorker0(queueSpec, runname, poolName),
	)
}
