package sharder

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/rand"
	"sort"
)

type ShardResult struct {
	Bytes         []byte
	OffsetIndices []int // start index of each content offset
	Offsets       []int // raw offsets in add order
}

type Chunk struct {
	Order          int
	Offset         int
	Limit          int
	OffsetInBuffer int
}

type ShardContent struct {
	Data                  []byte
	Chunks                []*Chunk
	AvailableChunkIndices []int
	InsertedChunks        []*Chunk
}

type Sharder struct {
	Contents                []*ShardContent
	AvailableContentIndices []int
	TotalSize               int
	NumChunks               int
	Threshold               float64
	Seed                    int64
	Rand                    *rand.Rand
}

func NewSharder(seed int64) *Sharder {
	return &Sharder{
		Contents:                make([]*ShardContent, 0),
		AvailableContentIndices: make([]int, 0),
		Threshold:               0.25,
		Rand:                    rand.New(rand.NewSource(seed)),
	}
}

func (s *Sharder) AddRandomData(size int32) {
	if size > 0 {
		data := make([]byte, size)
		s.Rand.Read(data)
		s.Add(data)
	}
}

func (s *Sharder) Add(value []byte) (ptr uint64) {
	if s.Rand.Int()%2 == 0 {
		s.AddRandomData(s.Rand.Int31n(int32(len(value)/2) + 1))
	}

	s.AvailableContentIndices = append(s.AvailableContentIndices, len(s.Contents))

	content := &ShardContent{
		Data:           value,
		InsertedChunks: make([]*Chunk, 0),
	}

	s.splitData(content)

	index := 0
	for _, c := range s.Contents {
		index += len(c.Chunks)
	}
	ptr = (uint64(len(content.Chunks)) << 32) | uint64(index)
	// fmt.Printf("[ADD DATA] [Chunks offset=%d] [Number of chunks=%d]\n", index, len(content.Chunks))

	s.Contents = append(s.Contents, content)

	return
}

func (s *Sharder) splitData(c *ShardContent) {
	dataSize := len(c.Data)
	s.TotalSize += dataSize

	chunkSize := 0
	if s.Threshold > 0 && s.Threshold < 1 {
		chunkSize = int(float64(dataSize) * s.Threshold)
	} else if s.Threshold >= 1 {
		chunkSize = int(math.Ceil(s.Threshold))
	}

	if chunkSize == 0 {
		chunkSize = dataSize
	}

	numChunks := dataSize / chunkSize
	remainder := dataSize % chunkSize
	if remainder != 0 {
		numChunks++
	}

	s.NumChunks += numChunks
	c.Chunks = make([]*Chunk, numChunks)
	c.AvailableChunkIndices = make([]int, numChunks)

	for i := 0; i < numChunks; i++ {
		c.AvailableChunkIndices[i] = i
		c.Chunks[i] = &Chunk{
			Order:  i,
			Offset: i * chunkSize,
			Limit:  chunkSize,
		}

		if i+1 == numChunks && remainder != 0 {
			c.Chunks[i].Limit = remainder
		}
	}
}

func (s *Sharder) getRandomChunk() (content *ShardContent, chunk *Chunk) {
	aci := s.Rand.Intn(len(s.AvailableContentIndices))
	contentIndex := s.AvailableContentIndices[aci]
	content = s.Contents[contentIndex]

	achi := s.Rand.Intn(len(content.AvailableChunkIndices))
	chunkIndex := content.AvailableChunkIndices[achi]
	chunk = content.Chunks[chunkIndex]

	content.AvailableChunkIndices = removeIntSlice(content.AvailableChunkIndices, achi)
	if len(content.AvailableChunkIndices) == 0 {
		s.AvailableContentIndices = removeIntSlice(s.AvailableContentIndices, aci)
	}

	return
}

func (s *Sharder) Shard() *ShardResult {
	result := &ShardResult{
		OffsetIndices: make([]int, 0),
		Offsets:       make([]int, 0),
	}

	buffer := bytes.Buffer{}
	for i := 0; i < s.NumChunks; i++ {
		content, chunk := s.getRandomChunk()
		chunk.OffsetInBuffer = buffer.Len()
		content.InsertedChunks = append(content.InsertedChunks, chunk)
		buffer.Write(content.Data[chunk.Offset : chunk.Offset+chunk.Limit])

		// fmt.Printf("[SHARD LOOP] [Index=%d] [Selected content=%s] [Selected chunk order=%d] [Selected chunk data=%s] [Chunk offset=%d]\n",
		// 	i, string(content.Data), chunk.Order, string(content.Data[chunk.Offset:chunk.Offset+chunk.Limit]), chunk.OffsetInBuffer)
	}

	for _, content := range s.Contents {
		result.OffsetIndices = append(result.OffsetIndices, len(result.Offsets))
		sort.Slice(content.InsertedChunks, func(i, j int) bool {
			chi, chj := content.InsertedChunks[i], content.InsertedChunks[j]
			return chi.Order < chj.Order
		})
		for _, chunk := range content.InsertedChunks {
			result.Offsets = append(result.Offsets, chunk.OffsetInBuffer)
		}
	}

	result.Bytes = buffer.Bytes()

	return result
}

func (s *ShardResult) RawData() []byte { // num offsets, offsets (each 4 bytes), bytes
	buffer := bytes.Buffer{}

	lenBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(lenBytes, uint32(len(s.Offsets)))
	buffer.Write(lenBytes)

	for _, offset := range s.Offsets {
		offsetBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(offsetBytes, uint32(offset))
		buffer.Write(offsetBytes)
	}

	buffer.Write(s.Bytes)

	return buffer.Bytes()
}
