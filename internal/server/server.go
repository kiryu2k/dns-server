package server

import (
	"context"
	"log/slog"
	"net"
	"os"
	"strings"
	"time"

	"github.com/kiryu2k/dns-server/internal/domain"
	"github.com/pkg/errors"
)

const (
	maxBufferSize     = 512
	packetReadTimeout = 500 * time.Millisecond
)

type logger interface {
	DebugContext(ctx context.Context, msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

type dns struct {
	udpAddr *net.UDPAddr
	udpConn *net.UDPConn
	logger  logger
}

func NewDns(host string, port string, logger logger) (*dns, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(host, port))
	if err != nil {
		return nil, errors.WithMessage(err, "resolve udp address")
	}
	return &dns{
		udpAddr: udpAddr,
		udpConn: nil,
		logger:  logger,
	}, nil
}

func (s *dns) ListenAndServe(ctx context.Context) error {
	/* todo: implement connection pool */
	udpConn, err := net.ListenUDP("udp", s.udpAddr)
	if err != nil {
		return errors.WithMessage(err, "listen udp")
	}
	s.udpConn = udpConn

	ch := make(chan packet)
	go s.startReadingPackets(ctx, ch)

	for {
		select {
		case <-ctx.Done():
			return errors.WithMessage(ctx.Err(), "context done")
		case v, ok := <-ch:
			if !ok {
				s.logger.InfoContext(ctx, "done") // todo: understand why it doesn't work with ctrl+c
				return nil
			}
			go s.handlePacket(ctx, v)
		}
	}
}

func (s *dns) Close() {
	if s.udpConn == nil {
		return
	}
	if err := s.udpConn.Close(); err != nil {
		s.logger.WarnContext(context.Background(), errors.WithMessage(err, "close udp connection").Error())
	}
}

type packet struct {
	msg  domain.DnsMessage
	from *net.UDPAddr
}

func (s *dns) startReadingPackets(ctx context.Context, out chan<- packet) {
	defer close(out)

	buf := make([]byte, maxBufferSize)
	for {
		select {
		case <-ctx.Done():
			s.logger.ErrorContext(ctx, errors.WithMessage(ctx.Err(), "context done").Error())
			return
		default:
		}

		err := s.udpConn.SetReadDeadline(time.Now().Add(packetReadTimeout))
		if err != nil {
			s.logger.ErrorContext(ctx, errors.WithMessage(err, "set read deadline").Error())
			continue
		}

		size, source, err := s.udpConn.ReadFromUDP(buf)
		switch {
		case errors.Is(err, os.ErrDeadlineExceeded):
			continue
		case isUdpConnectionClosed(err):
			s.logger.DebugContext(ctx, errors.WithMessage(err, "udp connection is closed").Error())
			return
		case err != nil:
			s.logger.ErrorContext(ctx, errors.WithMessage(err, "read from udp").Error())
			continue
		}

		msg, err := domain.MessageFromBytes(buf[:size])
		if err != nil {
			s.logger.ErrorContext(ctx, errors.WithMessage(err, "message from bytes").Error())
			continue
		}
		s.logger.DebugContext(ctx, "received dns packet", slog.Any("id", msg.Header.Id))

		out <- packet{
			msg:  msg,
			from: source,
		}
	}
}

func (s *dns) handlePacket(ctx context.Context, p packet) {
	response := domain.NewMessage(p.msg.Header.Id).
		AsReply().
		WithQuestion(p.msg.Question).
		Encode()
	if _, err := s.udpConn.WriteToUDP(response, p.from); err != nil {
		s.logger.ErrorContext(ctx, errors.WithMessage(err, "write to udp").Error())
		return
	}
	s.logger.DebugContext(ctx, "send dns packet", slog.Any("id", p.msg.Header.Id))
}

func isUdpConnectionClosed(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "use of closed network connection")
}
