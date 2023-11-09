/*
* Copyright (C) 2021 THL A29 Limited, a Tencent company.  All rights reserved.
* This source code is licensed under the Apache License Version 2.0.
 */

package cmds

import (
	"bufio"
	"fmt"
	"github.com/lsutils/nocalhost/cmd/nhctl/cmds/common"
	"github.com/lsutils/nocalhost/internal/nhctl/daemon_client"
	"github.com/lsutils/nocalhost/internal/nhctl/daemon_server/command"
	"github.com/lsutils/nocalhost/internal/nhctl/vpn/driver"
	"github.com/lsutils/nocalhost/internal/nhctl/vpn/util"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"k8s.io/client-go/tools/clientcmd"
	"strings"
)

var workloads string

func init() {
	connectCmd.Flags().StringVar(&common.KubeConfig, "kubeconfig", clientcmd.RecommendedHomeFile, "kubeconfig")
	connectCmd.Flags().StringVarP(&common.NameSpace, "namespace", "n", "", "namespace")
	connectCmd.Flags().StringVar(&workloads, "workloads", "", "workloads, like: services/tomcat, deployment/nginx, replicaset/tomcat...")
	vpnCmd.AddCommand(connectCmd)
}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "connect",
	Long:  `connect`,
	PreRun: func(*cobra.Command, []string) {
		util.InitLogger(util.Debug)
		if util.IsWindows() {
			_ = driver.InstallWireGuardTunDriver()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// if not sudo and sudo daemon is not running, needs sudo permission
		if !util.IsAdmin() && !util.IsSudoDaemonServing() {
			if err := util.RunWithElevated(); err != nil {
				log.Warn(err)
				return
			}
		}
		_, err := daemon_client.GetDaemonClient(true)
		if err != nil {
			log.Warn(err)
			return
		}
		client, err := daemon_client.GetDaemonClient(false)
		if err != nil {
			log.Warn(err)
			return
		}
		must(common.Prepare())
		err = client.SendVPNOperateCommand(common.KubeConfig, common.NameSpace, command.Connect, workloads, f)
		if err != nil {
			log.Warn(err)
		}
	},
}

var f = func(reader io.Reader) error {
	stream := bufio.NewReader(reader)
	for {
		if line, _, err := stream.ReadLine(); errors.Is(err, io.EOF) {
			return nil
		} else {
			if len(line) == 0 {
				continue
			}
			if strings.Contains(string(line), util.EndSignOK) {
				return nil
			} else if strings.Contains(string(line), util.EndSignFailed) {
				return nil
			}
			fmt.Println(string(line))
		}
	}
}
