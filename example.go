package snowflake

import (
	"lsh-mcp-lcr-cops/pkg/util/snowflake"
"sync"
)

var instance *snowflake.Snowflake
var once sync.Once

func Instance() *snowflake.Snowflake {
	once.Do(func() {
		// TODO: 目前不支持节点横向扩展
		instance = snowflake.NewSnowflake(1)
	})
	return instance
}
