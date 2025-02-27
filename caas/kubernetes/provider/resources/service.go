// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package resources

import (
	"context"
	"time"

	"github.com/juju/errors"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"

	k8sconstants "github.com/DavinZhang/juju/caas/kubernetes/provider/constants"
	"github.com/DavinZhang/juju/core/status"
)

// Service extends the k8s service.
type Service struct {
	corev1.Service
}

// NewService creates a new service resource.
func NewService(name string, namespace string, in *corev1.Service) *Service {
	if in == nil {
		in = &corev1.Service{}
	}
	in.SetName(name)
	in.SetNamespace(namespace)
	return &Service{*in}
}

// Clone returns a copy of the resource.
func (s *Service) Clone() Resource {
	clone := *s
	return &clone
}

// ID returns a comparable ID for the Resource
func (s *Service) ID() ID {
	return ID{"Service", s.Name, s.Namespace}
}

// Apply patches the resource change.
func (s *Service) Apply(ctx context.Context, client kubernetes.Interface) error {
	api := client.CoreV1().Services(s.Namespace)
	data, err := runtime.Encode(unstructured.UnstructuredJSONScheme, &s.Service)
	if err != nil {
		return errors.Trace(err)
	}
	res, err := api.Patch(ctx, s.Name, types.StrategicMergePatchType, data, metav1.PatchOptions{
		FieldManager: JujuFieldManager,
	})
	if k8serrors.IsNotFound(err) {
		res, err = api.Create(ctx, &s.Service, metav1.CreateOptions{
			FieldManager: JujuFieldManager,
		})
	}
	if k8serrors.IsConflict(err) {
		return errors.Annotatef(errConflict, "service %q", s.Name)
	}
	if err != nil {
		return errors.Trace(err)
	}
	s.Service = *res
	return nil
}

// Get refreshes the resource.
func (s *Service) Get(ctx context.Context, client kubernetes.Interface) error {
	api := client.CoreV1().Services(s.Namespace)
	res, err := api.Get(ctx, s.Name, metav1.GetOptions{})
	if k8serrors.IsNotFound(err) {
		return errors.NewNotFound(err, "k8s")
	} else if err != nil {
		return errors.Trace(err)
	}
	s.Service = *res
	return nil
}

// Delete removes the resource.
func (s *Service) Delete(ctx context.Context, client kubernetes.Interface) error {
	api := client.CoreV1().Services(s.Namespace)
	err := api.Delete(ctx, s.Name, metav1.DeleteOptions{
		PropagationPolicy: k8sconstants.DefaultPropagationPolicy(),
	})
	if k8serrors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return errors.Trace(err)
	}
	return nil
}

// Events emitted by the resource.
func (s *Service) Events(ctx context.Context, client kubernetes.Interface) ([]corev1.Event, error) {
	return ListEventsForObject(ctx, client, s.Namespace, s.Name, "Service")
}

// ComputeStatus returns a juju status for the resource.
func (s *Service) ComputeStatus(ctx context.Context, client kubernetes.Interface, now time.Time) (string, status.Status, time.Time, error) {
	if s.DeletionTimestamp != nil {
		return "", status.Terminated, s.DeletionTimestamp.Time, nil
	}
	return "", status.Active, s.CreationTimestamp.Time, nil
}
