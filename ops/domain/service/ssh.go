package service

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"

	"golang.org/x/crypto/ssh"
)

type Connection struct {
	*ssh.Client
	password string
}

func Connect(addr, user, password string) (*Connection, error) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}
	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, err
	}
	return &Connection{conn, password}, nil
}

func (conn *Connection) SendCommands(cmds ...string) ([]byte, error) {
	session, err := conn.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	err = session.RequestPty("xterm", 150, 100, modes)
	if err != nil {
		return []byte{}, err
	}

	in, err := session.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	out, err := session.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	var output []byte
	done := make(chan bool)
	go func(in io.WriteCloser, out io.Reader, output *[]byte) {
		var (
			line string
			r    = bufio.NewReader(out)
		)
		for {
			b, err := r.ReadByte()
			if err != nil {
				break
			}
			*output = append(*output, b)
			if b == byte('\n') {
				line = ""
				continue
			}
			line += string(b)
			if strings.HasPrefix(line, "[sudo] password for ") && strings.HasSuffix(line, ": ") {
				_, err = in.Write([]byte(conn.password + "\n"))
				if err != nil {
					break
				}
			}
		}
		done <- true
	}(in, out, &output)
	cmd := strings.Join(cmds, "; ")
	_, err = session.Output(cmd)
	if err != nil {
		return []byte{}, err
	}
	<-done
	return output, nil
}
