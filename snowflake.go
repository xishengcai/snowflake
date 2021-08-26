package snowflake

import (
	"fmt"
	"sync"
	"time"
)

// timestamp 无位数限制， workerID 限制位1， maxSequence 1毫秒内序列号生成最大值
// 由于
const (
	workerIDBits = 1 // 10bit 工作机器ID中的 5bit workerID
	sequenceBits = 6

	maxWorkerID = int64(1)<<workerIDBits - 1  // 节点ID的最大值 用于防止溢出
	maxSequence = int64(-1)<<sequenceBits - 1 // ID序列号的最大值 用于防止溢出

	timeLeft = workerIDBits + sequenceBits // 时间戳向左偏移量
	workLeft = sequenceBits                // 节点IDx向左偏移量

	// startTime 41bit位的常量时间戳，单位毫秒。
	// 2020-05-20 05:20:00 +0800 CST
	startTime = int64(1589923200000)
)

type Snowflake struct {
	mu        sync.Mutex
	LastStamp int64 // 记录上一次ID的时间戳
	WorkerID  int64 // 该节点的ID
	Sequence  int64 // 当前毫秒已经生成的ID序列号(从0 开始累加) 1毫秒内最多生成4096个ID
}

// NewSnowflake 分布式情况下,我们应通过外部配置文件或其他方式为每台机器分配独立的id
func NewSnowflake(workerID int64) *Snowflake {
	if workerID > maxWorkerID {
		panic("workID 超出范围 ")
	}
	return &Snowflake{
		WorkerID:  workerID,
		LastStamp: 0,
		Sequence:  0,
	}
}

func (s *Snowflake) getMilliSeconds() int64 {
	return time.Now().UnixNano() / 1e6 // 毫秒
}

func (s *Snowflake) NextID() (uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.nextID()
}

func (s *Snowflake) nextID() (uint64, error) {
	timeStamp := s.getMilliSeconds()
	if timeStamp < s.LastStamp {
		return 0, fmt.Errorf("machine time is not accurate, time now is: %v", time.Now().UTC())
	}
	if s.LastStamp == timeStamp {
		s.Sequence = (s.Sequence + 1) & maxSequence
		if s.Sequence == 0 {
			for timeStamp <= s.LastStamp {
				timeStamp = s.getMilliSeconds()
			}
		}
	} else {
		s.Sequence = 0
	}
	s.LastStamp = timeStamp
	id := ((timeStamp - startTime) << timeLeft) |
		(s.WorkerID << workLeft) |
		s.Sequence
	return uint64(id), nil
}
