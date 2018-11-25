package main

import (
	"bytes"
	"crypto/sha256"
	"log"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	log.Printf("headers ---> %+v \n", headers)
	hash := sha256.Sum256(headers)
	log.Printf("hash ---> %v \n", hash)
	b.Hash = hash[:]
	log.Printf("b.Hash ---> %x \n", b.Hash)

}

func NewBlock(data string, prevBlockHash []byte) *Block {
	log.Printf("data to byte ---> %+v \n", []byte(data))
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}}
	block.SetHash()
	return block
}
