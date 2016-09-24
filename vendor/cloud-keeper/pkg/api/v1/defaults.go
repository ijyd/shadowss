package v1

import "gofreezer/pkg/runtime"

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return scheme.AddDefaultingFuncs(SetDefaults_NodeUser,
		SetDefaults_APIServer)
}

func SetDefaults_NodeUser(obj *NodeUser) {
}

func SetDefaults_APIServer(obj *APIServer) {
}

func SetDefaults_UserService(obj *UserService) {
}

func SetDefaults_Node(obj *Node) {
	if len(obj.Spec.Server.Method) == 0 {
		obj.Spec.Server.Method = "aes-256-cfb"
	}
}
