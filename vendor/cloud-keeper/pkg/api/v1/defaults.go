package v1

import "gofreezer/pkg/runtime"

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return scheme.AddDefaultingFuncs(SetDefaults_NodeUser,
		SetDefaults_APIServer)
}

func SetDefaults_NodeUser(obj *NodeUser) {
	obj.Spec.User.Port = 0
}

func SetDefaults_APIServer(obj *APIServer) {
	obj.Spec.Server.Port = 0
}

func SetDefaults_UserService(obj *UserService) {
	obj.Spec.NodeCnt = 0
}
