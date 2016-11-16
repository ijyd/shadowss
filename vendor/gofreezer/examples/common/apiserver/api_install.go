package apiserver

import (
	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
)

func (apis *APIServer) installSwaggerAPI(container *restful.Container) {
	hostAndPort := apis.Host + string(":") + strconv.Itoa(apis.Port)
	//protocol := "https://"
	protocol := "http://"
	webServicesUrl := protocol + hostAndPort

	// Enable swagger UI and discovery API
	swaggerConfig := swagger.Config{
		WebServicesUrl:  webServicesUrl,
		WebServices:     container.RegisteredWebServices(),
		ApiPath:         "/swaggerapi/",
		SwaggerPath:     "/swaggerui/",
		SwaggerFilePath: apis.SwaggerPath,
		SchemaFormatHandler: func(typeName string) string {
			switch typeName {
			case "unversioned.Time", "*unversioned.Time":
				return "date-time"
			}
			return ""
		},
	}
	swagger.RegisterSwaggerService(swaggerConfig, container)
}

func (apis *APIServer) install(container *restful.Container) error {

	if len(apis.SwaggerPath) > 0 {
		apis.installSwaggerAPI(container)
	}

	return nil
}
