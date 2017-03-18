// +build !ignore_autogenerated

/*
Copyright 2016 The cloud-keeper Authors.
*/

// This file was autogenerated by conversion-gen. Do not edit it manually!

package v1alpha1

import (
	api "cloud-keeper/pkg/api"
	v1 "cloud-keeper/pkg/api/v1"
	batch "cloud-keeper/pkg/apis/batch"
	conversion "gofreezer/pkg/conversion"
	runtime "gofreezer/pkg/runtime"
)

func init() {
	SchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(scheme *runtime.Scheme) error {
	return scheme.AddGeneratedConversionFuncs(
		Convert_v1alpha1_BatchAccServer_To_batch_BatchAccServer,
		Convert_batch_BatchAccServer_To_v1alpha1_BatchAccServer,
		Convert_v1alpha1_BatchAccServerSpec_To_batch_BatchAccServerSpec,
		Convert_batch_BatchAccServerSpec_To_v1alpha1_BatchAccServerSpec,
		Convert_v1alpha1_BatchResumeUsers_To_batch_BatchResumeUsers,
		Convert_batch_BatchResumeUsers_To_v1alpha1_BatchResumeUsers,
		Convert_v1alpha1_BatchResumeUsersSpec_To_batch_BatchResumeUsersSpec,
		Convert_batch_BatchResumeUsersSpec_To_v1alpha1_BatchResumeUsersSpec,
	)
}

func autoConvert_v1alpha1_BatchAccServer_To_batch_BatchAccServer(in *BatchAccServer, out *batch.BatchAccServer, s conversion.Scope) error {
	// TODO: Inefficient conversion - can we improve it?
	if err := s.Convert(&in.ObjectMeta, &out.ObjectMeta, 0); err != nil {
		return err
	}
	if err := Convert_v1alpha1_BatchAccServerSpec_To_batch_BatchAccServerSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

func Convert_v1alpha1_BatchAccServer_To_batch_BatchAccServer(in *BatchAccServer, out *batch.BatchAccServer, s conversion.Scope) error {
	return autoConvert_v1alpha1_BatchAccServer_To_batch_BatchAccServer(in, out, s)
}

func autoConvert_batch_BatchAccServer_To_v1alpha1_BatchAccServer(in *batch.BatchAccServer, out *BatchAccServer, s conversion.Scope) error {
	// TODO: Inefficient conversion - can we improve it?
	if err := s.Convert(&in.ObjectMeta, &out.ObjectMeta, 0); err != nil {
		return err
	}
	if err := Convert_batch_BatchAccServerSpec_To_v1alpha1_BatchAccServerSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

func Convert_batch_BatchAccServer_To_v1alpha1_BatchAccServer(in *batch.BatchAccServer, out *BatchAccServer, s conversion.Scope) error {
	return autoConvert_batch_BatchAccServer_To_v1alpha1_BatchAccServer(in, out, s)
}

func autoConvert_v1alpha1_BatchAccServerSpec_To_batch_BatchAccServerSpec(in *BatchAccServerSpec, out *batch.BatchAccServerSpec, s conversion.Scope) error {
	if in.ServerList != nil {
		in, out := &in.ServerList, &out.ServerList
		*out = make([]api.AccServer, len(*in))
		for i := range *in {
			// TODO: Inefficient conversion - can we improve it?
			if err := s.Convert(&(*in)[i], &(*out)[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.ServerList = nil
	}
	return nil
}

func Convert_v1alpha1_BatchAccServerSpec_To_batch_BatchAccServerSpec(in *BatchAccServerSpec, out *batch.BatchAccServerSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_BatchAccServerSpec_To_batch_BatchAccServerSpec(in, out, s)
}

func autoConvert_batch_BatchAccServerSpec_To_v1alpha1_BatchAccServerSpec(in *batch.BatchAccServerSpec, out *BatchAccServerSpec, s conversion.Scope) error {
	if in.ServerList != nil {
		in, out := &in.ServerList, &out.ServerList
		*out = make([]v1.AccServer, len(*in))
		for i := range *in {
			// TODO: Inefficient conversion - can we improve it?
			if err := s.Convert(&(*in)[i], &(*out)[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.ServerList = nil
	}
	return nil
}

func Convert_batch_BatchAccServerSpec_To_v1alpha1_BatchAccServerSpec(in *batch.BatchAccServerSpec, out *BatchAccServerSpec, s conversion.Scope) error {
	return autoConvert_batch_BatchAccServerSpec_To_v1alpha1_BatchAccServerSpec(in, out, s)
}

func autoConvert_v1alpha1_BatchResumeUsers_To_batch_BatchResumeUsers(in *BatchResumeUsers, out *batch.BatchResumeUsers, s conversion.Scope) error {
	// TODO: Inefficient conversion - can we improve it?
	if err := s.Convert(&in.ObjectMeta, &out.ObjectMeta, 0); err != nil {
		return err
	}
	if err := Convert_v1alpha1_BatchResumeUsersSpec_To_batch_BatchResumeUsersSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

func Convert_v1alpha1_BatchResumeUsers_To_batch_BatchResumeUsers(in *BatchResumeUsers, out *batch.BatchResumeUsers, s conversion.Scope) error {
	return autoConvert_v1alpha1_BatchResumeUsers_To_batch_BatchResumeUsers(in, out, s)
}

func autoConvert_batch_BatchResumeUsers_To_v1alpha1_BatchResumeUsers(in *batch.BatchResumeUsers, out *BatchResumeUsers, s conversion.Scope) error {
	// TODO: Inefficient conversion - can we improve it?
	if err := s.Convert(&in.ObjectMeta, &out.ObjectMeta, 0); err != nil {
		return err
	}
	if err := Convert_batch_BatchResumeUsersSpec_To_v1alpha1_BatchResumeUsersSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

func Convert_batch_BatchResumeUsers_To_v1alpha1_BatchResumeUsers(in *batch.BatchResumeUsers, out *BatchResumeUsers, s conversion.Scope) error {
	return autoConvert_batch_BatchResumeUsers_To_v1alpha1_BatchResumeUsers(in, out, s)
}

func autoConvert_v1alpha1_BatchResumeUsersSpec_To_batch_BatchResumeUsersSpec(in *BatchResumeUsersSpec, out *batch.BatchResumeUsersSpec, s conversion.Scope) error {
	out.SchedulingTime = in.SchedulingTime
	return nil
}

func Convert_v1alpha1_BatchResumeUsersSpec_To_batch_BatchResumeUsersSpec(in *BatchResumeUsersSpec, out *batch.BatchResumeUsersSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_BatchResumeUsersSpec_To_batch_BatchResumeUsersSpec(in, out, s)
}

func autoConvert_batch_BatchResumeUsersSpec_To_v1alpha1_BatchResumeUsersSpec(in *batch.BatchResumeUsersSpec, out *BatchResumeUsersSpec, s conversion.Scope) error {
	out.SchedulingTime = in.SchedulingTime
	return nil
}

func Convert_batch_BatchResumeUsersSpec_To_v1alpha1_BatchResumeUsersSpec(in *batch.BatchResumeUsersSpec, out *BatchResumeUsersSpec, s conversion.Scope) error {
	return autoConvert_batch_BatchResumeUsersSpec_To_v1alpha1_BatchResumeUsersSpec(in, out, s)
}