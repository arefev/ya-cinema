package handler

import (
	"math/rand/v2"
	"net/http"
	"net/http/httputil"
	"net/url"

	"proxy/internal/application"
)

type proxyHandler struct {
	app                 *application.App
	target              string
	useGradualMigration bool
}

func NewProxyHandler(app *application.App, target string) *proxyHandler {
	return &proxyHandler{app: app, target: target}
}

func (ph *proxyHandler) Proxy(w http.ResponseWriter, r *http.Request) {
	target, err := url.Parse(ph.getTarget())
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	r.URL.Host = target.Host
	r.URL.Scheme = target.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = target.Host

	ph.app.Log.Sugar().Infof("Proxing request to %s with method %s", r.URL, r.Method)
	proxy.ServeHTTP(w, r)
}

func (ph *proxyHandler) UseGradualMigration() {
	ph.useGradualMigration = true
}

func (ph *proxyHandler) getTarget() string {
	if !ph.useGradualMigration {
		return ph.target
	}

	percent := rand.IntN(100)
	migrationPercent := ph.app.Conf.MoviesMigrationPercent

	if percent > migrationPercent {
		return ph.app.Conf.MonolithUrl
	}

	return ph.target
}
