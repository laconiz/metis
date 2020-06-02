package websocket

import "time"

// 响应路径
const Path = "/ws"

// 握手超时时间
const HandshakeTimeout = time.Second * 3
