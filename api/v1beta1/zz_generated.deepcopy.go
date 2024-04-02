//go:build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KafkaSchema) DeepCopyInto(out *KafkaSchema) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KafkaSchema.
func (in *KafkaSchema) DeepCopy() *KafkaSchema {
	if in == nil {
		return nil
	}
	out := new(KafkaSchema)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KafkaSchema) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KafkaSchemaData) DeepCopyInto(out *KafkaSchemaData) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KafkaSchemaData.
func (in *KafkaSchemaData) DeepCopy() *KafkaSchemaData {
	if in == nil {
		return nil
	}
	out := new(KafkaSchemaData)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KafkaSchemaList) DeepCopyInto(out *KafkaSchemaList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]KafkaSchema, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KafkaSchemaList.
func (in *KafkaSchemaList) DeepCopy() *KafkaSchemaList {
	if in == nil {
		return nil
	}
	out := new(KafkaSchemaList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KafkaSchemaList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KafkaSchemaSpec) DeepCopyInto(out *KafkaSchemaSpec) {
	*out = *in
	out.SchemaRegistry = in.SchemaRegistry
	out.Data = in.Data
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KafkaSchemaSpec.
func (in *KafkaSchemaSpec) DeepCopy() *KafkaSchemaSpec {
	if in == nil {
		return nil
	}
	out := new(KafkaSchemaSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KafkaSchemaStatus) DeepCopyInto(out *KafkaSchemaStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KafkaSchemaStatus.
func (in *KafkaSchemaStatus) DeepCopy() *KafkaSchemaStatus {
	if in == nil {
		return nil
	}
	out := new(KafkaSchemaStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReadyReason) DeepCopyInto(out *ReadyReason) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReadyReason.
func (in *ReadyReason) DeepCopy() *ReadyReason {
	if in == nil {
		return nil
	}
	out := new(ReadyReason)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SchemaRegistry) DeepCopyInto(out *SchemaRegistry) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SchemaRegistry.
func (in *SchemaRegistry) DeepCopy() *SchemaRegistry {
	if in == nil {
		return nil
	}
	out := new(SchemaRegistry)
	in.DeepCopyInto(out)
	return out
}