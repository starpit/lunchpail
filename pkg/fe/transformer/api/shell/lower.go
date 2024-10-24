package shell

import (
	"encoding/base64"
	"fmt"
	"path/filepath"
	"strconv"

	"lunchpail.io/pkg/build"
	q "lunchpail.io/pkg/fe/linker/queue"
	"lunchpail.io/pkg/fe/transformer/api"
	"lunchpail.io/pkg/ir/hlir"
	"lunchpail.io/pkg/ir/llir"
	"lunchpail.io/pkg/lunchpail"
)

func Lower(buildName string, ctx llir.Context, app hlir.Application, opts build.Options) (llir.Component, error) {
	var component lunchpail.Component
	switch app.Spec.Role {
	case "worker":
		component = lunchpail.WorkersComponent
	default:
		component = lunchpail.DispatcherComponent
	}

	return LowerAsComponent(buildName, ctx, app, llir.ShellComponent{Component: component}, opts)
}

func LowerAsComponent(buildName string, ctx llir.Context, app hlir.Application, component llir.ShellComponent, opts build.Options) (llir.Component, error) {
	component.Application = app
	if component.Sizing.Workers == 0 {
		component.Sizing = api.ApplicationSizing(app, opts)
	}
	if component.InstanceName == "" {
		component.InstanceName = ctx.Run.RunName
	}

	if app.Spec.Env == nil {
		app.Spec.Env = hlir.Env{}
	}
	app.Spec.Env["LUNCHPAIL_RUN_NAME"] = ctx.Run.RunName
	app.Spec.Env["LUNCHPAIL_STEP"] = strconv.Itoa(ctx.Run.Step)
	app.Spec.Env["LUNCHPAIL_QUEUE_BUCKET"] = ctx.Queue.Bucket

	for _, needs := range app.Spec.Needs {
		var req string

		if needs.Requirements != "" {
			req = "--requirements " + base64.StdEncoding.EncodeToString([]byte(needs.Requirements))
			if opts.Log.Verbose {
				fmt.Printf("Setting requirements for needs \n")
			}
		}

		component.Spec.Command = fmt.Sprintf(`PATH=$($LUNCHPAIL_EXE needs %s %s %s --verbose=%v):$PATH %s`, needs.Name, needs.Version, req, opts.Log.Verbose, component.Spec.Command)
	}

	for _, dataset := range app.Spec.Datasets {
		if dataset.S3.Rclone.RemoteName != "" && dataset.S3.CopyIn.Path != "" {
			// We were asked to copy data in from s3, so
			// we will use the secrets attached to an
			// initContainer
			isValid, spec, err := q.SpecFromRcloneRemoteName(dataset.S3.Rclone.RemoteName, "", ctx.Run.RunName, ctx.Queue.Port)

			if err != nil {
				return nil, err
			} else if !isValid {
				return nil, fmt.Errorf("Error: invalid or missing rclone config for given remote=%s for Application=%s", dataset.S3.Rclone.RemoteName, app.Metadata.Name)
			}

			// sleep to delay the copy-out, if requested
			component.Spec.Command = fmt.Sprintf(`sleep %d
env lunchpail_queue_endpoint=%s lunchpail_queue_accessKeyID=%s lunchpail_queue_secretAccessKey=%s $LUNCHPAIL_EXE queue download %s %s/%s
%s`, dataset.S3.CopyIn.Delay, spec.Endpoint, spec.AccessKey, spec.SecretKey, dataset.S3.CopyIn.Path, dataset.Name, filepath.Base(dataset.S3.CopyIn.Path), component.Spec.Command)
		}
	}
	return component, nil
}
