package wav

import (
	"encoding/hex"
	"reflect"
	"strings"
	"testing"
)

func hex2bytes(data string) []byte {
	data = strings.Map(func(r rune) rune {
		if '0' <= r && r <= '9' {
			return r
		}
		if 'a' <= r && r <= 'f' {
			return r
		}
		if 'A' <= r && r <= 'F' {
			return r
		}
		return -1
	}, data)

	x, err := hex.DecodeString(data)
	if err != nil {
		panic(err)
	}
	return append(x, make([]byte, 2048)...)
}

var ExampleDataShort = hex2bytes(`
	52 49 46 46 24 08 00 00
		57 41 56 45
	66 6d 74 20 10 00 00 00
		01 00 02 00 22 56 00 00 88 58 01 00 04 00 10 00
	64 61 74 61 00 08 00 00
		00 00 00 00	24 17 1e f3 3c 13 3c 14 16 f9 18 f9 34 e7 23 a6 3c f2 24 f2 11 ce 1a 0d 
`)

func TestBytesReader(t *testing.T) {
	rd, err := NewBytesReader(ExampleDataShort)
	if err != nil {
		t.Errorf("failed to read: %v", err)
		return
	}

	hdr := header{
		ChunkID:   [4]byte{'R', 'I', 'F', 'F'},
		ChunkSize: 2084,
		Format:    [4]byte{'W', 'A', 'V', 'E'},
	}

	fmt := format{
		ChunkID:       [4]byte{'f', 'm', 't', ' '},
		ChunkSize:     16,
		AudioFormat:   1,
		NumChannels:   2,
		SampleRate:    22050,
		ByteRate:      88200,
		BlockAlign:    4,
		BitsPerSample: 16,
	}

	if !reflect.DeepEqual(rd.header, hdr) {
		t.Errorf("invalid header got:\n%#+v\nexpected:\n%#+v", rd.header, hdr)
	}
	if !reflect.DeepEqual(rd.format, fmt) {
		t.Errorf("invalid format got:\n%#+v\nexpected:\n%#+v", rd.format, fmt)
	}

	if rd.data.ChunkID != [4]byte{'d', 'a', 't', 'a'} {
		t.Errorf("invalid data.ChunkID got:\n%#+v", rd.data)
	}
	if rd.data.ChunkSize != 2048 {
		t.Errorf("invalid data.ChunkSize got:\n%#+v", rd.data)
	}
}
