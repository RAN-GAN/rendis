package main

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

type Generator struct {
	rng *rand.Rand
}

func NewGenerator() *Generator {
	return &Generator{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (g *Generator) RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[g.rng.Intn(len(charset))]
	}
	return string(b)
}

func (g *Generator) RandomKey() string {
	return "bench:" + g.RandomString(16)
}
