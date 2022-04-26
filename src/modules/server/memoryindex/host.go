package memoryindex

import (
	"Distributed-System-Awareness-Platform/src/common"
	"Distributed-System-Awareness-Platform/src/models"
	"encoding/json"
	"fmt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	mem "github.com/ning1875/inverted-index"
	"github.com/ning1875/inverted-index/labels"
	"strconv"
	"strings"
)

type HostIndex struct {
	Ir      *mem.HeadIndexReader
	Logger  log.Logger
	Modulus int // 静态分片
	Num     int
}

func (hi *HostIndex) FlushIndex() {
	// 数个数
	r := new(models.ResourceHost)
	total := int(r.Count())
	ids := ""
	for i := 0; i < total; i++ {

		// 先写单点逻辑
		if hi.Modulus == 0 {
			ids += fmt.Sprintf("%d,", i)
			continue
		}
		// 分片匹配中了 ，keep的逻辑
		if i%hi.Modulus == hi.Num {
			ids += fmt.Sprintf("%d,", i)
			continue
		}
	}
	ids = strings.TrimRight(ids, ",")
	inSql := fmt.Sprintf("id in (%s) ", ids)
	objs, err := models.ResourceHostGetMany(inSql)
	if err != nil {
		return
	}
	thisH := mem.NewHeadReader()

	//自动刷path的map
	thisGPAS := map[string]struct{}{}

	for _, item := range objs {
		m := make(map[string]string)
		m["hash"] = item.Hash
		tags := make(map[string]string)
		// 数组型 内网ips 公网ips 安全组
		prIps := []string{}
		puIps := []string{}

		// 当个kv
		m["uid"] = item.Uid
		m["name"] = item.Name
		m["cloud_provider"] = item.CloudProvider
		m["charging_mode"] = item.ChargingMode
		m["region"] = item.Region
		m["instance_type"] = item.InstanceType
		m["availability_zone"] = item.AvailabilityZone
		m["vpc_id"] = item.VpcId
		m["subnet_id"] = item.SubnetId
		m["status"] = item.Status

		m["account_id"] = strconv.FormatInt(item.AccountId, 10)

		// json列表型
		json.Unmarshal([]byte(item.PrivateIps), &prIps)
		json.Unmarshal([]byte(item.PublicIps), &puIps)

		// json map型
		json.Unmarshal([]byte(item.Tags), &tags)

		// g.p.a
		m["stree_group"] = item.StreeGroup
		m["stree_product"] = item.StreeProduct
		m["stree_app"] = item.StreeApp

		thisGPAS[fmt.Sprintf("%s.%s.%s", item.StreeGroup, item.StreeProduct, item.StreeApp)] = struct{}{}

		// 调用倒排索引库刷新索引
		thisH.GetOrCreateWithID(uint64(item.Id), item.Hash, mapTolsets(m))
		thisH.GetOrCreateWithID(uint64(item.Id), item.Hash, mapTolsets(tags))

		// 数组型
		for _, i := range prIps {
			mp := map[string]string{
				"private_ip": i,
			}
			thisH.GetOrCreateWithID(uint64(item.Id), item.Hash, mapTolsets(mp))
		}

		for _, i := range puIps {
			mp := map[string]string{
				"private_ip": i,
			}
			thisH.GetOrCreateWithID(uint64(item.Id), item.Hash, mapTolsets(mp))
		}
		for _, i := range prIps {
			mp := map[string]string{
				"public_ip": i,
			}
			thisH.GetOrCreateWithID(uint64(item.Id), item.Hash, mapTolsets(mp))
		}
	}

	hi.Ir.Reset(thisH)

	go func() {
		level.Info(hi.Logger).Log("msg", "FlushIndex.Add.GPA.To.PATH",
			"num", len(thisGPAS),
		)
		for node := range thisGPAS {
			inputs := common.NodeCommonReq{
				Node: node,
			}
			models.StreePathAddOne(&inputs, hi.Logger)
		}
	}()

}

func (hi *HostIndex) GetIndexReader() *mem.HeadIndexReader {
	return hi.Ir
}

func (hi *HostIndex) GetLogger() log.Logger {
	return hi.Logger
}

func mapTolsets(m map[string]string) labels.Labels {
	var lset labels.Labels
	for k, v := range m {
		l := labels.Label{
			Name:  k,
			Value: v,
		}
		lset = append(lset, l)
	}
	return lset
}
