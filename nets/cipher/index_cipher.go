package cipher

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"
)

type IndexCipher struct {
	sender   uint32
	receiver uint32
	rand     *rand.Rand
}

func (cipher *IndexCipher) Encode(raw []byte) ([]byte, error) {

	var buf bytes.Buffer

	rand := cipher.rand.Uint32()
	binary.Write(&buf, binary.LittleEndian, rand)
	cipher.sender++
	binary.Write(&buf, binary.LittleEndian, cipher.sender)
	buf.Write(raw)

	stream := buf.Bytes()
	mixer := mixBytes(stream)
	for i := randSize; i < len(stream); i++ {
		stream[i] ^= mixer[i%randSize]
	}

	return stream, nil
}

func (cipher *IndexCipher) Decode(stream []byte) ([]byte, error) {

	const header = randSize + flagSize
	if len(stream) < header {
		return nil, fmt.Errorf("invalid stream size %d", len(stream))
	}

	mixer := mixBytes(stream)
	for i := randSize; i < len(stream); i++ {
		stream[i] ^= mixer[i%randSize]
	}

	cipher.receiver++
	flag := binary.LittleEndian.Uint32(stream[randSize:header])
	if flag != cipher.receiver {
		return nil, fmt.Errorf("invalid index %d != %d", flag, cipher.receiver)
	}

	return stream[header:], nil
}

const (
	randSize = 4
	flagSize = 4
)

func mixBytes(stream []byte) []byte {
	return []byte{
		stream[0]>>6<<6 | stream[1]>>4<<6>>2 | stream[2]>>2<<6>>4 | stream[3]<<6>>6,
		stream[1]>>6<<6 | stream[2]>>4<<6>>2 | stream[3]>>2<<6>>4 | stream[0]<<6>>6,
		stream[2]>>6<<6 | stream[3]>>4<<6>>2 | stream[0]>>2<<6>>4 | stream[1]<<6>>6,
		stream[3]>>6<<6 | stream[0]>>4<<6>>2 | stream[1]>>2<<6>>4 | stream[2]<<6>>6,
	}
}

type IndexMaker struct {
}

func (maker IndexMaker) New() Cipher {
	return &IndexCipher{rand: rand.New(rand.NewSource(time.Now().UnixNano()))}
}
