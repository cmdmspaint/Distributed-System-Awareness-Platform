package web

import (
	"Distributed-System-Awareness-Platform/src/common"
	"Distributed-System-Awareness-Platform/src/models"
	"Distributed-System-Awareness-Platform/src/modules/server/memoryindex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/log"
	"math"
	"strconv"
)

func ResourceMount(c *gin.Context) {

	var inputs common.ResourceMountReq
	if err := c.BindJSON(&inputs); err != nil {
		common.JSONR(c, 400, err)
		return
	}
	logger := c.MustGet("logger").(log.Logger)

	// 校验 资源的名
	ok := models.CheckResources(inputs.ResourceType)
	if !ok {
		common.JSONR(c, 400, fmt.Errorf("resource_node_exist:%v", inputs.ResourceType))
		return
	}

	// 校验g.p.a是否存在
	qReq := &common.NodeCommonReq{
		Node:      inputs.TargetPath,
		QueryType: 4,
	}

	gpa := models.StreePathQuery(qReq, logger)
	if len(gpa) == 0 {
		common.JSONR(c, 400, fmt.Errorf("target_path_not_exist:%v", inputs.TargetPath))
		return
	}

	// 绑定的动作
	rowsAff, err := models.ResourceMount(&inputs, logger)
	if err != nil {
		common.JSONR(c, 500, err)
		return
	}

	common.JSONR(c, 200, fmt.Sprintf("rowsAff:%d", rowsAff))
	return

}

func ResourceQuery(c *gin.Context) {

	var inputs common.ResourceQueryReq
	if err := c.BindJSON(&inputs); err != nil {
		common.JSONR(c, 400, err)
		return
	}
	ok := memoryindex.JudgeResourceIndexExists(inputs.ResourceType)
	if !ok {
		common.JSONR(c, 400, fmt.Errorf("ResourceType_not_exists:%v", inputs.ResourceType))
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "100"))
	if err != nil {
		common.JSONR(c, 400, fmt.Errorf("invalid_page_size"))
		return
	}
	currentPage, err := strconv.Atoi(c.DefaultQuery("current_page", "1"))
	if err != nil {
		common.JSONR(c, 400, fmt.Errorf("invalid current_page"))
		return
	}

	offset := 0
	limit := 0
	limit = pageSize
	if currentPage > 1 {
		offset = (currentPage - 1) * limit
	}
	matchIds := memoryindex.GetMatchIdsByIndex(inputs)
	totalCount := len(matchIds)
	logger := c.MustGet("logger").(log.Logger)

	pageCount := int(math.Ceil(float64(totalCount) / float64(limit)))
	resp := common.QueryResponse{
		Code:        200,
		CurrentPage: currentPage,
		PageSize:    pageSize,
		PageCount:   pageCount,
		TotalCount:  totalCount,
	}
	res, err := models.ResourceQuery(inputs.ResourceType, matchIds, logger, limit, offset)
	if err != nil {
		resp.Code = 500
		resp.Result = err
	}
	resp.Result = res
	common.JSONR(c, resp)
}

func ResourceGroup(c *gin.Context) {
	resourceType := c.DefaultQuery("resource_type", common.RESOURCE_HOST)
	label := c.DefaultQuery("label", "region")

	ok := memoryindex.JudgeResourceIndexExists(resourceType)
	if !ok {
		common.JSONR(c, 400, fmt.Errorf("ResourceType_not_exists:%v", resourceType))
		return
	}
	_, ri := memoryindex.GetResourceIndexReader(resourceType)
	res := ri.GetIndexReader().GetGroupByLabel(label)
	common.JSONR(c, res)

}

func GetLabelDistribution(c *gin.Context) {

	var inputs common.ResourceQueryReq
	if err := c.BindJSON(&inputs); err != nil {
		common.JSONR(c, 400, err)
		return
	}
	ok, ri := memoryindex.GetResourceIndexReader(inputs.ResourceType)
	if !ok {
		common.JSONR(c, 400, fmt.Errorf("ResourceType_not_exists:%v", inputs.ResourceType))
		return
	}

	matchIds := memoryindex.GetMatchIdsByIndex(inputs)
	fmt.Println(inputs, matchIds)
	res := ri.GetIndexReader().GetGroupDistributionByLabel(inputs.TargetLabel, matchIds)
	common.JSONR(c, res)

}
