// Package tools Useful general purpose tools
package tools

import (
	"crypto/rand"
	"math/big"

	log "github.com/sirupsen/logrus"
)

// Try Probe the error and return bool, optionally log the message.
func Try(err error, verbose ...bool) bool {
	if err != nil {
		if len(verbose) > 0 && verbose[0] {
			log.WithError(err).Errorln("")
		}
		return false
	}
	return true
}

// Must Tolerates no errors.
func Must(err error) {
	if !Try(err) {
		log.WithError(err).Panicln("fatal error")
	}
}

// RandInt Return a random in specified range.
func RandInt(min, max int) int {
	bInt, err := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	Must(err)
	bInt = bInt.Add(bInt, big.NewInt(int64(min)))
	return int(bInt.Int64())
}

func IntPtr(n int) *int {
	return &n
}
