package sshclient

import (
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

func ConnectByPrivateKey(keyPath, user, host string) (*ssh.Client, *ssh.Session, error) {
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, nil, err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		// TODO: do not use ssh.InsecureIgnoreHostKey()
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}

	return client, session, nil
}

func ConnectByPassword(user, host string) (*ssh.Client, *ssh.Session, error) {
	var pass string
	fmt.Print("Password: ")
	fmt.Scanf("%s\n", &pass)

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		// TODO: do not use ssh.InsecureIgnoreHostKey()
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}

	return client, session, nil
}
