package iac

import (
	"context"
	"fmt"
)

func Run(ctx context.Context, command Runnable, args ...any) error {

	output := command.Run(ctx, args...)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case msg, ok := <-output:
			if !ok {
				return nil
			}
			handleMessage(msg)
		}
	}
}

func handleMessage(msg Message) {
	switch msg.Type {

	case MessageOK:
		fmt.Println(msg.Timestamp.Format("2006-01-02T15:04:05"), "info:", string(msg.Content))

	case MessageError:
		fmt.Println(msg.Timestamp.Format("2006-01-02T15:04:05"), "error:", string(msg.Content))
	}
}
