// +build !ignore_autogenerated

// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	commonv1alpha1 "github.com/elastic/cloud-on-k8s/pkg/apis/common/v1alpha1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Logstash) DeepCopyInto(out *Logstash) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	if in.assocConf != nil {
		in, out := &in.assocConf, &out.assocConf
		*out = new(commonv1alpha1.AssociationConf)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Logstash.
func (in *Logstash) DeepCopy() *Logstash {
	if in == nil {
		return nil
	}
	out := new(Logstash)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Logstash) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LogstashList) DeepCopyInto(out *LogstashList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Logstash, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LogstashList.
func (in *LogstashList) DeepCopy() *LogstashList {
	if in == nil {
		return nil
	}
	out := new(LogstashList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LogstashList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LogstashSpec) DeepCopyInto(out *LogstashSpec) {
	*out = *in
	out.ElasticsearchRef = in.ElasticsearchRef
	if in.Config != nil {
		in, out := &in.Config, &out.Config
		*out = (*in).DeepCopy()
	}
	in.HTTP.DeepCopyInto(&out.HTTP)
	in.PodTemplate.DeepCopyInto(&out.PodTemplate)
	if in.SecureSettings != nil {
		in, out := &in.SecureSettings, &out.SecureSettings
		*out = make([]commonv1alpha1.SecretSource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LogstashSpec.
func (in *LogstashSpec) DeepCopy() *LogstashSpec {
	if in == nil {
		return nil
	}
	out := new(LogstashSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LogstashStatus) DeepCopyInto(out *LogstashStatus) {
	*out = *in
	out.ReconcilerStatus = in.ReconcilerStatus
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LogstashStatus.
func (in *LogstashStatus) DeepCopy() *LogstashStatus {
	if in == nil {
		return nil
	}
	out := new(LogstashStatus)
	in.DeepCopyInto(out)
	return out
}
