package rpc

import (
	"Distributed-System-Awareness-Platform/src/models"
	"github.com/toolkits/pkg/logger"
)

// input = agent.hostName
func (*Server) LogJobSync(input string, output *[]*models.LogStrategy) error {
	ljs, _ := models.LogJobGets("id>0")

	*output = ljs
	logger.Infof("LogJobSync.call.receive ljs :%v %v %v", ljs, input, output)
	return nil

}
