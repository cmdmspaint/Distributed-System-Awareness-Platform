package memoryindex

import (
	"Distributed-System-Awareness-Platform/src/common"
	"Distributed-System-Awareness-Platform/src/modules/server/config"
	"Distributed-System-Awareness-Platform/src/modules/server/metric"
	"context"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	mem "github.com/ning1875/inverted-index"
	"github.com/ning1875/inverted-index/index"
	"github.com/prometheus/client_golang/prometheus"
	"strings"
	"sync"
	"time"
)

type ResourceIndexer interface {
	FlushIndex()                          // 刷新索引的方法
	GetIndexReader() *mem.HeadIndexReader // 获取内置的索引reader
	GetLogger() log.Logger
}

var indexContainer = make(map[string]ResourceIndexer)

func iRegister(name string, ri ResourceIndexer) {
	indexContainer[name] = ri
}

func JudgeResourceIndexExists(name string) bool {
	_, ok := indexContainer[name]
	return ok
}

func Init(logger log.Logger, ims []*config.IndexModuleConf) {

	loadNum := 0
	loadResource := make([]string, 0)
	for _, i := range ims {
		if !i.Enable {

			continue
		}
		level.Info(logger).Log("msg", "memoryindex.init", "name", i.ResourceName)
		loadNum += 1
		loadResource = append(loadResource, i.ResourceName)
		switch i.ResourceName {
		case common.RESOURCE_HOST:
			mi := &HostIndex{
				Ir:      mem.NewHeadReader(),
				Logger:  logger,
				Modulus: i.Modulus,
				Num:     i.Num,
			}
			iRegister(i.ResourceName, mi)
		case common.RESOURCE_RDS:
			mi := &HostIndex{
				Ir:      mem.NewHeadReader(),
				Logger:  logger,
				Modulus: i.Modulus,
				Num:     i.Num,
			}
			iRegister(i.ResourceName, mi)

		}
	}
	level.Info(logger).Log("msg", "mem-index.init.summary", "loadNum", loadNum, "detail", strings.Join(loadResource, " "))
}

func GetAllResourceIndexReader() (make map[string]ResourceIndexer) {
	return indexContainer

}

func GetMatchIdsByIndex(req common.ResourceQueryReq) (matchIds []uint64) {
	ri, ok := indexContainer[req.ResourceType]
	if !ok {
		return
	}
	matcher := common.FormatLabelMatcher(req.Labels)

	p, err := mem.PostingsForMatchers(ri.GetIndexReader(), matcher...)
	if err != nil {
		level.Error(ri.GetLogger()).Log("msg", "ii.PostingsForMatchers.error", "ResourceType", req.ResourceType, "err", err)
		return
	}
	matchIds, err = index.ExpandPostings(p)
	if err != nil {
		level.Error(ri.GetLogger()).Log("msg", "index.ExpandPostings.error", "ResourceType", req.ResourceType, "err", err)
		return
	}
	return
}

func RevertedIndexSyncManager(ctx context.Context, logger log.Logger) error {
	level.Info(logger).Log("msg", "RevertedIndexSyncManager.start", "resource_num", len(indexContainer))
	ticker := time.NewTicker(15 * time.Second)
	doIndexFlush()
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			level.Info(logger).Log("msg", "RevertedIndexSyncManager.exit.receive_quit_signal", "resource_num", len(indexContainer))
			return nil
		case <-ticker.C:
			level.Info(logger).Log("msg", "doIndexFlush.cron", "resource_num", len(indexContainer))

			doIndexFlush()
		}
	}

}

func doIndexFlush() {
	var wg sync.WaitGroup
	wg.Add(len(indexContainer))
	for name, ir := range indexContainer {
		name := name
		ir := ir
		go func() {
			defer wg.Done()
			start := time.Now()
			ir.FlushIndex()
			metric.IndexFlushDuration.With(prometheus.Labels{common.LABEL_RESOURCE_TYPE: name}).Set(float64(time.Since(start).Seconds()))

		}()
	}
	wg.Wait()
}

func GetResourceIndexReader(name string) (bool, ResourceIndexer) {
	ri, ok := indexContainer[name]
	return ok, ri

}
