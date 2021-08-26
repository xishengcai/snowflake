package snowflake

import (
	"fmt"
	"lsh-mcp-lcr-cops/pkg/util/snowflake"
"sync"
	"testing"
	"time"
)

func TestSnowflakeConcurrence(t *testing.T) {
	plan := 10000
	wg := sync.WaitGroup{}
	s := snowflake.NewSnowflake(1)
	ch := make(chan uint64, plan)
	wg.Add(plan)
	defer close(ch)
	//并发 count个goroutine 进行 snowflake ID 生成
	for i := 0; i < plan; i++ {
		go func() {
			defer wg.Done()
			id, _ := s.NextID()
			ch <- id
		}()
	}
	wg.Wait()
	m := make(map[uint64]int)
	for i := 0; i < plan; i++ {
		id := <-ch
		// 如果 map 中存在为 id 的 key, 说明生成的 snowflake ID 有重复
		_, ok := m[id]
		if ok {
			fmt.Printf("repeat id %d\n", id)
			return
		}
		// 将 id 作为 key 存入 map
		m[id] = i
	}
	// 成功生成 snowflake ID
	t.Log("All", len(m), "snowflake ID Get Success!")
}

func TestSnowflake(t *testing.T) {
	s := snowflake.NewSnowflake(1)
	tick := time.NewTicker(time.Nanosecond * 1e4)

	i := 20
	for {
		if i < 0 {
			break
		}
		select {
		case <-tick.C:
			id, err := s.NextID()
			if err != nil {
				t.Fatal(err)
			}
			t.Log("id: ", id)
			t.Logf("%+v,    now:%d", s, time.Now().UnixNano()/1e6)
		}
		i--
	}
}
