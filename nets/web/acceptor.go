package web

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/laconiz/metis/log"
	"net/http"
	"sync"
	"time"
)

func NewAcceptor(option *AcceptorOption, logger log.Logger) (*Acceptor, error) {

	option.parse()

	logger = logger.Level(option.Level).Field("acceptor", option.Name)

	engine := gin.New()
	engine.Use(gin.Recovery())

	invoker := NewInvoker(logger).Arguments(option.Params...).Builders(option.Creators...)
	if err := invoker.Register(engine, option.Nodes); err != nil {
		return nil, err
	}

	return &Acceptor{listener: &http.Server{Addr: option.Addr, Handler: engine}, logger: logger}, nil
}

type Acceptor struct {
	running  bool
	listener *http.Server // 侦听器
	logger   log.Logger   // 日志接口
	mutex    sync.RWMutex
}

func (acceptor *Acceptor) Running() bool {
	acceptor.mutex.RLock()
	defer acceptor.mutex.RUnlock()
	return acceptor.running
}

func (acceptor *Acceptor) Engine() *gin.Engine {
	acceptor.mutex.RLock()
	defer acceptor.mutex.RUnlock()
	return acceptor.listener.Handler.(*gin.Engine)
}

func (acceptor *Acceptor) Run() {

	acceptor.mutex.Lock()
	defer acceptor.mutex.Unlock()
	if acceptor.running {
		return
	}
	acceptor.running = true

	acceptor.listener = &http.Server{Addr: acceptor.listener.Addr, Handler: acceptor.listener.Handler}
	acceptor.logger.Data("addr", acceptor.listener.Addr).Info("start")

	go func() {

		err := acceptor.listener.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			acceptor.logger.Data("error", err).Error("listen error")
		}

		acceptor.mutex.Lock()
		defer acceptor.mutex.Unlock()

		acceptor.running = false
		acceptor.logger.Info("stopped")
	}()
}

func (acceptor *Acceptor) Stop() {

	acceptor.mutex.Lock()
	defer acceptor.mutex.Unlock()
	if !acceptor.running {
		return
	}
	acceptor.running = false

	acceptor.logger.Info("shutting down")

	context, _ := context.WithTimeout(context.Background(), time.Second*2)
	if err := acceptor.listener.Shutdown(context); err != nil {
		acceptor.logger.Data("error", err).Error("shutdown error")
	}
}

// ---------------------------------------------------------------------------------------------------------------------

func init() {
	gin.SetMode(gin.ReleaseMode)
}
