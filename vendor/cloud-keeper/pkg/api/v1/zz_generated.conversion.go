package v1

import (
	"cloud-keeper/pkg/api"

	"gofreezer/pkg/api/prototype"
	"gofreezer/pkg/conversion"
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/types"
)

func init() {
	SchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(scheme *runtime.Scheme) error {
	return scheme.AddGeneratedConversionFuncs(
		Convert_api_ObjectMeta_To_v1_ObjectMeta,
		Convert_v1_ObjectMeta_To_api_ObjectMeta,
		Convert_v1_OwnerReference_To_api_OwnerReference,
		Convert_api_OwnerReference_To_v1_OwnerReference,
	)
}

func autoConvert_v1_ObjectMeta_To_api_ObjectMeta(in *ObjectMeta, out *prototype.ObjectMeta, s conversion.Scope) error {
	out.Name = in.Name
	out.GenerateName = in.GenerateName
	out.Namespace = in.Namespace
	out.SelfLink = in.SelfLink
	out.UID = types.UID(in.UID)
	out.ResourceVersion = in.ResourceVersion
	out.Generation = in.Generation
	if err := prototype.Convert_unversioned_Time_To_unversioned_Time(&in.CreationTimestamp, &out.CreationTimestamp, s); err != nil {
		return err
	}
	out.DeletionTimestamp = in.DeletionTimestamp
	out.DeletionGracePeriodSeconds = in.DeletionGracePeriodSeconds
	out.Labels = in.Labels
	out.Annotations = in.Annotations
	if in.OwnerReferences != nil {
		in, out := &in.OwnerReferences, &out.OwnerReferences
		*out = make([]prototype.OwnerReference, len(*in))
		for i := range *in {
			if err := Convert_v1_OwnerReference_To_api_OwnerReference(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.OwnerReferences = nil
	}
	out.Finalizers = in.Finalizers
	out.ClusterName = in.ClusterName
	return nil
}

func Convert_v1_ObjectMeta_To_api_ObjectMeta(in *ObjectMeta, out *prototype.ObjectMeta, s conversion.Scope) error {
	return autoConvert_v1_ObjectMeta_To_api_ObjectMeta(in, out, s)
}

func autoConvert_api_ObjectMeta_To_v1_ObjectMeta(in *prototype.ObjectMeta, out *ObjectMeta, s conversion.Scope) error {
	out.Name = in.Name
	out.GenerateName = in.GenerateName
	out.Namespace = in.Namespace
	out.SelfLink = in.SelfLink
	out.UID = types.UID(in.UID)
	out.ResourceVersion = in.ResourceVersion
	out.Generation = in.Generation
	if err := prototype.Convert_unversioned_Time_To_unversioned_Time(&in.CreationTimestamp, &out.CreationTimestamp, s); err != nil {
		return err
	}
	out.DeletionTimestamp = in.DeletionTimestamp
	out.DeletionGracePeriodSeconds = in.DeletionGracePeriodSeconds
	out.Labels = in.Labels
	out.Annotations = in.Annotations
	if in.OwnerReferences != nil {
		in, out := &in.OwnerReferences, &out.OwnerReferences
		*out = make([]OwnerReference, len(*in))
		for i := range *in {
			if err := Convert_api_OwnerReference_To_v1_OwnerReference(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.OwnerReferences = nil
	}
	out.Finalizers = in.Finalizers
	out.ClusterName = in.ClusterName
	return nil
}

func Convert_api_ObjectMeta_To_v1_ObjectMeta(in *prototype.ObjectMeta, out *ObjectMeta, s conversion.Scope) error {
	return autoConvert_api_ObjectMeta_To_v1_ObjectMeta(in, out, s)
}

func autoConvert_v1_OwnerReference_To_api_OwnerReference(in *OwnerReference, out *prototype.OwnerReference, s conversion.Scope) error {
	out.APIVersion = in.APIVersion
	out.Kind = in.Kind
	out.Name = in.Name
	out.UID = types.UID(in.UID)
	out.Controller = in.Controller
	return nil
}

func Convert_v1_OwnerReference_To_api_OwnerReference(in *OwnerReference, out *prototype.OwnerReference, s conversion.Scope) error {
	return autoConvert_v1_OwnerReference_To_api_OwnerReference(in, out, s)
}

func autoConvert_api_OwnerReference_To_v1_OwnerReference(in *prototype.OwnerReference, out *OwnerReference, s conversion.Scope) error {
	out.APIVersion = in.APIVersion
	out.Kind = in.Kind
	out.Name = in.Name
	out.UID = types.UID(in.UID)
	out.Controller = in.Controller
	return nil
}

func Convert_api_OwnerReference_To_v1_OwnerReference(in *prototype.OwnerReference, out *OwnerReference, s conversion.Scope) error {
	return autoConvert_api_OwnerReference_To_v1_OwnerReference(in, out, s)
}

func autoConvert_api_NodeUser_To_v1_NodeUser(in *api.NodeUser, out *NodeUser, s conversion.Scope) error {
	if err := prototype.Convert_unversioned_TypeMeta_To_unversioned_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := Convert_api_ObjectMeta_To_v1_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	if err := autoConvert_api_NodeUserSpec_To_v1_NodeUserSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1_NodeUser_To_api_NodeUser(in *NodeUser, out *api.NodeUser, s conversion.Scope) error {
	SetDefaults_NodeUser(in)
	if err := prototype.Convert_unversioned_TypeMeta_To_unversioned_TypeMeta(&in.TypeMeta, &out.TypeMeta, s); err != nil {
		return err
	}
	if err := Convert_v1_ObjectMeta_To_api_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, s); err != nil {
		return err
	}
	if err := autoConvert_v1_NodeUserSpec_To_api_NodeUserSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}

	return nil
}

func autoConvert_v1_NodeUserSpec_To_api_NodeUserSpec(in *NodeUserSpec, out *api.NodeUserSpec, s conversion.Scope) error {

	out.User.ID = in.User.ID
	out.User.Method = in.User.Method
	out.User.DownloadTraffic = in.User.DownloadTraffic
	out.User.UploadTraffic = in.User.UploadTraffic
	out.User.EnableOTA = in.User.EnableOTA
	out.User.Password = in.User.Password
	out.User.Port = in.User.Port
	out.NodeName = in.NodeName

	return nil
}

func autoConvert_api_NodeUserSpec_To_v1_NodeUserSpec(in *api.NodeUserSpec, out *NodeUserSpec, s conversion.Scope) error {
	//out.User = NodeUserSpec(in.User)
	out.User.ID = in.User.ID
	out.User.Method = in.User.Method
	out.User.DownloadTraffic = in.User.DownloadTraffic
	out.User.UploadTraffic = in.User.UploadTraffic
	out.User.EnableOTA = in.User.EnableOTA
	out.User.Password = in.User.Password
	out.User.Port = in.User.Port
	out.NodeName = in.NodeName
	return nil
}

func autoConvert_v1_APIServerSpec_To_api_APIServerSpec(in *APIServerSpec, out *api.APIServerSpec, s conversion.Scope) error {

	out.Server.ID = in.Server.ID
	out.Server.Name = in.Server.Name
	out.Server.Host = in.Server.Host
	out.Server.Port = in.Server.Port
	out.Server.Status = in.Server.Status
	out.Server.CreateTime = in.Server.CreateTime

	return nil
}

func autoConvert_api_APIServerSpec_To_v1_APIServerSpec(in *api.APIServerSpec, out *APIServerSpec, s conversion.Scope) error {
	out.Server.ID = in.Server.ID
	out.Server.Name = in.Server.Name
	out.Server.Host = in.Server.Host
	out.Server.Port = in.Server.Port
	out.Server.Status = in.Server.Status
	out.Server.CreateTime = in.Server.CreateTime
	return nil
}

func autoConvert_v1_UserServiceSpec_To_api_UserServiceSpec(in *UserServiceSpec, out *api.UserServiceSpec, s conversion.Scope) error {

	out.NodeCnt = in.NodeCnt
	if in.NodeUserReference != nil {
		out.NodeUserReference = make(map[string]api.UserReferences)
		for k, v := range in.NodeUserReference {
			outVal := api.UserReferences{
				ID:              v.ID,
				Port:            v.Port,
				Method:          v.Method,
				Password:        v.Password,
				EnableOTA:       v.EnableOTA,
				UploadTraffic:   v.UploadTraffic,
				DownloadTraffic: v.DownloadTraffic,
			}
			out.NodeUserReference[k] = outVal
		}
	} else {
		out.NodeUserReference = nil
	}

	return nil
}

func autoConvert_api_UserServiceSpec_To_v1_UserServiceSpec(in *api.UserServiceSpec, out *UserServiceSpec, s conversion.Scope) error {
	out.NodeCnt = in.NodeCnt
	if in.NodeUserReference != nil {
		out.NodeUserReference = make(map[string]UserReferences)
		for k, v := range in.NodeUserReference {
			outVal := UserReferences{
				ID:              v.ID,
				Port:            v.Port,
				Method:          v.Method,
				Password:        v.Password,
				EnableOTA:       v.EnableOTA,
				UploadTraffic:   v.UploadTraffic,
				DownloadTraffic: v.DownloadTraffic,
			}
			out.NodeUserReference[k] = outVal
		}
	} else {
		out.NodeUserReference = nil
	}

	return nil
}
