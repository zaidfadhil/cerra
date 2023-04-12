package goatq

import "errors"

var (
	ErrEmtpyQueue    = errors.New("goatq: empty queue")
	ErrInActiveQueue = errors.New("goatq: inactive queue")
)
