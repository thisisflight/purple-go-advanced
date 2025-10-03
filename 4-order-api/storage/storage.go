package storage

import (
	"encoding/json"
	"fmt"
)

type ByteReader interface {
	Read() ([]byte, error)
}

type ByteWriter interface {
	Write(content []byte)
}

type DB interface {
	ByteReader
	ByteWriter
}

type TokenRecord struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type TokenRepository struct {
	Data []TokenRecord `json:"data"`
	db   DB
}

func (tr *TokenRepository) ToBytes() ([]byte, error) {
	return json.MarshalIndent(tr, "", "\t")
}

func (tr *TokenRepository) save() {
	data, err := tr.ToBytes()
	if err != nil {
		fmt.Println("Не удалось преобразовать")
		return
	}
	tr.db.Write(data)
}

func (tr *TokenRepository) AddTokenRecord(record *TokenRecord) {
	tr.Data = append(tr.Data, *record)
	tr.save()
}

func (tr *TokenRepository) RemoveRecordByToken(token string) bool {
	for i, rec := range tr.Data {
		if rec.Token == token {
			tr.Data = append(tr.Data[:i], tr.Data[i+1:]...)
			tr.save()
			return true
		}
	}
	return false
}

func NewTokenRepository(db DB) *TokenRepository {
	file, err := db.Read()
	if err != nil {
		return &TokenRepository{
			Data: []TokenRecord{},
			db:   db,
		}
	}

	var verifyData TokenRepository
	err = json.Unmarshal(file, &verifyData)
	if err != nil {
		return &TokenRepository{
			Data: []TokenRecord{},
			db:   db,
		}
	}
	verifyData.db = db
	return &verifyData
}
