package models

import "regexp"

// 用户配置的日志策略，可以是agent 本地的yaml，也可以是通过接口过来的
type LogStrategy struct {
	ID         int64             `json:"id" yaml:"-"`
	MetricName string            `json:"metric_name" yaml:"metric_name"`
	MetricHelp string            `json:"metric_help" yaml:"metric_help"`
	FilePath   string            `json:"file_path" yaml:"file_path"`
	Pattern    string            `json:"pattern" yaml:"pattern"`
	Func       string            `json:"func" yaml:"func"`
	Tags       map[string]string `json:"tags" yaml:"tags"`
	Creator    string            `json:"creator"`
	// 上面是yaml或者前端的配置

	PatternReg *regexp.Regexp            `json:"-" yaml:"-"` // 主正则
	TagRegs    map[string]*regexp.Regexp `json:"-" yaml:"-"` // 标签的正则

}
