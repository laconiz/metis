package session

import (
	cellutil "github.com/davyxu/cellnet/util"
	"github.com/laconiz/metis/log"
	"github.com/laconiz/metis/nets/cipher"
	"github.com/laconiz/metis/nets/encoder"
	"github.com/laconiz/metis/nets/event"
	"github.com/laconiz/metis/nets/queue"
	"sync"
	"sync/atomic"
	"time"
)

var unique uint64

func NewID() uint64 {
	return atomic.AddUint64(&unique, 1)
}

func New(conn Conn, option *Option, logger log.Logger) *Session {

	ses := &Session{
		id:      NewID(),
		conn:    conn,
		queue:   queue.New(option.Queue),
		cipher:  option.Cipher.New(),
		encoder: option.Encoder.New(),
		option:  option,
	}

	const fieldSession = "session"
	ses.logger = logger.Field(fieldSession, ses.id)

	return ses
}

type Session struct {
	id      uint64          // 会话ID
	conn    Conn            // 连接
	queue   *queue.Queue    // 写入队列
	cipher  cipher.Cipher   // 加密器
	encoder encoder.Encoder // 编码器
	logger  log.Logger      // 日志接口
	option  *Option         // 配置信息
	data    sync.Map        // 附加信息
}

func (ses *Session) ID() uint64 {
	return ses.id
}

func (ses *Session) Addr() string {
	return ses.conn.Addr()
}

func (ses *Session) Data() *sync.Map {
	return &ses.data
}

func (ses *Session) Close() error {
	return ses.queue.Close()
}

func (ses *Session) Send(msg interface{}) error {

	// 打包
	pkt, err := ses.encoder.Marshal(msg)
	if err != nil {
		ses.logger.Data("message", msg).Data("error", err).Error("marshal error")
		return err
	}

	return ses.queue.Push(pkt.Stream)
}

func (ses *Session) SendRaw(raw []byte) error {
	return ses.queue.Push(raw)
}

// 读取线程
func (ses *Session) read() {

	for {

		// 重置超时时间
		deadline := time.Now().Add(ses.option.Timeout)
		ses.conn.Deadline(deadline)

		// 读取流
		stream, err := ses.conn.Read()
		if err != nil {
			ses.logger.Data("error", err).Info("read error")
			break
		}

		// 解密
		raw, err := ses.cipher.Decode(stream)
		if err != nil {
			ses.logger.Data(string(stream), err).Warn("decode error")
			return
		}

		// 解包
		pkt, err := ses.encoder.Unmarshal(raw)
		if err != nil {
			ses.logger.Data(string(raw), err).Warn("unmarshal error")
			return
		}

		ses.logger.Data("raw", string(raw)).Debug("recv message")
		ses.Invoke(&event.Event{Ses: ses, Packet: pkt})
	}
}

// 写入线程
func (ses *Session) write() {

	for {

		// 重置超时时间
		deadline := time.Now().Add(ses.option.Timeout)
		ses.conn.Deadline(deadline)

		// 读取写入队列
		event, closed := ses.queue.Pop()
		if closed {
			break
		}

		raw := event.([]byte)

		// 加密
		stream, err := ses.cipher.Encode(raw)
		if err != nil {
			ses.logger.Data(string(raw), err).Error("encode error")
			continue
		}

		// 写入流
		if err := ses.conn.Write(stream); err != nil {
			ses.logger.Data(string(stream), err).Warn("write error")
			closed = true
			break
		}

		ses.logger.Data("raw", string(raw)).Debug("send message")
	}
}

// 启动会话
func (ses *Session) Run(callback func(*Session)) {

	ses.logger.Data("addr", ses.Addr()).Info("connected")

	// 运行写入线程
	go func() {
		ses.write()
		ses.conn.Close()
	}()

	// 连接事件
	ses.Invoke(event.NewConnected(ses))
	// 读取线程
	ses.read()
	// 关闭写入队列
	ses.queue.Close()

	// 关闭回调
	callback(ses)
	// 断开事件
	ses.Invoke(event.NewDisconnected(ses))

	ses.logger.Info("disconnected")
}

// 事件派发
func (ses *Session) Invoke(event *event.Event) {

	defer func() {
		if err := recover(); err != nil {
			ses.logger.Data("error", err).Errorf("invoke panic %s", cellutil.StackToString(5))
		}
	}()

	ses.option.Invoker.Invoke(event)
}
