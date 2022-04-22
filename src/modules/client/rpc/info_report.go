package rpc

import (
	"Distributed-System-Awareness-Platform/src/models"
	"github.com/go-kit/log/level"
)

func (r *RpcCli) HostInfoReport(info models.AgentCollectInfo) {
	var msg string
	err := r.GetCli()
	if err != nil {
		level.Error(r.logger).Log("msg", "get cli error", "serverAddr", r.ServerAddr, "err", err)
		return
	}
	err = r.Cli.Call("Server.HostInfoReport", info, &msg)
	if err != nil {
		level.Error(r.logger).Log("msg", "Server.HostInfoReport.error", "serverAddr", r.ServerAddr, "err", err)
		return
	}

}
