package models

//机器上shell 采集到的字段
type AgentCollectInfo struct {
	SN       string `json:"sn"`       //sn 号
	CPU      string `json:"cpu"`      //cpu 核数
	Mem      string `json:"mem"`      //内存多少g
	Disk     string `json:"disk"`     //磁盘多少g
	IpAddr   string `json:"ip_addr"`  //ip
	HostName string `json:"hostname"` //hostname
}
