package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"depulse/internal/vendorcrypto/argon2"
)

const (
	argonTime    uint32 = 3
	argonMemory  uint32 = 64 * 1024
	argonThreads uint8  = 4
	argonSaltLen        = 16
	argonKeyLen  uint32 = 32
)

func validateNewPassword(password string) error {
	if len(password) < 12 {
		return errors.New("password must be at least 12 characters")
	}
	if len(password) > 1024 {
		return errors.New("password is too long")
	}
	return nil
}

func hashPasswordArgon2id(password string) (string, error) {
	if err := validateNewPassword(password); err != nil {
		return "", err
	}
	salt := make([]byte, argonSaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	key := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)
	b64 := base64.RawStdEncoding
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, argonMemory, argonTime, argonThreads, b64.EncodeToString(salt), b64.EncodeToString(key)), nil
}

func verifyPasswordArgon2id(password, encoded string) (bool, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, errors.New("invalid password hash format")
	}
	if parts[2] != "v=19" {
		return false, errors.New("unsupported argon2 version")
	}
	var mem, timeCost uint64
	var threads uint64
	for _, kv := range strings.Split(parts[3], ",") {
		p := strings.SplitN(kv, "=", 2)
		if len(p) != 2 {
			continue
		}
		v, e := strconv.ParseUint(p[1], 10, 32)
		if e != nil {
			return false, e
		}
		switch p[0] {
		case "m":
			mem = v
		case "t":
			timeCost = v
		case "p":
			threads = v
		}
	}
	// Reject attacker-controlled PHC parameters outside the policy envelope before allocating memory.
	if mem < 8*1024 || mem > 256*1024 || timeCost < 1 || timeCost > 10 || threads < 1 || threads > 16 {
		return false, errors.New("argon2 parameters outside policy")
	}
	b64 := base64.RawStdEncoding
	salt, err := b64.DecodeString(parts[4])
	if err != nil || len(salt) < 8 || len(salt) > 64 {
		return false, errors.New("invalid argon2 salt")
	}
	expected, err := b64.DecodeString(parts[5])
	if err != nil || len(expected) < 16 || len(expected) > 64 {
		return false, errors.New("invalid argon2 key")
	}
	actual := argon2.IDKey([]byte(password), salt, uint32(timeCost), uint32(mem), uint8(threads), uint32(len(expected)))
	return subtle.ConstantTimeCompare(actual, expected) == 1, nil
}
