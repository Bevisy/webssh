/*
Copyright © 2021 Binbin Zhang <binbin36520{at}gamil.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"

	sshClient "github.com/bevisy/webssh/util/sshclient"
)

// sshCmd represents the ssh command
var (
	sshCmd = &cobra.Command{
		Use:   "ssh",
		Short: "An sshlike command tool",
		Long:  `An sshlike command tool.`,
		Run: func(cmd *cobra.Command, args []string) {
			sshConnect(args)
		},
	}

	privateKey string
	user       string
	host       string
	port       string
)

func init() {
	rootCmd.AddCommand(sshCmd)

	sshCmd.PersistentFlags().StringVarP(&privateKey, "privateKey", "i", "", "openssh private key")
	sshCmd.PersistentFlags().StringVarP(&host, "server", "s", "localhost", "login server host")
	sshCmd.PersistentFlags().StringVarP(&port, "port", "p", "22", "login serevr port")
	sshCmd.PersistentFlags().StringVarP(&user, "user", "u", "root", "login user")
}

func sshConnect(args []string) {
	var (
		client  *ssh.Client
		session *ssh.Session
		err     error
	)

	server := host + ":" + port
	if privateKey != "" {
		client, session, err = sshClient.ConnectByPrivateKey(privateKey, user, server)
		if err != nil {
			log.Fatalf("connect failed: %s\n", err)
		}
	} else {
		client, session, err = sshClient.ConnectByPassword(user, server)
		if err != nil {
			log.Fatalf("connect failed: %s\n", err)
		}

	}

	defer client.Close()
	defer session.Close()

	if len(args) == 0 {
		// 1. enter terminal on server
		session.Stdout = os.Stdout // 会话输出关联到系统标准输出设备
		session.Stderr = os.Stderr // 会话错误输出关联到系统标准错误输出设备
		session.Stdin = os.Stdin   // 会话输入关联到系统标准输入设备
		modes := ssh.TerminalModes{
			ssh.ECHO:          1,     // 禁用回显（0禁用，1启动）
			ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
			ssh.TTY_OP_OSPEED: 14400, //output speed = 14.4kbaud
		}
		if err = session.RequestPty("linux", 32, 160, modes); err != nil {
			log.Fatalf("request pty error: %s", err.Error())
		}
		if err = session.Shell(); err != nil {
			log.Fatalf("start shell error: %s", err.Error())
		}
		if err = session.Wait(); err != nil {
			log.Fatalf("return error: %s", err.Error())
		}
	} else {
		// 2. excute commands on server
		var cmdstr string
		for i := range args {
			cmdstr = cmdstr + args[i] + " "
		}

		out, err := session.CombinedOutput(cmdstr)
		if err != nil {
			log.Fatalf("excute commnd error: %s", err.Error())
		}
		fmt.Println(string(out))
	}
}
