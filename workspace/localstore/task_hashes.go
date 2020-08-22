package localstore

import (
	"encoding/binary"
	"github.com/unravela/delvin/api"
	"io"
	"os"
	"path/filepath"
)

func (ls *LocalStore) PutTaskHash(taskRef api.Ref, hash api.TaskHash) error {
	refHash := taskRef.GetHash()
	ls.hashes[refHash] = hash
	return ls.saveHashes()
}

func (ls *LocalStore) GetTaskHash(taskRef api.Ref) api.TaskHash {
	refHash := taskRef.GetHash()
	return ls.hashes[refHash]
}

func (ls *LocalStore) saveHashes() error {
	path := filepath.Join(ls.dir, "task_hashes")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	bytes := make([]byte, 16)
	for k, v := range ls.hashes {
		binary.LittleEndian.PutUint64(bytes[:8], k)
		binary.LittleEndian.PutUint64(bytes[8:], uint64(v))

		_, err := f.Write(bytes)
		if err != nil {
			return err
		}
	}

	return nil
}

func loadTaskHashes(dir string) (map[uint64]api.TaskHash, error) {
	path := filepath.Join(dir, "task_hashes")
	hashes := make(map[uint64]api.TaskHash)

	// If file doesn't exist, the empty storage is used
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return hashes, nil
	}

	// open the file
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// read the hashes
	refHash := make([]byte, 8)
	taskHash := make([]byte, 8)
	for {
		_, err := f.Read(refHash)
		if err == io.EOF {
			break
		}

		_, err = f.Read(taskHash)
		if err == io.EOF {
			break
		}

		hashes[binary.LittleEndian.Uint64(refHash)] = api.TaskHash(binary.LittleEndian.Uint64(taskHash))
	}

	return hashes, nil
}
