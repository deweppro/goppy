package web

import (
	"sync"

	"github.com/deweppro/go-sdk/errors"
	"github.com/deweppro/go-sdk/log"
	"github.com/deweppro/goppy/plugins"
)

func WithWebsocketServerPool(options ...WebsocketServerOption) plugins.Plugin {
	return plugins.Plugin{
		Inject: func(l log.Logger) (*wssPool, WebsocketServerPool) {
			wssp := &wssPool{
				options: options,
				pool:    make(map[string]*wssProvider, 10),
				log:     l,
			}
			return wssp, wssp
		},
	}
}

type (
	wssPool struct {
		options []WebsocketServerOption
		pool    map[string]*wssProvider
		log     log.Logger
		mux     sync.Mutex
	}

	WebsocketServerPool interface {
		Create(name string) WebsocketServer
	}
)

func (v *wssPool) Create(name string) WebsocketServer {
	v.mux.Lock()
	defer v.mux.Unlock()

	if p, ok := v.pool[name]; ok {
		return p
	}

	u := newWebsocketUpgrader()
	for _, option := range v.options {
		option(u)
	}
	p := newWsServerProvider(v.log, u)
	v.pool[name] = p

	if err := p.Up(); err != nil {
		v.log.WithFields(log.Fields{
			"err":  err,
			"name": name,
		}).Errorf("Create Websocket Server in pool")
	}

	return p
}

func (v *wssPool) Up() error {
	return nil
}

func (v *wssPool) Down() error {
	v.mux.Lock()
	defer v.mux.Unlock()

	var err error
	for _, item := range v.pool {
		if e := item.Down(); e != nil {
			err = errors.Wrap(err, e)
		}
	}

	return err
}
