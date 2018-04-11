package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/Encrypt-S/navpi-go/app/conf"
)

// InitMetaHandlers starts the meta api handlers
func InitMetaHandlers(r *mux.Router, prefix string) {
	nameSpace := "meta"

	metaErrorCodePath := RouteBuilder(prefix, nameSpace, "v1", "errorcodes")
	OpenRouteHandler(metaErrorCodePath, r, metaErrorDisplayHandler())

	metaCoinPath := RouteBuilder(prefix, nameSpace, "v1", "coins")
	OpenRouteHandler(metaCoinPath, r, coinMetaHandler())

}


// metaErrorDisplayHandler displays all the application errors to frontend
func coinMetaHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appResp := Response{}
		appResp.Data = conf.AppConf.Coins
		appResp.Send(w)

	})
}

// metaErrorDisplayHandler displays all the application errors to frontend
func metaErrorDisplayHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appResp := Response{}
		appResp.Data = AppRespErrors
		appResp.Send(w)

	})
}
