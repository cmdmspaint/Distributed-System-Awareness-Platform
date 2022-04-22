package rpc

import (
	"Distributed-System-Awareness-Platform/src/models"
	"log"
)

func (*Server) HostInfoReport(input models.AgentCollectInfo, output *string) error {
	log.Printf("[HostInfoReport][input:%+v]", input)
	*output = "I know"
	return nil
}
