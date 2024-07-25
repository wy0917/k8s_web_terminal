package k8s

import (
	v1 "k8s.io/api/core/v1"
)

type WatchContainerLogRequest struct {
	Namespace string `json:"namespace" validate:"required"`
	PodName   string `json:"pod_name" validate:"required"`
	*v1.PodLogOptions
}

func NewWatchContainerLogRequest() *WatchContainerLogRequest {
	return &WatchContainerLogRequest{
		PodLogOptions: &v1.PodLogOptions{
			Follow:                       true,
			Previous:                     false,
			InsecureSkipTLSVerifyBackend: true,
		},
	}
}
