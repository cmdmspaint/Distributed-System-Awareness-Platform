package web

import (
	"Distributed-System-Awareness-Platform/src/common"
	"Distributed-System-Awareness-Platform/src/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/log"
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
