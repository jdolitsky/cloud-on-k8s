// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package sset

import (
	"github.com/elastic/cloud-on-k8s/pkg/apis/elasticsearch/v1beta1"
	"github.com/elastic/cloud-on-k8s/pkg/controller/common/hash"
	"github.com/elastic/cloud-on-k8s/pkg/controller/elasticsearch/label"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type TestSset struct {
	Namespace   string
	Name        string
	ClusterName string
	Version     string
	Replicas    int32
	Master      bool
	Data        bool
	Ingest      bool
	Status      appsv1.StatefulSetStatus
}

func (t TestSset) Pods() []runtime.Object {
	podNames := PodNames(t.Build())
	pods := make([]runtime.Object, t.Replicas)
	for i, podName := range podNames {
		pods[i] = TestPod{
			Namespace:       t.Namespace,
			Name:            podName,
			StatefulSetName: t.Name,
			Master:          t.Master,
			Data:            t.Data,
			Ingest:          t.Ingest,
			Version:         t.Version,
			ClusterName:     t.ClusterName,
		}.BuildPtr()
	}
	return pods
}

func (t TestSset) Build() appsv1.StatefulSet {
	labels := map[string]string{
		label.VersionLabelName:     t.Version,
		label.ClusterNameLabelName: t.ClusterName,
	}
	label.NodeTypesMasterLabelName.Set(t.Master, labels)
	label.NodeTypesDataLabelName.Set(t.Data, labels)
	label.NodeTypesIngestLabelName.Set(t.Ingest, labels)
	statefulSet := appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      t.Name,
			Namespace: t.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &t.Replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
			},
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: "OnDelete",
			},
		},
		Status: t.Status,
	}
	statefulSet.Labels = hash.SetTemplateHashLabel(statefulSet.Labels, statefulSet.Spec)
	return statefulSet
}

func (t TestSset) BuildPtr() *appsv1.StatefulSet {
	built := t.Build()
	return &built
}

type TestPod struct {
	Namespace       string
	Name            string
	ClusterName     string
	StatefulSetName string
	Version         string
	Revision        string
	Master          bool
	Data            bool
	Ingest          bool
	Ready           bool
	RestartCount    int32
	Phase           corev1.PodPhase
}

func (t TestPod) Build() corev1.Pod {
	labels := map[string]string{
		label.VersionLabelName:          t.Version,
		label.ClusterNameLabelName:      t.ClusterName,
		label.StatefulSetNameLabelName:  t.StatefulSetName,
		appsv1.StatefulSetRevisionLabel: t.Revision,
	}
	label.NodeTypesMasterLabelName.Set(t.Master, labels)
	label.NodeTypesDataLabelName.Set(t.Data, labels)
	label.NodeTypesIngestLabelName.Set(t.Ingest, labels)

	status := corev1.PodStatus{
		// assume Running by default
		Phase: corev1.PodRunning,
	}
	// unless specified otherwise
	if t.Phase != "" {
		status.Phase = t.Phase
	}
	if t.Ready {
		status.Conditions = []corev1.PodCondition{
			{
				Status: corev1.ConditionTrue,
				Type:   corev1.ContainersReady,
			},
			{
				Status: corev1.ConditionTrue,
				Type:   corev1.PodReady,
			},
		}
	}
	status.ContainerStatuses = []corev1.ContainerStatus{
		{
			Name:         v1beta1.ElasticsearchContainerName,
			RestartCount: t.RestartCount,
			Ready:        t.Ready,
		},
	}
	return corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: t.Namespace,
			Name:      t.Name,
			Labels:    labels,
		},
		Status: status,
	}
}

func (t TestPod) BuildPtr() *corev1.Pod {
	pod := t.Build()
	return &pod
}
