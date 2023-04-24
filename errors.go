package cerra

import "errors"

var (
	ErrEmtpyQueue    = errors.New("cerra: empty queue")
	ErrInActiveQueue = errors.New("cerra: inactive queue")
	ErrQueueClosed   = errors.New("cerra: queue closed")
)
