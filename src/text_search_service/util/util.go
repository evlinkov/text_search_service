package util

import (
	"github.com/satori/go.uuid"
	"hash/crc32"
)

func GenerateUUID() uuid.UUID {
	return uuid.NewV4()
}

func ParseStringToUUID(text string) (uuid.UUID, error) {
	return uuid.FromString(text)
}

func GetHash(text string) uint32 {
	return crc32.ChecksumIEEE([]byte(text))
}
