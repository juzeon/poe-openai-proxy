package util

import (
	"github.com/op/go-logging"
	"math/rand"
	"os"
	"time"
)

var Logger = logging.MustGetLogger("common")

func init() {
	rand.Seed(time.Now().UnixNano())
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	logBackend := logging.NewBackendFormatter(backend, logging.MustStringFormatter("%{level}: %{message}"))
	logging.SetBackend(logBackend)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
