package app

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func (c *App) MatrixAdminProxy() http.HandlerFunc {

	endpoint := fmt.Sprintf("http://%s/_synapse/", c.Config.Matrix.Homeserver)

	target, _ := url.Parse(endpoint)

	proxy := httputil.NewSingleHostReverseProxy(target)

	return func(w http.ResponseWriter, r *http.Request) {

		user_id := c.AuthenticatedUser(r)

		access_token := c.AuthenticatedAccessToken(r)

		c.Log.Debug().
			Msg(fmt.Sprintf("Admin user %s", *user_id))

		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *access_token))
		w.Header().Del("Access-Control-Allow-Origin")

		proxy.ServeHTTP(w, r)
	}

}
