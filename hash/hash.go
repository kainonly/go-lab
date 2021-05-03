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
	def = argon2Option{
		time:    uint32(4),
		memory:  uint32(65536),
		threads: uint8(1),
	}

	ErrInvalidHash         = errors.New(`the encoded hash is not in the correct format`)
	ErrIncompatibleVersion = errors.New(`incompatible version of argon2id`)
)

type Option interface {
	apply(param *argon2Option)
}

type argon2Option struct {

	// Maximum memory (in kibibytes) that may be used to compute the Argon2 hash
	time uint32

	// Maximum amount of time it may take to compute the Argon2 hash
	memory uint32

	// Number of threads to use for computing the Argon2 hash
	threads uint8
}

func Time(value uint32) Option {
	return time(value)
}

type time uint32

func (c time) apply(opt *argon2Option) {
	opt.time = uint32(c)
}

func Memory(value uint32) Option {
	return memory(value)
}

type memory uint32

func (c memory) apply(opt *argon2Option) {
	opt.memory = uint32(c)
}

func Threads(value uint8) Option {
	return threads(value)
}

type threads uint8

func (c threads) apply(opt *argon2Option) {
	opt.threads = uint8(c)
}

// Make Use argon2id hash to generate user password
//	@param `password` string user password
//	@param `options` ...Option the algorithm using the memory, time, and threads options
//	@return `hashedPassword` string hash password
func Make(password string, options ...Option) (hashedPassword string, err error) {
	salt := make([]byte, 16)
	if _, err = rand.Read(salt); err != nil {
		return
	}
	option := def
	for _, value := range options {
		value.apply(&option)
	}
	hash := argon2.IDKey([]byte(password), salt, option.time, option.memory, option.threads, 32)
	hashedPassword = "$argon2id$v=" + strconv.Itoa(argon2.Version)
	hashedPassword += "$m=" + strconv.Itoa(int(option.memory))
	hashedPassword += ",t=" + strconv.Itoa(int(option.time))
	hashedPassword += ",p=" + strconv.Itoa(int(option.threads))
	hashedPassword += "$" + base64.RawStdEncoding.EncodeToString(salt)
	hashedPassword += "$" + base64.RawStdEncoding.EncodeToString(hash)
	return
}

// Verify Verifying that a password matches a hash
//	@param `password` string user password
//	@param `hashedPassword` string hash password
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
	option := argon2Option{
		memory:  uint32(memory),
		time:    uint32(time),
		threads: uint8(threads),
	}
	var decodeSalt []byte
	if decodeSalt, err = base64.RawStdEncoding.DecodeString(args[6]); err != nil {
		return false, err
	}
	var hash []byte
	if hash, err = base64.RawStdEncoding.DecodeString(args[7]); err != nil {
		return false, err
	}
	newHash := argon2.IDKey([]byte(password), decodeSalt, option.time, option.memory, option.threads, 32)
	if subtle.ConstantTimeCompare(hash, newHash) == 1 {
		return true, nil
	}
	return false, nil
}
