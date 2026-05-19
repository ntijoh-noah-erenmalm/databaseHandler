package storage

import (
    "os"
    "encoding/json"
)

const dataPath = "data.db"

type Storage struct {
    file       *os.File
    nextOffset int64
}

func NewStorage() (*Storage, error) {
	file, error := os.OpenFile(dataPath, os.O_RDWR|os.O_CREATE, 0666)
	if error != nil {
		return nil, error
	}
	info, error := file.Stat()
	if error != nil {
		return nil, error
	}
	return &Storage{file:file, nextOffset: info.Size()}, nil
}


func (s *Storage) Write(record map[string]string) (int64, int, error) {
	data, error := json.Marshal(record)
	if error != nil {
		return 0,0,nil
	}
	offset := s.nextOffset
	size := len(data)
	_, error = s.file.WriteAt(data, offset)
	if error != nil {
		return 0,0,error
	}
	s.nextOffset += int64(size)
	return offset, size, nil
}

func (s *Storage) Read(offset int64, size int) (map[string]string, error) {
	buffer := make([]byte, size)
	_, error := s.file.ReadAt(buffer, offset)
	if error != nil {
		return nil, error
	}
	var record map[string]string
	error = json.Unmarshal(buffer, &record)
	return record, error
}

func (s *Storage) Close() {
	s.file.Close()
}






