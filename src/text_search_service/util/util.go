package util

import (
	"github.com/satori/go.uuid"
	"hash/crc32"
)

func GenerateUUID() uuid.UUID {
	return uuid.NewV4()
}

func GetHash(text string) uint32 {
	return crc32.ChecksumIEEE([]byte(text))
}
