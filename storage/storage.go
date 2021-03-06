package storage

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

const filename = "motd_storage.json"

var (
	storagePath = filepath.Join(os.TempDir(), filename)
	initialData = "quidquid Latine dictum sit altum videtur"
	mutex       = sync.RWMutex{}
)

// Datastore interface
type Datastore interface {
	InitData(path string) error
	Read() (*[]string, error)
	Add(message string) error
}

type Storage struct{}

// Messages contains a list of all stored messages of the day
type data struct {
	Messages []string `json:"messages"`
}

// InitData initializes storage
func (s Storage) InitData(path string) error {
	if len(path) > 0 {
		storagePath = path
	}
	dat, err := s.Read()
	if os.IsNotExist(err) || (err == nil && len(*dat) == 0) {
		f, _ := os.Create(storagePath)
		f.Close()
		return s.Add(initialData)
	}
	return err
}

func parseFile(content []byte) (*data, error) {
	m := new(data)
	err := json.Unmarshal(content, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Read returns current persisted messages
func (s Storage) Read() (*[]string, error) {
	mutex.RLock()
	defer mutex.RUnlock()
	return readFile()
}

func readFile() (*[]string, error) {
	dat, err := ioutil.ReadFile(storagePath)
	if err != nil {
		return nil, err
	}
	if len(dat) == 0 {
		return new([]string), nil
	}
	data, err := parseFile(dat)
	if err != nil {
		return nil, err
	}
	return &data.Messages, nil
}

func makeJSON(messages []string) ([]byte, error) {
	d := data{Messages: messages}
	j, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// Add persists a new message
func (s Storage) Add(message string) error {
	if len(message) == 0 {
		return errors.New("Empty message")
	}
	mutex.Lock()
	defer mutex.Unlock()
	dat, err := readFile()
	if err != nil {
		return err
	}
	*dat = append(*dat, message)
	b, err := makeJSON(*dat)
	if err == nil {
		err = ioutil.WriteFile(storagePath, b, 0)
	}
	return err
}
