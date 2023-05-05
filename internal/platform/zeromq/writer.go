package zeromq

import (
	"fmt"

	"github.com/pebbe/zmq4"
)

// SocketWriter is a struct that implements io.WriteCloser interface that writes
// data to ZMQ socket as individual push messages.
type SocketWriter struct {
	sock *zmq4.Socket
}

// Connect connects to the remote ZMQ socket.
func Connect(addr string) (*SocketWriter, error) {
	sock, err := connectSocket(addr, zmq4.PUSH)
	if err != nil {
		return nil, fmt.Errorf("zeromq: connecting socket: %w", err)
	}

	return &SocketWriter{
		sock: sock,
	}, nil
}

// Write pushes every message as individual write operation to ZMQ socket.
func (r *SocketWriter) Write(b []byte) (int, error) {
	n, err := r.sock.SendMessage(b)
	if err != nil {
		return 0, fmt.Errorf("zeromq: writing bytes: %w", err)
	}

	return n, nil
}

// Close closes the ZMQ socket connection.
func (r *SocketWriter) Close() error {
	if err := r.sock.Close(); err != nil {
		return fmt.Errorf("zeromq: closing socket: %w", err)
	}

	return nil
}

func connectSocket(addr string, t zmq4.Type) (*zmq4.Socket, error) {
	zmqCtx, err := zmq4.NewContext()
	if err != nil {
		return nil, fmt.Errorf("zeromq: creating context: %w", err)
	}

	zmqSock, err := zmqCtx.NewSocket(t)
	if err != nil {
		return nil, fmt.Errorf("zeromq: creating socket: %w", err)
	}

	if err := zmqSock.Connect(addr); err != nil {
		return nil, fmt.Errorf("zeromq: connecting socket: %w", err)
	}

	return zmqSock, nil
}
