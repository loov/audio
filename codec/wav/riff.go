package wav

func findSubChunk(data []byte, target [4]byte) []byte {
	var id [4]byte
	var size uint32
	for len(data) >= 8 {
		copy(id[:], data)
		readU32LE(&size, data[4:])
		if id == target {
			return data
		}

		chunkSize := int(size)
		if id == [4]byte{'R', 'I', 'F', 'F'} {
			chunkSize = 4
		}

		data = data[8+chunkSize:]
	}
	return nil
}

type header struct {
	ChunkID   [4]byte // "RIFF"
	ChunkSize uint32  // 4 + (8 + fmt.ChunkSize) + (8 + data.ChunkSize)
	Format    [4]byte // "WAVE"
}

type format struct {
	ChunkID       [4]byte // "fmt "
	ChunkSize     uint32  // 16 for PCM, size rest of header
	AudioFormat   uint16  // 1 == linear, ...
	NumChannels   uint16  // 1, 2, ...
	SampleRate    uint32  // 8000, 41000 ...
	ByteRate      uint32  // SampleRate * NumChannels * BitsPerSample / 8
	BlockAlign    uint16  // NumChannels * BitsPerSample / 8
	BitsPerSample uint16  // 8, 16
	// Extra Parameters
}

type data struct {
	ChunkID   [4]byte // "data"
	ChunkSize uint32  // NumSamples * NumChannels * BitsPerSample / 8. bytes in data
	Data      []byte
}

func readU32LE(r *uint32, v []byte) int {
	*r = uint32(v[0])<<0 | uint32(v[1])<<8 | uint32(v[2])<<16 | uint32(v[3])<<24
	return 4
}

func readU16LE(r *uint16, v []byte) int {
	*r = uint16(v[0]<<0) | uint16(v[1]<<8)
	return 2
}

func readU8LE(r *uint8, v []byte) int {
	*r = uint8(v[0])
	return 1
}

func (chunk *header) Read(data []byte) (rest []byte) {
	p := 0
	p += copy(chunk.ChunkID[:], data[p:])
	p += readU32LE(&chunk.ChunkSize, data[p:])
	p += copy(chunk.Format[:], data[p:])
	return data[p:]
}

func (chunk *format) Read(data []byte) (rest []byte) {
	data = findSubChunk(data, [4]byte{'f', 'm', 't', ' '})
	p := 0
	p += copy(chunk.ChunkID[:], data[p:])
	p += readU32LE(&chunk.ChunkSize, data[p:])
	p += readU16LE(&chunk.AudioFormat, data[p:])
	p += readU16LE(&chunk.NumChannels, data[p:])
	p += readU32LE(&chunk.SampleRate, data[p:])
	p += readU32LE(&chunk.ByteRate, data[p:])
	p += readU16LE(&chunk.BlockAlign, data[p:])
	p += readU16LE(&chunk.BitsPerSample, data[p:])
	return data[8+int(chunk.ChunkSize):]
}

func (chunk *data) Read(data []byte) (rest []byte) {
	data = findSubChunk(data, [4]byte{'d', 'a', 't', 'a'})

	p := 0
	p += copy(chunk.ChunkID[:], data[p:])
	p += readU32LE(&chunk.ChunkSize, data[p:])

	k := p + int(chunk.ChunkSize)
	if k > len(data) {
		k = len(data)
	}
	chunk.Data = data[p:k]
	return data[k:]
}
