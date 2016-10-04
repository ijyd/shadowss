package v1

import (
	"cloud-keeper/pkg/api"

	"github.com/golang/glog"

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

	autoConvert_v1_UserReferences_To_api_UserReferences(&in.User, &out.User, s)
	out.NodeName = in.NodeName
	out.Phase = api.NodeUserPhase(in.Phase)

	return nil
}

func autoConvert_api_NodeUserSpec_To_v1_NodeUserSpec(in *api.NodeUserSpec, out *NodeUserSpec, s conversion.Scope) error {
	autoConvert_api_UserReferences_To_v1_UserReferences(&in.User, &out.User, s)

	out.NodeName = in.NodeName
	out.Phase = NodeUserPhase(in.Phase)

	glog.V(5).Infof("convert %+v to %+v", *in, *out)
	return nil
}

func autoConvert_v1_APIServerSpec_To_api_APIServerSpec(in *APIServerSpec, out *api.APIServerSpec, s conversion.Scope) error {

	out.Server.ID = in.Server.ID
	out.Server.Name = in.Server.Name

	out.Server.Port = in.Server.Port
	out.Server.Status = in.Server.Status
	out.Server.CreateTime = in.Server.CreateTime
	out.Server.Host = in.Server.Host

	if in.HostList != nil {
		in, out := &in.HostList, &out.HostList
		*out = make([]string, len(*in))
		for i := range *in {
			(*out)[i] = (*in)[i]
		}
	}

	return nil
}

func autoConvert_api_APIServerSpec_To_v1_APIServerSpec(in *api.APIServerSpec, out *APIServerSpec, s conversion.Scope) error {
	out.Server.ID = in.Server.ID
	out.Server.Name = in.Server.Name
	out.Server.Port = in.Server.Port
	out.Server.Status = in.Server.Status
	out.Server.CreateTime = in.Server.CreateTime
	out.Server.Host = in.Server.Host

	if in.HostList != nil {
		in, out := &in.HostList, &out.HostList
		*out = make([]string, len(*in))
		for i := range *in {
			(*out)[i] = (*in)[i]
		}
	}
	return nil
}

func autoConvert_v1_UserServiceSpec_To_api_UserServiceSpec(in *UserServiceSpec, out *api.UserServiceSpec, s conversion.Scope) error {

	out.NodeCnt = in.NodeCnt
	if in.NodeUserReference != nil {
		out.NodeUserReference = make(map[string]api.NodeReferences)
		for k, v := range in.NodeUserReference {
			var outVal api.UserReferences
			autoConvert_v1_UserReferences_To_api_UserReferences(&v.User, &outVal, s)
			nodeRefer := api.NodeReferences{
				User: outVal,
				Host: v.Host,
			}
			out.NodeUserReference[k] = nodeRefer
		}
	} else {
		out.NodeUserReference = nil
	}

	glog.V(5).Infof("in user %v out %v", *in, *out)
	return nil
}

func autoConvert_api_UserServiceSpec_To_v1_UserServiceSpec(in *api.UserServiceSpec, out *UserServiceSpec, s conversion.Scope) error {
	out.NodeCnt = in.NodeCnt
	if in.NodeUserReference != nil {
		out.NodeUserReference = make(map[string]NodeReferences)
		for k, v := range in.NodeUserReference {
			var outVal UserReferences
			autoConvert_api_UserReferences_To_v1_UserReferences(&v.User, &outVal, s)
			nodeRefer := NodeReferences{
				User: outVal,
				Host: v.Host,
			}
			out.NodeUserReference[k] = nodeRefer
		}
	} else {
		out.NodeUserReference = nil
	}

	glog.V(5).Infof("in user %v out %v", *in, *out)

	return nil
}

func autoConvert_v1_NodeSpec_To_api_NodeSpec(in *NodeSpec, out *api.NodeSpec, s conversion.Scope) error {

	glog.Infof("convert v1 node spec to api nodespec")
	out.Server = api.NodeServer{
		ID:            in.Server.ID,
		Name:          in.Server.Name,
		EnableOTA:     in.Server.EnableOTA,
		Host:          in.Server.Host,
		Method:        in.Server.Method,
		Status:        in.Server.Status,
		Location:      in.Server.Location,
		AccServerID:   in.Server.AccServerID,
		AccServerName: in.Server.AccServerName,
		Description:   in.Server.Description,
		TrafficLimit:  in.Server.TrafficLimit,
		Upload:        in.Server.Upload,
		Download:      in.Server.Download,
		TrafficRate:   in.Server.TrafficRate,
	}

	return nil
}

func autoConvert_api_NodeSpec_To_v1_NodeSpec(in *api.NodeSpec, out *NodeSpec, s conversion.Scope) error {

	glog.Infof("convert ai node spec to v1 nodespec")
	out.Server = NodeServer{
		ID:            in.Server.ID,
		Name:          in.Server.Name,
		EnableOTA:     in.Server.EnableOTA,
		Host:          in.Server.Host,
		Method:        in.Server.Method,
		Status:        in.Server.Status,
		Location:      in.Server.Location,
		AccServerID:   in.Server.AccServerID,
		AccServerName: in.Server.AccServerName,
		Description:   in.Server.Description,
		TrafficLimit:  in.Server.TrafficLimit,
		Upload:        in.Server.Upload,
		Download:      in.Server.Download,
		TrafficRate:   in.Server.TrafficRate,
	}

	return nil
}

func autoConvert_api_UserReferences_To_v1_UserReferences(in *api.UserReferences, out *UserReferences, s conversion.Scope) error {

	if in != nil {
		*out = UserReferences{
			ID:              in.ID,
			Port:            in.Port,
			Method:          in.Method,
			Password:        in.Password,
			Name:            in.Name,
			EnableOTA:       in.EnableOTA,
			UploadTraffic:   in.UploadTraffic,
			DownloadTraffic: in.DownloadTraffic,
		}
	} else {
		out = nil
	}
	return nil
}

func autoConvert_v1_UserReferences_To_api_UserReferences(in *UserReferences, out *api.UserReferences, s conversion.Scope) error {
	if in != nil {
		*out = api.UserReferences{
			ID:              in.ID,
			Port:            in.Port,
			Method:          in.Method,
			Password:        in.Password,
			Name:            in.Name,
			EnableOTA:       in.EnableOTA,
			UploadTraffic:   in.UploadTraffic,
			DownloadTraffic: in.DownloadTraffic,
		}
	} else {
		out = nil
	}
	return nil
}
