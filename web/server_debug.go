/*
 *  Copyright (c) 2022-2024 Mikhail Knyazhev <markus621@yandex.ru>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package web

import (
	"net/http"
	"net/http/pprof"

	"go.osspkg.com/goppy/xc"
	"go.osspkg.com/goppy/xlog"
)

// ServerDebug service model
type ServerDebug struct {
	server *ServerHttp
	route  *BaseRouter
}

// NewServerDebug init debug service
func NewServerDebug(c ConfigHttp, l xlog.Logger) *ServerDebug {
	route := NewBaseRouter()
	return &ServerDebug{
		server: NewServerHttp(c, route, l),
		route:  route,
	}
}

// Up start service
func (o *ServerDebug) Up(ctx xc.Context) error {
	o.route.Route("/debug/pprof", pprof.Index, http.MethodGet)
	o.route.Route("/debug/pprof/goroutine", pprof.Index, http.MethodGet)
	o.route.Route("/debug/pprof/allocs", pprof.Index, http.MethodGet)
	o.route.Route("/debug/pprof/block", pprof.Index, http.MethodGet)
	o.route.Route("/debug/pprof/heap", pprof.Index, http.MethodGet)
	o.route.Route("/debug/pprof/mutex", pprof.Index, http.MethodGet)
	o.route.Route("/debug/pprof/threadcreate", pprof.Index, http.MethodGet)
	o.route.Route("/debug/pprof/cmdline", pprof.Cmdline, http.MethodGet)
	o.route.Route("/debug/pprof/profile", pprof.Profile, http.MethodGet)
	o.route.Route("/debug/pprof/symbol", pprof.Symbol, http.MethodGet)
	o.route.Route("/debug/pprof/trace", pprof.Trace, http.MethodGet)
	return o.server.Up(ctx)
}

// Down stop service
func (o *ServerDebug) Down() error {
	return o.server.Down()
}
