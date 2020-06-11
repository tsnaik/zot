package api

import (
	"fmt"
	"net/http"

	"github.com/anuvu/zot/pkg/extensions/accesscontrol"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
)

func AccessHandler(c *Controller) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username := fmt.Sprintf("%v", context.Get(r, "username"))
			if accesscontrol.IsAuthorized(username, r.Method, r.RequestURI, c.Config.HTTP.AccessControlConfigPath, c.Log) {
				next.ServeHTTP(w, r)
			} else {
				var json = jsoniter.ConfigCompatibleWithStandardLibrary
				data, err := json.Marshal(NewError(DENIED))

				if err != nil {
					c.Log.Panic().Err(err).Msg("Error marshalling json: " + err.Error())
				}
				http.Error(w, string(data), http.StatusForbidden)
			}
		})
	}
}
