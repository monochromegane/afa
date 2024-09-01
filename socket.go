package main

import (
	"bytes"
	"encoding/gob"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
)

type Server struct {
	Addr     string
	Listener net.Listener
}

func NewServer(path string) (*Server, error) {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return nil, err
	}

	return &Server{
		Addr: path,
	}, nil
}

func (s *Server) Listen() error {
	listener, err := net.Listen("unix", s.Addr)
	if err != nil {
		return err
	}
	s.Listener = listener
	return nil
}

func (s *Server) Accept() (net.Conn, error) {
	defer s.Listener.Close()
	return s.Listener.Accept()
}

type Client struct {
	Commands []string
}

func NewClient(path string, commands []string) (*Client, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	return &Client{
		Commands: append(commands, absolutePath),
	}, nil
}

func (c *Client) Run() error {
	cmd := exec.Command(c.Commands[0], c.Commands[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

type SocketMessageReader struct {
	Conn    net.Conn
	Decoder *gob.Decoder
	buffer  []byte
	offset  int
}

func NewSocketMessageReader(conn net.Conn) *SocketMessageReader {
	return &SocketMessageReader{
		Conn:    conn,
		Decoder: gob.NewDecoder(conn),
	}
}

func (r *SocketMessageReader) Read(p []byte) (int, error) {
	var err error
	if r.offset >= len(r.buffer) {
		var data []byte
		if err = r.Decoder.Decode(&data); err != nil {
			return 0, nil
		}
		r.buffer = data
		r.offset = 0
	}

	n := copy(p, r.buffer[r.offset:])
	r.offset += n

	if r.offset >= len(r.buffer) && err == io.EOF {
		return n, io.EOF
	}
	return n, nil
}

type SocketMessageWriter struct {
	Conn    net.Conn
	Encoder *gob.Encoder
}

func NewSocketMessageWriter(conn net.Conn) *SocketMessageWriter {
	return &SocketMessageWriter{
		Conn:    conn,
		Encoder: gob.NewEncoder(conn),
	}
}

func (w *SocketMessageWriter) Write(p []byte) (int, error) {
	buf := bytes.NewBuffer(p)
	if err := w.Encoder.Encode(buf.Bytes()); err != nil {
		return 0, err
	}
	return len(p), nil
}

func (w *SocketMessageWriter) Close() error {
	return w.Conn.Close()
}

func (w *SocketMessageWriter) Prompt() error {
	if _, err := w.Write([]byte("__AFA_PROMPT__")); err != nil {
		return err
	}
	return nil
}
