package main

import (
	"errors"
	"time"
	"math"
)

const (
	MaxMessage     = 65536
	UpdateInterval = 100 * time.Millisecond
	PingInterval   = time.Second
	PingTimeout    = 20 * PingInterval
	AckTimeout     = 500 * time.Millisecond
	JoinRetry      = 5
	PlayerSize = 50
	MonsterSize = 80

	PistolCooldown = time.Second

	deg = math.Pi / 180
)

var (
	EConnect = errors.New("connection error")
)
