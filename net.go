package gojob

import (
	"time"
	"fmt"

	"golang.org/x/crypto/ssh"
	"github.com/pkg/sftp"
)

type Conn struct {
	Connect *ssh.Client
	Client	*sftp.Client
}

func NewConnect(user, password, host, port string) (*Conn, error) {
	// ssh client config
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		// allow any host key to be used (non-prod)
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),

		// verify host public key
		// HostKeyCallback: ssh.FixedHostKey(hostKey),
		// optional tcp connect timeout
		Timeout: 5 * time.Second,
	}

	// connect
	sshConn, err := ssh.Dial("tcp", host+":"+port, config)
	if err != nil {
		return nil, fmt.Errorf("gojob: ssh dial %v", err)
	}

	client, err := sftp.NewClient(sshConn)
	if err != nil {
		return nil, fmt.Errorf("gojob: sftp client %v", err)
	}

	conn := &Conn{
		Connect: sshConn,
		Client: client,
	}

	return conn, nil
}

func (c *Conn) Session() (*ssh.Session, error) {
	conn := c.Connect
	session, err := conn.NewSession()
	if err != nil {
		return nil, fmt.Errorf("new session %v", err)
	}
	return session, nil
}

func (c *Conn) Close() {
	c.Client.Close()
	c.Connect.Close()
}
