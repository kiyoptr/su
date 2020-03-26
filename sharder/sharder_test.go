package sharder

import (
	"crypto/rand"
	"testing"
	"time"
)

func TestShard(t *testing.T) {
	for i := 0; i < 1000; i++ {
		s := NewSharder(time.Now().UnixNano())

		s.Threshold = s.Rand.Float64()
		if s.Threshold < 0.1 {
			s.Threshold *= 10
		}

		numData := s.Rand.Int31n(50)
		values := make(map[string]uint64)
		for i := int32(0); i < numData; i++ {
			dataLen := s.Rand.Int31n(1024*10) + 1
			data := make([]byte, dataLen)
			rand.Read(data)
			values[string(data)] = 0
		}

		for k, _ := range values {
			values[k] = s.Add([]byte(k))
		}

		result := s.Shard()
		raw := result.RawData()

		for k, v := range values {
			data := GetData(raw, v)
			if string(data) != k {
				t.Fatalf("Invalid data")
			}
		}
	}
}
