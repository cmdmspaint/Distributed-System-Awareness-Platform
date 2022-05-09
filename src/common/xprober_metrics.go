package common

const (
	// ping
	MetricsNamePingLatency       = `ping_latency_millonseconds` //延迟毫秒数
	MetricsNamePingPackageDrop   = `ping_packageDrop_rate`      //延迟丢包率
	MetricsNamePingTargetSuccess = `ping_target_success`        //ping成功

	// http
	MetricsNameHttpResolvedurationMillonseconds    = `http_resolveDuration_millonseconds`    //DNS解析的时间
	MetricsNameHttpTlsDurationMillonseconds        = `http_tlsDuration_millonseconds`        //tls握手时间
	MetricsNameHttpConnectDurationMillonseconds    = `http_connectDuration_millonseconds`    //链接建立时间
	MetricsNameHttpProcessingDurationMillonseconds = `http_processingDuration_millonseconds` //服务端处理时间
	MetricsNameHttpTransferDurationMillonseconds   = `http_transferDuration_millonseconds`   //任务处理时间
	MetricsNameHttpInterfaceSuccess                = `http_interface_success`
)
