package log

import (
	"github.com/laconiz/metis/log/context"
	"time"
)

// 日志信息
type Log struct {
	Level   Level            // 等级
	Time    time.Time        // 事件
	Data    *context.Data    // 数据
	Context *context.Context // 上下文
	Message string           // 格式化信息
}
