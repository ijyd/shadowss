// +build !ignore_autogenerated

/*
Copyright 2016 The cloud-keeper Authors.
*/

// This file was autogenerated by deepcopy-gen. Do not edit it manually!

package v1alpha1

import (
	api_v1 "cloud-keeper/pkg/api/v1"
	v1 "gofreezer/pkg/api/v1"
	conversion "gofreezer/pkg/conversion"
	runtime "gofreezer/pkg/runtime"
	reflect "reflect"
)

func init() {
	SchemeBuilder.Register(RegisterDeepCopies)
}

// RegisterDeepCopies adds deep-copy functions to the given scheme. Public
// to allow building arbitrary schemes.
func RegisterDeepCopies(scheme *runtime.Scheme) error {
	return scheme.AddGeneratedDeepCopyFuncs(
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_v1alpha1_BatchAccServer, InType: reflect.TypeOf(&BatchAccServer{})},
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_v1alpha1_BatchAccServerSpec, InType: reflect.TypeOf(&BatchAccServerSpec{})},
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_v1alpha1_BatchResumeUsers, InType: reflect.TypeOf(&BatchResumeUsers{})},
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_v1alpha1_BatchResumeUsersSpec, InType: reflect.TypeOf(&BatchResumeUsersSpec{})},
	)
}

func DeepCopy_v1alpha1_BatchAccServer(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*BatchAccServer)
		out := out.(*BatchAccServer)
		out.TypeMeta = in.TypeMeta
		if newVal, err := c.DeepCopy(&in.ObjectMeta); err != nil {
			return err
		} else {
			out.ObjectMeta = *newVal.(*v1.ObjectMeta)
		}
		if newVal, err := c.DeepCopy(&in.Spec); err != nil {
			return err
		} else {
			out.Spec = *newVal.(*BatchAccServerSpec)
		}
		return nil
	}
}

func DeepCopy_v1alpha1_BatchAccServerSpec(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*BatchAccServerSpec)
		out := out.(*BatchAccServerSpec)
		if in.ServerList != nil {
			in, out := &in.ServerList, &out.ServerList
			*out = make([]api_v1.AccServer, len(*in))
			for i := range *in {
				if err := api_v1.DeepCopy_v1_AccServer(&(*in)[i], &(*out)[i], c); err != nil {
					return err
				}
			}
		} else {
			out.ServerList = nil
		}
		return nil
	}
}

func DeepCopy_v1alpha1_BatchResumeUsers(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*BatchResumeUsers)
		out := out.(*BatchResumeUsers)
		out.TypeMeta = in.TypeMeta
		if newVal, err := c.DeepCopy(&in.ObjectMeta); err != nil {
			return err
		} else {
			out.ObjectMeta = *newVal.(*v1.ObjectMeta)
		}
		if newVal, err := c.DeepCopy(&in.Spec); err != nil {
			return err
		} else {
			out.Spec = *newVal.(*BatchResumeUsersSpec)
		}
		return nil
	}
}

func DeepCopy_v1alpha1_BatchResumeUsersSpec(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*BatchResumeUsersSpec)
		out := out.(*BatchResumeUsersSpec)
		out.SchedulingTime = in.SchedulingTime.DeepCopy()
		return nil
	}
}