package sharder

import (
	"bytes"
	"encoding/binary"
	"math"
	"sort"
)

func GetData(rawShardedData []byte, ptr uint64) []byte {
	buffer := bytes.NewBuffer(rawShardedData)

	lenBytes := make([]byte, 4)
	buffer.Read(lenBytes)
	numOffsets := int(binary.LittleEndian.Uint32(lenBytes))

	offsets := make([]int, numOffsets)
	for i := 0; i < numOffsets; i++ {
		offsetBytes := make([]byte, 4)
		buffer.Read(offsetBytes)
		offsets[i] = int(binary.LittleEndian.Uint32(offsetBytes))
	}

	dataLength := buffer.Len()
	data := make([]byte, dataLength)
	buffer.Read(data)

	offsetIndex := int(ptr & math.MaxUint32)
	numChunks := int(ptr >> 32)

	sortedOffsets := make([]int, numOffsets)
	copy(sortedOffsets, offsets)
	sort.Slice(sortedOffsets, func(i, j int) bool { return sortedOffsets[i] < sortedOffsets[j] })

	buffer = &bytes.Buffer{}
	for i := 0; i < numChunks; i++ {
		chunkStart := offsets[offsetIndex+i]
		sortedIndex := findOffsetIndex(sortedOffsets, chunkStart)

		chunkEnd := dataLength
		if sortedIndex+1 < len(sortedOffsets) {
			chunkEnd = sortedOffsets[sortedIndex+1]
		}

		chunk := data[chunkStart:chunkEnd]
		buffer.Write(chunk)

		// fmt.Printf("[UNSHARD LOOP] [Chunk start index=%d] [Chunk end index=%d] [Index in sorted offsets=%d] [Chunk data=%s]\n",
		// 	chunkStart, chunkEnd, sortedIndex, string(chunk))
	}

	// fmt.Printf("[UNSHARD] [All data len=%d] [All chunks=%d] [Chunks offset=%d] [Num chunks=%d] [Retrieved bytes len=%d]\n",
	// 	dataLength, numOffsets, offsetIndex, numChunks, buffer.Len())

	return buffer.Bytes()
}

func findOffsetIndex(offsets []int, offset int) int {
	for i, index := range offsets {
		if index == offset {
			return i
		}
	}
	return -1
}
