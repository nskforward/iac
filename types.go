package iac

import (
	"context"
	"time"
)

type Runnable interface {
	Run(ctx context.Context, args ...any) Output
}

type Message struct {
	Type      MessageType
	Timestamp time.Time
	Content   []byte
}

type MessageType uint8

const (
	MessageUndefined MessageType = iota
	MessageOK
	MessageError
)
