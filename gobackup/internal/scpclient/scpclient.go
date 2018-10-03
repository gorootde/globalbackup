package scpclient

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

//SCPClient is a client for transfering files via SCP to another machine
type SCPClient struct {
	clientConfig *ssh.ClientConfig
	port         int
	host         string
	session      *ssh.Session
}

//New creates a new client
func New(username string, host string, port int) *SCPClient {
	signer, _ := ssh.ParsePrivateKey([]byte(privateKey))
	clientConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}
	return &SCPClient{clientConfig: clientConfig, host: host, port: port}
}

//Connect to the server
func (instance *SCPClient) Connect() {
	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", instance.host, instance.port), instance.clientConfig)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	instance.session = session
}

//Copy a file to the server
func (instance *SCPClient) Copy(src string, dst string) {
	w, _ := instance.session.StdinPipe()
	defer w.Close()
	fmt.Fprintln(w, "D0755", 0, dst) // mkdir
	fmt.Fprintln(w, "C0644", len(content), "testfile1")
	fmt.Fprint(w, content)
	fmt.Fprint(w, "\x00") // transfer end with \x00
	if err := instance.session.Run("/usr/bin/scp -tr ./"); err != nil {
		panic("Failed to run: " + err.Error())
	}
}

//Disconnect from server
func (instance *SCPClient) Disconnect() {
	instance.session.Close()
}
