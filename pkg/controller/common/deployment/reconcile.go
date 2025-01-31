// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package deployment

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/elastic/cloud-on-k8s/pkg/controller/common"
	"github.com/elastic/cloud-on-k8s/pkg/controller/common/hash"
	"github.com/elastic/cloud-on-k8s/pkg/controller/common/reconciler"
	"github.com/elastic/cloud-on-k8s/pkg/utils/k8s"
)

var (
	defaultRevisionHistoryLimit int32
)

// Params to specify a Deployment specification.
type Params struct {
	Name            string
	Namespace       string
	Selector        map[string]string
	Labels          map[string]string
	PodTemplateSpec corev1.PodTemplateSpec
	Replicas        int32
}

// New creates a Deployment from the given params.
func New(params Params) appsv1.Deployment {
	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      params.Name,
			Namespace: params.Namespace,
			Labels:    params.Labels,
		},
		Spec: appsv1.DeploymentSpec{
			RevisionHistoryLimit: common.Int32(defaultRevisionHistoryLimit),
			Selector: &metav1.LabelSelector{
				MatchLabels: params.Selector,
			},
			Template: params.PodTemplateSpec,
			Replicas: &params.Replicas,
		},
	}
}

// ReconcileDeployment creates or updates the given deployment for the specified owner.
func Reconcile(
	k8sClient k8s.Client,
	scheme *runtime.Scheme,
	expected appsv1.Deployment,
	owner metav1.Object,
) (appsv1.Deployment, error) {
	// label the deployment with a hash of itself
	expected = WithTemplateHash(expected)

	reconciled := &appsv1.Deployment{}
	err := reconciler.ReconcileResource(reconciler.Params{
		Client:     k8sClient,
		Scheme:     scheme,
		Owner:      owner,
		Expected:   &expected,
		Reconciled: reconciled,
		NeedsUpdate: func() bool {
			// compare hash of the deployment at the time it was built
			return hash.GetTemplateHashLabel(reconciled.Labels) != hash.GetTemplateHashLabel(expected.Labels)
		},
		UpdateReconciled: func() {
			expected.DeepCopyInto(reconciled)
		},
	})
	return *reconciled, err
}

// WithTemplateHash returns a new deployment with a hash of its template to ease comparisons.
func WithTemplateHash(d appsv1.Deployment) appsv1.Deployment {
	dCopy := *d.DeepCopy()
	dCopy.Labels = hash.SetTemplateHashLabel(dCopy.Labels, dCopy)
	return dCopy
}
