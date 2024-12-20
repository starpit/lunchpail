package parser

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"lunchpail.io/pkg/defaults/application"
	"lunchpail.io/pkg/ir/hlir"
	"os"
	"strings"
)

func Parse(yamls string) (hlir.HLIR, error) {
	model := hlir.HLIR{}
	d := yaml.NewDecoder(strings.NewReader(yamls))

	for {
		var m hlir.UnknownResource
		if err := d.Decode(&m); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: skipping yaml with parse error %v", err)
			continue
		} else if len(m) == 0 {
			continue
		}

		kind, err := stringVal("kind", m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		}

		bytes, err := yaml.Marshal(m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: skipping invalid yaml %v", err)
			continue
		}

		switch kind {
		case "Application":
			var r hlir.Application
			if err := yaml.Unmarshal(bytes, &r); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: skipping yaml with invalid Application resource %v", err)
				continue
			} else {
				model.Applications = append(model.Applications, application.WithDefaults(r))
			}

		case "WorkerPool":
			var r hlir.WorkerPool
			if err := yaml.Unmarshal(bytes, &r); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: skipping yaml with invalid WorkerPool resource %v\n!!!!\n%s\n!!!!\n", err, string(bytes))
				continue
			} else {
				model.WorkerPools = append(model.WorkerPools, r)
			}

		default:
			model.Others = append(model.Others, m)
		}
	}

	return model, nil
}

func stringVal(key string, m hlir.UnknownResource) (string, error) {
	uval, ok := m[key]
	if !ok {
		return "", fmt.Errorf("Warning: skipping yaml with missing %s in %v", key, m)
	}

	val, ok := uval.(string)
	if !ok {
		return "", fmt.Errorf("Warning: skipping yaml with invalid %s in %v", key, uval)
	}

	return val, nil
}
