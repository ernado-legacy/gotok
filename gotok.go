package gotok

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
)

type Storage interface {
	Get(hexToken string) (*Token, error)
	Generate(id bson.ObjectId) (*Token, error)
	Remove(token *Token) error
}

func Generate(id bson.ObjectId) *Token {
	var bytes = make([]byte, 100)
	var hash = sha256.New()
	rand.Read(bytes)
	hash.Write(bytes)
	hash.Write([]byte(id))
	return &Token{id, hex.EncodeToString(hash.Sum(nil))}
}

type Token struct {
	Id    bson.ObjectId `json:"id"     bson:"user,omitempty"`
	Token string        `json:"token"  bson:"_id"`
}

func (t *Token) GetCookie() *http.Cookie {
	return &http.Cookie{Name: "token", Value: t.Token, Path: "/"}
}

type StorageMemory struct {
	tokens *mgo.Collection
	cache  map[string]bson.ObjectId
}

func (storage *StorageMemory) Get(hexToken string) (*Token, error) {
	t := &Token{}
	value, ok := storage.cache[hexToken]
	if ok {
		t.Id = value
		t.Token = hexToken
		return t, nil
	}
	err := storage.tokens.Find(bson.M{"_id": hexToken}).One(t)
	if err == mgo.ErrNotFound {
		return nil, nil
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}

	storage.cache[hexToken] = t.Id
	return t, nil
}

func (storage *StorageMemory) Generate(id bson.ObjectId) (*Token, error) {
	t := Generate(id)
	err := storage.tokens.Insert(t)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	storage.cache[t.Token] = t.Id
	return t, nil
}

func (storage *StorageMemory) Remove(token *Token) error {
	_, ok := storage.cache[token.Token]
	if ok {
		delete(storage.cache, token.Token)
	}
	return storage.tokens.Remove(bson.M{"_id": token.Token})
}

func New(collection *mgo.Collection) *StorageMemory {
	return &StorageMemory{collection, make(map[string]bson.ObjectId)}
}
