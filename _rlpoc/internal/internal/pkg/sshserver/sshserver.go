package sshserver

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/sync/errgroup"
)

const defaultSessionEstablishmentTimeout = time.Second

type ChannelHandler interface {
	Allow(ctx context.Context, pubKey ssh.PublicKey) error
	Session(ctx context.Context, h Handle) error
}

type Handle interface {
	io.Reader
	io.Writer
	AuthKey() string
	Size() (rows, cols uint)
}

type SshServerConfig struct {
	PrivateKeyPath string
	ListenAddr     string
}

type SshServer struct {
	privateKeyPath string
	listenAddr     string
	handler        ChannelHandler
}

func NewSshServer(cfg SshServerConfig, handler ChannelHandler) *SshServer {
	return &SshServer{
		privateKeyPath: cfg.PrivateKeyPath,
		listenAddr:     cfg.ListenAddr,
		handler:        handler,
	}
}

func (s *SshServer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	config := &ssh.ServerConfig{
		NoClientAuth: false,
		PublicKeyCallback: func(conn ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {
			fingerprint := ssh.FingerprintSHA256(pubKey)
			if err := s.handler.Allow(ctx, pubKey); err != nil {
				log.Printf("error allowing public key %s: %v", fingerprint, err)
				return nil, fmt.Errorf("connection disapproved: %w", err)
			}
			return &ssh.Permissions{
				// Record the public key used for authentication.
				Extensions: map[string]string{
					"pubkey-fp": fingerprint,
				},
			}, nil
		},
	}
	privateBytes, err := ioutil.ReadFile(s.privateKeyPath)
	if err != nil {
		return fmt.Errorf("reading private key: %w", err)
	}
	privateKey, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		return fmt.Errorf("parsing private key: %w", err)
	}
	config.AddHostKey(privateKey)

	listener, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return fmt.Errorf("creating tcp listener: %w", err)
	}
	defer listener.Close()

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		<-ctx.Done()
		log.Printf("closing listener")
		if err := listener.Close(); err != nil {
			return fmt.Errorf("closing listener: %w", err)
		}
		return ctx.Err()
	})
	eg.Go(func() error {
		defer log.Printf("ending listener")
		for {
			conn, err := listener.Accept()
			if err != nil {
				return fmt.Errorf("accepting connection: %w", err)
			}
			eg.Go(func() error {
				defer log.Printf("ending connection processor")
				err := s.processConn(ctx, conn, config)
				log.Printf("error processing connection: %v", err)
				return nil
			})
		}
	})
	if err := eg.Wait(); err != nil {
		return fmt.Errorf("waiting: %w", err)
	}
	return nil
}

