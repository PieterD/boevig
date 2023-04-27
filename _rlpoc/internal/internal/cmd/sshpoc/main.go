package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/crypto/ssh"

	"github.com/PieterD/rlpoc/old/internal/internal/pkg/ansi"
	"github.com/PieterD/rlpoc/old/internal/internal/pkg/sshserver"
)

func main() {
	ctx := context.Background()
	s := sshserver.NewSshServer(sshserver.SshServerConfig{
		PrivateKeyPath: "id_rsa",
		ListenAddr:     "0.0.0.0:2222",
	}, &customHandler{})
	if err := s.Run(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

type customHandler struct{}

func (h *customHandler) Allow(ctx context.Context, pubKey ssh.PublicKey) error {
	return nil
}

func (h *customHandler) Session(ctx context.Context, handle sshserver.Handle) error {
	log.Printf("session with pubkey: %s", handle.AuthKey())
	term := ansi.NewGenerator()
	term.ClearScreen()
	s := "press a key"
	rows, cols := handle.Size()
	log.Printf("rows: %d, cols: %d", rows, cols)
	term.CursorPos(rows/2, cols/2-uint(len(s))/2)
	term.Printf("%s", s)

	for {
		if err := term.Flush(handle); err != nil {
			return fmt.Errorf("flushing: %w", err)
		}
		buf := make([]byte, 1)
		if _, err := handle.Read(buf); err != nil {
			return fmt.Errorf("reading: %w", err)
		}
		switch buf[0] {
		case 3:
			return fmt.Errorf("received ETX (end of text): %w", io.EOF)
		case 4:
			return fmt.Errorf("received EOT (end of transmission): %w", io.EOF)
		}
		term.ClearScreen()
		term.CursorPos(1, 1)
		term.BackgroundColor(ansi.Red)
		term.ForegroundColor(ansi.Cyan, false)
		term.Printf("cyan ")
		term.ForegroundColor(ansi.Cyan, true)
		term.Printf("bright cyan ")
		term.ForegroundColor(ansi.Blue, false)
		term.Printf("blue ")
		term.ForegroundColor(ansi.Blue, true)
		term.Printf("bright blue ")
		term.Reset()
		term.CursorPos(2, 1)
		term.Printf("Extended Ascii: copyright symbol %c", 169)
		term.CursorPos(3, 1)
		term.Printf("read: %c (%d)", buf[0], buf[0])
	}
}
