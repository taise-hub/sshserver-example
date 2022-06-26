package main

import (
	"io"
	"strings"
	"os/exec"
	"golang.org/x/crypto/ssh"
	"github.com/creack/pty"
	"log"
	"net"
)

func main() {
	config := &ssh.ServerConfig{
		NoClientAuth: true,
	}

	private, err := ssh.ParsePrivateKey([]byte(key))
	if err != nil {
		log.Fatal("Failed to parse private key: ", err)
	}

	config.AddHostKey(private)

	listener, err := net.Listen("tcp", "0.0.0.0:2022")
	if err != nil {
		log.Fatal("failed to listen for connection: ", err)
	}

	// SSHd
	for {
		nConn, err := listener.Accept()
		if err != nil {
			log.Fatal("failed to accept incoming connection: ", err)
		}

		_, chans, reqs, err := ssh.NewServerConn(nConn, config)

		if err != nil {
			log.Fatal("failed to handshake: ", err)
		}

		go ssh.DiscardRequests(reqs)
		go handleChannels(chans)
	}
}

func handleChannels(chans <-chan ssh.NewChannel) {
	for newChannel := range chans {
		go handleChannel(newChannel)
	}
}

func handleChannel(newChan ssh.NewChannel) {
	if newChan.ChannelType() != "session" {
		newChan.Reject(ssh.UnknownChannelType, "unknown channel type")
		return
	}

	channel, _, err := newChan.Accept()
	if err != nil {
		log.Fatalf("Could not accept channel: %v", err)
	}
	log.Printf("open channel: %s \n", newChan.ChannelType())

	// shell
	shell(channel)
}

// 命名適当
func shell(channel ssh.Channel) {
	defer channel.Close()
	c := exec.Command("bash")
	ptmx, err := pty.Start(c)
	io.Copy(channel, strings.NewReader("Welecome to my SSH Server!!! build by Taise\n"))
    if err != nil {
		log.Fatal("Could not start pty: ", err)
    }
    defer func() { _ = ptmx.Close() }() // Best effort.

	// 関数が戻る前に次の入力を待ち続けてくれるからこれでOK
	go func() { _, _ = io.Copy(ptmx, channel) }()
    _, _ = io.Copy(channel, ptmx)
}

var key = `
-----BEGIN OPENSSH PRIVATE KEY-----
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
bAECAw==
-----END OPENSSH PRIVATE KEY-----`
