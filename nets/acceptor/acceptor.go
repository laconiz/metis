package acceptor

import (
	"github.com/laconiz/metis/log"
	"github.com/laconiz/metis/nets/session"
	"strings"
	"sync"
)

const StrClosed = "use of closed network connection"

type Listener interface {
	Accept() (session.Conn, error)
	Close() error
}

type Dialer func(addr string) (Listener, error)

// 创建侦听器
func New(addr string, dialer Dialer, option Option, logger log.Logger) *Acceptor {

	const fieldName = "acceptor"

	option.parse()

	return &Acceptor{
		addr:     addr,
		dialer:   dialer,
		option:   &option,
		logger:   logger.Field(fieldName, option.Name),
		sessions: Sessions{},
	}
}

// 会话列表
type Sessions map[uint64]*session.Session

type Acceptor struct {
	addr     string     // 侦听地址
	dialer   Dialer     // 侦听器生成函数
	option   *Option    // 配置信息
	logger   log.Logger // 日志接口
	listener Listener   // 侦听器
	sessions Sessions   // 会话列表
	mutex    sync.RWMutex
}

// 是否运行
func (acc *Acceptor) Running() bool {

	acc.mutex.RLock()
	defer acc.mutex.RUnlock()

	return acc.listener != nil
}

// 当前会话数量
func (acc *Acceptor) Count() int64 {

	acc.mutex.RLock()
	defer acc.mutex.RUnlock()

	return int64(len(acc.sessions))
}

// 启动
func (acc *Acceptor) Run() {

	acc.mutex.Lock()
	defer acc.mutex.Unlock()

	// 已启动
	if acc.listener != nil {
		return
	}

	acc.logger.Data(acc.addr).Info("accept")

	// 启动侦听器
	listener, err := acc.dialer(acc.addr)
	if err != nil {
		acc.logger.Data(err).Error("dial error")
		return
	}
	acc.listener = listener

	go func() {

		for {

			conn, err := listener.Accept()
			if err != nil {
				// 非正常关闭
				if !strings.Contains(err.Error(), StrClosed) {
					acc.logger.Data(err).Error("accept error")
				}
				break
			}

			ses := session.New(conn, &acc.option.Session, acc.logger)

			// 添加会话
			acc.mutex.Lock()
			acc.sessions[ses.ID()] = ses
			acc.mutex.Unlock()

			// 运行会话
			go ses.Run(func(ses *session.Session) {
				// 移除会话
				acc.mutex.Lock()
				delete(acc.sessions, ses.ID())
				acc.mutex.Unlock()
			})
		}

		acc.mutex.Lock()
		defer acc.mutex.Unlock()

		if acc.listener == listener {
			acc.listener = nil
		}

		acc.logger.Info("stopped")
	}()
}

// 关闭
func (acc *Acceptor) Stop() {

	acc.mutex.Lock()
	defer acc.mutex.Unlock()

	// 已关闭
	if acc.listener == nil {
		return
	}

	// 关闭所有会话
	for _, ses := range acc.sessions {
		ses.Close()
	}

	// 关闭侦听器
	acc.listener.Close()
	acc.listener = nil
}

// 广播消息
func (acc *Acceptor) Broadcast(msg interface{}) {
	for _, ses := range acc.sessions {
		ses.Send(msg)
	}
}

// 广播消息流
func (acc *Acceptor) BroadcastRaw(raw []byte) {
	for _, ses := range acc.sessions {
		ses.SendRaw(raw)
	}
}
