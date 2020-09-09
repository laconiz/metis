package connector

import (
	"errors"
	"github.com/laconiz/metis/log"
	"github.com/laconiz/metis/nets/event"
	"github.com/laconiz/metis/nets/session"
	"sync"
	"time"
)

func New(addr string, dialer Dialer, option Option, logger log.Logger) *Connector {

	const fieldName = "name"

	option.parse()

	return &Connector{
		addr:   addr,
		dialer: dialer,
		option: &option,
		logger: logger.Field(fieldName, option.Name),
	}
}

type Connector struct {
	addr      string           // 连接地址
	dialer    Dialer           // 连接器
	option    *Option          // 配置信息
	logger    log.Logger       // 日志接口
	reconnect bool             // 是否重连
	times     int              // 重连次数
	session   *session.Session // 会话
	mutex     sync.RWMutex
}

func (con *Connector) Connected() bool {

	con.mutex.RLock()
	defer con.mutex.RUnlock()

	return con.session != nil
}

func (con *Connector) Run() {

	con.mutex.Lock()
	defer con.mutex.Unlock()

	if con.session != nil {
		return
	}

	// 重置重连参数
	con.reconnect = con.option.Reconnect

	// 尝试连接
	con.connect()
}

func (con *Connector) connect() {

	con.logger.Data("addr", con.addr).Info("connect")

	conn, err := con.dialer(con.addr)
	ses := session.New(conn, &con.option.Session, con.logger)

	if err != nil {

		con.logger.Data("error", err).Error("dial error")

		// 延时重连
		go con.delay()
		// 连接失败事件
		go con.option.Session.Invoker.Invoke(event.NewConnectFailed())

		return
	}

	// 设置会话
	con.session = ses
	// 重置重连次数
	con.times = 0

	// 运行会话
	go con.session.Run(func(ses *session.Session) {

		// 会话终止
		con.mutex.Lock()
		defer con.mutex.Unlock()

		// 防止回调时已重新连接
		if con.session == ses {
			// 重置会话
			con.session = nil
			// 延时重连
			go con.delay()
		}
	})
}

// 延时重连
func (con *Connector) delay() {

	con.mutex.Lock()
	defer con.mutex.Unlock()

	// 无需重连 || 已经连接
	if !con.reconnect || con.session != nil {
		return
	}

	// 延时事件
	delays := con.option.Delays
	index := con.times
	if index >= len(delays) {
		index = len(delays) - 1
	}

	// 增加重连次数
	con.times++

	// 延时
	go func() {

		delay := delays[index]

		con.logger.Data("delay", delay.String()).Info("reconnect")

		// 等待
		<-time.After(delay)

		con.mutex.Lock()
		defer con.mutex.Unlock()

		// 无需重连 || 已经连接
		if !con.reconnect || con.session != nil {
			return
		}

		// 尝试连接
		con.connect()
	}()
}

func (con *Connector) Stop() {

	con.mutex.Lock()
	defer con.mutex.Unlock()

	// 不再重连
	con.reconnect = false

	if con.session == nil {
		return
	}

	// 关闭会话
	con.session.Close()
	con.session = nil
}

// 发送消息
func (con *Connector) Send(msg interface{}) error {

	con.mutex.RLock()
	defer con.mutex.RUnlock()

	if con.session == nil {
		return errDisconnected
	}

	return con.session.Send(msg)
}

// 发送消息流
func (con *Connector) SendRaw(raw []byte) error {

	con.mutex.RLock()
	defer con.mutex.RUnlock()

	if con.session == nil {
		return errDisconnected
	}

	return con.session.SendRaw(raw)
}

var errDisconnected = errors.New("disconnected")
