package rproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"

	"github.com/Juanmagc99/portalis/internal/balancer"
	"github.com/Juanmagc99/portalis/internal/registry"
	"github.com/labstack/echo/v4"
)

// Proxy maneja el balanceo round-robin dinámico y proxy inverso.
type Proxy struct {
	store     registry.Registry
	balancers map[string]*balancer.RoundRobin
	mu        sync.Mutex
}

// NewReverseProxyMiddleware crea el middleware de proxy con balanceadores inicializados.
func NewReverseProxyMiddleware(store registry.Registry) echo.MiddlewareFunc {
	p := &Proxy{
		store:     store,
		balancers: make(map[string]*balancer.RoundRobin),
	}
	return p.middleware()
}

func (p *Proxy) middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path

			if strings.HasPrefix(path, "/api/") {
				return next(c)
			}

			segments := strings.Split(strings.TrimPrefix(path, "/"), "/")
			if len(segments) == 0 || segments[0] == "" {
				return next(c)
			}
			svc := segments[0]

			p.mu.Lock()
			rr, ok := p.balancers[svc]
			if !ok {
				rr = &balancer.RoundRobin{}
				p.balancers[svc] = rr
			}
			p.mu.Unlock()

			insts, _ := p.store.List(svc)
			target := rr.Next(insts)
			if target == "" {
				return echo.NewHTTPError(http.StatusServiceUnavailable, "No instances available for "+svc)
			}

			u, err := url.Parse(target)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadGateway, "Invalid target URL")
			}

			callerSvc := c.Request().Header.Get("X-Service-Name")
			originalHost := c.Request().Host

			proxy := httputil.NewSingleHostReverseProxy(u)
			origDir := proxy.Director
			proxy.Director = func(req *http.Request) {
				origDir(req)
				req.Header.Set("X-Forwarded-Host", originalHost)
				if callerSvc != "" {
					req.Header.Set("X-Caller-Service", callerSvc)
				}
			}

			proxy.ServeHTTP(c.Response(), c.Request())
			return nil
		}
	}
}
