/*
* Copyright (C) 2021 THL A29 Limited, a Tencent company.  All rights reserved.
* This source code is licensed under the Apache License Version 2.0.
 */

package cmds

import (
	"github.com/lsutils/nocalhost/cmd/nhctl/cmds/common"
	"github.com/lsutils/nocalhost/internal/nhctl/daemon_client"
	"github.com/lsutils/nocalhost/internal/nhctl/daemon_server/command"
	"github.com/lsutils/nocalhost/internal/nhctl/utils"
	"github.com/lsutils/nocalhost/pkg/nhctl/log"
	"github.com/spf13/cobra"
	"io/ioutil"
)

func init() {
	kubeconfigCmd.AddCommand(kubeconfigAddCmd)
}

// Add kubeconfig
var kubeconfigAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add kubeconfig",
	Long:  `Add kubeconfig`,
	Run: func(cmd *cobra.Command, args []string) {
		daemonClient, err := daemon_client.GetDaemonClient(utils.IsSudoUser())
		if err != nil {
			log.FatalE(err, "")
		}
		if err := common.Prepare(); err != nil {
			return
		}
		if bytes, err := ioutil.ReadFile(common.KubeConfig); err == nil {
			if err = daemonClient.SendKubeconfigOperationCommand(bytes, common.NameSpace, command.OperationAdd); err != nil {
				log.Info(err)
			}
		}
	},
}
