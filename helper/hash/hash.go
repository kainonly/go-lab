package hash

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/argon2"
	"regexp"
	"strconv"
)

var (
	DefaultTime    = uint32(4)
	DefaultMemory  = uint32(64 * 1024)
	DefaultThreads = uint8(1)
)

var (
	ErrInvalidHash         = errors.New(`the encoded hash is not in the correct format`)
	ErrIncompatibleVersion = errors.New(`incompatible version of argon2`)
)

type Option struct {
	Time    uint32
	Memory  uint32
	Threads uint8
}

func Make(password string, option Option) (hashedPassword string, err error) {
	salt := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		return
	}
	if option.Time == 0 {
		option.Time = DefaultTime
	}
	if option.Memory == 0 {
		option.Memory = DefaultMemory
	}
	if option.Threads == 0 {
		option.Threads = DefaultThreads
	}
	hash := argon2.IDKey([]byte(password), salt, option.Time, option.Memory, option.Threads, 32)
	hashedPassword = "$argon2id$v=" + strconv.Itoa(argon2.Version)
	hashedPassword += "$m=" + strconv.Itoa(int(option.Memory))
	hashedPassword += ",t=" + strconv.Itoa(int(option.Time))
	hashedPassword += ",p=" + strconv.Itoa(int(option.Threads))
	hashedPassword += "$" + base64.RawStdEncoding.EncodeToString(salt)
	hashedPassword += "$" + base64.RawStdEncoding.EncodeToString(hash)
	return
}

func Verify(password string, hashedPassword string) (result bool, err error) {
	args := regexp.
		MustCompile(`^\$(\w+)\$v=(\d+)\$m=(\d+),t=(\d+),p=(\d+)\$(\S+)\$(\S+)`).
		FindStringSubmatch(hashedPassword)
	if len(args) != 8 {
		return false, ErrInvalidHash
	}
	if args[1] != `argon2id` {
		return false, ErrInvalidHash
	}
	var version int
	if version, err = strconv.Atoi(args[2]); err != nil {
		return false, err
	}
	if version != argon2.Version {
		return false, ErrIncompatibleVersion
	}
	var memory uint64
	if memory, err = strconv.ParseUint(args[3], 10, 32); err != nil {
		return false, err
	}
	var time uint64
	if time, err = strconv.ParseUint(args[4], 10, 32); err != nil {
		return false, err
	}
	var threads int
	if threads, err = strconv.Atoi(args[5]); err != nil {
		return false, err
	}
	option := Option{
		Memory:  uint32(memory),
		Time:    uint32(time),
		Threads: uint8(threads),
	}
	var decodeSalt []byte
	if decodeSalt, err = base64.RawStdEncoding.DecodeString(args[6]); err != nil {
		return false, err
	}
	var hash []byte
	if hash, err = base64.RawStdEncoding.DecodeString(args[7]); err != nil {
		return false, err
	}
	newHash := argon2.IDKey([]byte(password), decodeSalt, option.Time, option.Memory, option.Threads, 32)
	if subtle.ConstantTimeCompare(hash, newHash) == 1 {
		return true, nil
	}
	return false, nil
}
