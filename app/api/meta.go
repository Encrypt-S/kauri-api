package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

// InitMetaHandlers starts the meta api handlers
func InitMetaHandlers(r *mux.Router, prefix string) {
	nameSpace := "meta"

	metaErrorCodePath := RouteBuilder(prefix, nameSpace, "v1", "errorcodes")
	OpenRouteHandler(metaErrorCodePath, r, metaErrorDisplayHandler())

}

// metaErrorDisplayHandler displays all the application errors to frontend
func metaErrorDisplayHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appResp := Response{}
		appResp.Data = AppRespErrors
		appResp.Send(w)

	})
}
