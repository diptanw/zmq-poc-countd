package zeromq

import (
	"fmt"
	"io"

	"github.com/pebbe/zmq4"
)

// SocketReader is a struct that implements io.ReadCloser interface that reads
// messages from the ZMQ socket. The reader returns EOF when the message is fully
// read, reader should be recycled for the next message.
type SocketReader struct {
	sock *zmq4.Socket
}

// Bind binds the ZMQ socket listener to the provided address.
func Bind(bindAddr string) (SocketReader, error) {
	sock, err := bindSocket(bindAddr, zmq4.PULL)
	if err != nil {
		return SocketReader{}, err
	}

	return SocketReader{
		sock: sock,
	}, nil
}

// Read pulls every ZMQ socket message as individual read operation
// ending it with io.EOF marker.
func (r SocketReader) Read(p []byte) (int, error) {
	b, err := r.sock.RecvBytes(zmq4.SNDMORE)
	if err != nil {
		return 0, fmt.Errorf("zeromq: reading bytes: %w", err)
	}

	return copy(p, b), io.EOF
}

// Close closes the ZMQ socket.
func (r SocketReader) Close() error {
	if err := r.sock.Close(); err != nil {
		return fmt.Errorf("zeromq: closing socket: %w", err)
	}

	return nil
}

func bindSocket(bindAddr string, t zmq4.Type) (*zmq4.Socket, error) {
	zmqCtx, err := zmq4.NewContext()
	if err != nil {
		return nil, fmt.Errorf("zeromq: creating context: %w", err)
	}

	zmqSock, err := zmqCtx.NewSocket(t)
	if err != nil {
		return nil, fmt.Errorf("zeromq: creating socket: %w", err)
	}

	if err := zmqSock.Bind(bindAddr); err != nil {
		return nil, fmt.Errorf("zeromq: binding socket: %w", err)
	}

	return zmqSock, nil
}
