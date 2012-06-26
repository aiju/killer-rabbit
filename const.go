package main

import (
	"errors"
	"time"
)

const (
	MaxMessage     = 65536
	UpdateInterval = 100 * time.Millisecond
	PingInterval   = time.Second
	PingTimeout    = 20 * PingInterval
	AckTimeout     = 500 * time.Millisecond
	JoinRetry      = 5

	PistolCooldown = time.Second
)

var (
	EConnect = errors.New("connection error")
)
