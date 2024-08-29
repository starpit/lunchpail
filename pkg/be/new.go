package be

import (
	"fmt"

	"lunchpail.io/pkg/be/ibmcloud"
	"lunchpail.io/pkg/be/kubernetes"
	"lunchpail.io/pkg/compilation"
)

type TargetOptions struct {
	Namespace      string
	TargetPlatform Platform
}

func New(topts TargetOptions, aopts compilation.Options) (Backend, error) {
	var be Backend

	switch topts.TargetPlatform {
	case Kubernetes:
		be = kubernetes.Backend{Namespace: topts.Namespace}
	case IBMCloud:
		if ibm, err := ibmcloud.New(aopts); err != nil {
			return nil, err
		} else {
			be = ibm
		}
	default:
		return nil, fmt.Errorf("Unsupported backend %v", topts.TargetPlatform)
	}

	if err := be.Ok(); err != nil {
		return nil, err
	}

	return be, nil
}