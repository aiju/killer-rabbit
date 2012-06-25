package main

import (
	"errors"
	"time"
)

const (
	MaxMessage     = 65536
	UpdateInterval = 100 * time.Millisecond
	AckTimeout     = 500 * time.Millisecond
	JoinRetry      = 5
)

var (
	EConnect = errors.New("connection error")
)
