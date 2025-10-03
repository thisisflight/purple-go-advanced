package files

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type JSONDB struct {
	filename string
}

func NewJSONDB(filename string) *JSONDB {
	return &JSONDB{
		filename: filename,
	}
}

func (db *JSONDB) Read() ([]byte, error) {

	ext := strings.ToLower(filepath.Ext(db.filename))
	if ext != ".json" {
		return nil, fmt.Errorf("ожидалось расширение файла .json, получено: %s", ext)
	}

	data, err := os.ReadFile(db.filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (db *JSONDB) Write(content []byte) {
	file, err := os.Create(db.filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	_, err = file.Write(content)
	if err != nil {
		fmt.Println(err)
		return
	}
}