func (s *SshServer) processConn(ctx context.Context, netConn net.Conn, config *ssh.ServerConfig) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)
	sshConn, connChannelSource, connRequestSource, err := ssh.NewServerConn(netConn, config)
	if err != nil {
		return fmt.Errorf("creating new server conn: %w", err)
	}
	defer sshConn.Close()
	pubKey := sshConn.Permissions.Extensions["pubkey-fp"]
	log.Printf("login with key %s", pubKey)
	eg.Go(func() error {
		<-ctx.Done()
		log.Printf("closing ssh conn")
		if err := sshConn.Close(); err != nil {
			return fmt.Errorf("closing ssh conn: %w", err)
		}
		return ctx.Err()
	})
	eg.Go(func() error {
		defer log.Printf("ending conn processor")
		channelAccepted := false
		h := &handle{
			ReadWriter: nil,
			pubKey:     pubKey,
		}
		for {
			var connChannel ssh.NewChannel
			var ok bool
			select {
			case <-ctx.Done():
				return ctx.Err()
			case connRequest, ok := <-connRequestSource:
				if !ok {
					return fmt.Errorf("conn request source closed")
				}
				log.Printf("connection made a request: %s", connRequest.Type)
				if connRequest.WantReply {
					if err := connRequest.Reply(false, nil); err != nil {
						return fmt.Errorf("replying to connection request %s: %w", connRequest.Type, err)
					}
				}
				continue
			case connChannel, ok = <-connChannelSource:
				if !ok {
					return fmt.Errorf("conn channel source closed")
				}
			}
			if connChannel.ChannelType() != "session" {
				if err := connChannel.Reject(ssh.UnknownChannelType, "unsupported channel type"); err != nil {
					return fmt.Errorf("rejecting channel: %w", err)
				}
				continue
			}
			if channelAccepted {
				if err := connChannel.Reject(ssh.Prohibited, "only one session channel allowed at a time"); err != nil {
					return fmt.Errorf("rejecting second channel: %w", err)
				}
				continue
			}
			channel, channelRequestSource, err := connChannel.Accept()
			if err != nil {
				return fmt.Errorf("accepting channel request: %w", err)
			}
			h.ReadWriter = channel
			channelAccepted = true
			channelReadySignal := make(chan struct{})
			eg.Go(func() error {
				defer log.Printf("ending channel request processor")
				if err := processChannelRequests(ctx, channelRequestSource, h, channelReadySignal); err != nil {
					return fmt.Errorf("processing channel requests: %w", err)
				}
				return fmt.Errorf("processing channel requests returned no error")
			})
			eg.Go(func() error {
				defer log.Printf("ending session handler")
				defer channel.Close()
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(defaultSessionEstablishmentTimeout):
					return fmt.Errorf("client took too long to establish session")
				case <-channelReadySignal:
				}
				if err := s.handler.Session(ctx, h); err != nil {
					return fmt.Errorf("handling session: %w", err)
				}
				return fmt.Errorf("session handler returned no error")
			})
		}
	})
	if err := eg.Wait(); err != nil {
		log.Printf("wait error: %v", err)
		return fmt.Errorf("waiting: %w", err)
	}
	log.Printf("the wait is over")
	return nil
}

type sizeSetter interface {
	setSize(rows, cols uint)
}

func processChannelRequests(ctx context.Context, requestSource <-chan *ssh.Request, sizeSetter sizeSetter, readySignal chan<- struct{}) error {
	shellAccepted := false
	ptyAccepted := false
	for {
		if readySignal != nil && shellAccepted && ptyAccepted {
			close(readySignal)
			readySignal = nil
		}
		var request *ssh.Request
		var ok bool
		select {
		case <-ctx.Done():
			return ctx.Err()
		case request, ok = <-requestSource:
			if !ok {
				return fmt.Errorf("request source closed")
			}
		}
		log.Printf("channel made a request: %s (%t) %v", request.Type, request.WantReply, request.Payload)
		ok = false
		switch {
		case request.Type == "shell":
			if shellAccepted {
				log.Printf("another shell request received, but only expected one")
				break
			}
			shellAccepted = true
			ok = true
		case request.Type == "pty-req":
			if ptyAccepted {
				log.Printf("another pty-req request received, but only expected one")
				break
			}
			ptyAccepted = true
			ok = true
			p, err := parsePtyReqPayload(request.Payload)
			if err != nil {
				log.Printf("pty-req parser error: %v", err)
				return fmt.Errorf("parsing pty-req payload: %w", err)
			}
			sizeSetter.setSize(uint(p.Rows), uint(p.Columns))
			log.Printf("pty-req: %#v", p)
		case request.Type == "window-change":
			p, err := parseWindowChangePayload(request.Payload)
			if err != nil {
				log.Printf("window-change parser error: %v", err)
				return fmt.Errorf("parsing window-change payload: %w", err)
			}
			sizeSetter.setSize(uint(p.Rows), uint(p.Columns))
			log.Printf("window-change: %#v", p)
		}
		if request.WantReply {
			if err := request.Reply(ok, nil); err != nil {
				log.Fatalf("replying to %s: %v", request.Type, err)
			}
		}
	}
}
