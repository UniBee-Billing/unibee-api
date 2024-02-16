package utility

import (
	"sync"
	"time"
)

// Snowflake 结构定义
type Snowflake struct {
	mu        sync.Mutex
	timestamp int64 // 时间戳
	workerID  int64 // 工作机器ID
	sequence  int64 // 序列号
}

// NewSnowflake 创建一个新的Snowflake实例
func NewSnowflake(workerID int64) *Snowflake {
	return &Snowflake{
		timestamp: 0,
		workerID:  workerID,
		sequence:  0,
	}
}

// GenerateID 生成唯一ID
func (s *Snowflake) GenerateID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	currentTime := time.Now().UnixNano() / 1e6

	if s.timestamp == currentTime {
		s.sequence = (s.sequence + 1) & 4095
		if s.sequence == 0 {
			for currentTime <= s.timestamp {
				currentTime = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		s.sequence = 0
	}

	s.timestamp = currentTime

	id := (currentTime-epoch)<<22 | (s.workerID << 12) | s.sequence
	return id
}

const (
	epoch = 1597536000000 // 起始时间戳，这里设定为2020-08-16 00:00:00 UTC的毫秒数
)
