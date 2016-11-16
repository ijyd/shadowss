package api

import freezerapi "gofreezer/pkg/api"

//TODO import ObjectMeta and other type into here

var (
	Resource             = freezerapi.Resource
	ParameterCodec       = freezerapi.ParameterCodec
	NewDefaultRESTMapper = freezerapi.NewDefaultRESTMapper
	SimpleNameGenerator  = freezerapi.SimpleNameGenerator
)

// type Context interface {
// 	freezerapi.Context
// }
//
// type NameGenerator interface {
// 	freezerapi.NameGenerator
// }
