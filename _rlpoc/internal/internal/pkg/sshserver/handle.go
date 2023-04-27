package sshserver

import (
	"io"
	"sync"
)

type handle struct {
	io.ReadWriter
	pubKey string

	sizeLock   sync.Mutex
	rows, cols uint
}

func (h *handle) setSize(rows, cols uint) {
	h.sizeLock.Lock()
	defer h.sizeLock.Unlock()

	h.rows, h.cols = rows, cols
}

func (h *handle) Size() (rows, cols uint) {
	h.sizeLock.Lock()
	defer h.sizeLock.Unlock()

	rows, cols = h.rows, h.cols
	return rows, cols
}

func (h *handle) AuthKey() string {
	return h.pubKey
}

var _ Handle = &handle{}
