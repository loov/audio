package sequencer

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"strings"

	"github.com/loov/audio/codec/wav"
)

type Pool struct {
	Samples map[string]*wav.Reader
}

func NewPool() *Pool {
	return &Pool{
		Samples: make(map[string]*wav.Reader),
	}
}

func (pool *Pool) Subset(path string) []*wav.Reader {
	readers := []*wav.Reader{}
	for name, reader := range pool.Samples {
		if strings.HasPrefix(name, path) {
			readers = append(readers, reader)
		}
	}
	return readers
}

func (pool *Pool) LoadDirectory(root string) {
	//TODO: return errors
	filepath.Walk(root, func(name string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(name) != ".wav" {
			return nil
		}

		data, err := ioutil.ReadFile(name)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		reader, err := wav.NewBytesReader(data)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		name = strings.TrimPrefix(name, root)
		pool.Samples[filepath.ToSlash(name)] = reader
		return nil
	})
}
