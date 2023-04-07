// 核心配置

package config

import (
	"encoding/json"
	"strings"
)

type PointMap map[string]interface{} // 点位 Map，可转换为标准点位数据

// ToPoint 转换为标准点位数据
func (pm PointMap) ToPoint() Point {
	var p Point
	b, _ := json.Marshal(pm)
	_ = json.Unmarshal(b, &p)
	// 扩展参数
	p.Extends = make(map[string]interface{})
	for key, _ := range pm {
		if !strings.Contains("name,description,valueType,readWrite,defaultValue", key) {
			p.Extends[key] = pm[key]
		}
	}
	return p
}

// Config 配置
type Config struct {
	// 设备模型
	DeviceModels []DeviceModel `json:"deviceModels" validate:"required"`
	// 连接配置
	Connections map[string]interface{} `json:"connections" validate:"required"`
	// 协议名称（通过协议名称区分连接模式：客户端、服务端）
	ProtocolName string `json:"protocolName" validate:"required,oneof=http_server tcp_server modbus mqtt bacnet http_client"`
	// 配置唯一key，一般对应目录名称
	Key string `json:"-" validate:"-"`
	// 定时任务
	Tasks []TimerTask `json:"timerTasks" validate:"-"`
}

//------------------------------ 设备模型 ------------------------------

// DeviceModel 设备模型
type DeviceModel struct {
	// 模型名称
	Name string `json:"name" validate:"required"`
	// 云端模型 ID
	ModelID string `json:"modelId" validate:"required"`
	// 模型描述
	Description string `json:"description" validate:"required"`
	// 模型点位列表
	DevicePoints []PointMap `json:"devicePoints" validate:"required"`
	// 模型操作列表
	DeviceActions []DeviceAction `json:"deviceActions" validate:""`
	// 设备列表
	Devices []Device `json:"devices" validate:"required"`
}

// DeviceAction 设备操作
type DeviceAction struct {
	// 操作名称
	Name string `json:"name" validate:"required"`
	// 读写类型：RW、R、W
	ReadWrite string `json:"readWrite" validate:"required|oneof=R W RW"`
	// 资源操作列表
	ResourceOperations []ResourceOperation `json:"resourceOperations" validate:""`
}

// ResourceOperation 资源操作
type ResourceOperation struct {
	// 设备资源名称
	DeviceResource string `json:"deviceResource" validate:"-"`
	// 资源默认值
	DefaultValue string `json:"defaultValue" validate:"-"`
	// 资源扩展参数
	Mappings map[string]string `json:"mappings" validate:"-"`
}

// Point 点位数据
type Point struct {
	// 点位名称
	Name string `json:"name" validate:"required"`
	// 点位描述
	Description string `json:"description" validate:"required"`
	// 值类型
	ValueType string `json:"valueType" validate:"required,oneof=int float string"`
	// 读写模式
	ReadWrite string `json:"readWrite" validate:"required,oneof=int float string"`
	// 默认值
	DefaultValue string `json:"defaultValue" validate:"-"`
	// 实时上报开关
	RealReport bool `json:"realReport" validate:"required,boolean"`
	// 定时上报
	TimerReport string `json:"timerReport" validate:"required"`
	// 扩展参数
	Extends map[string]interface{} `json:"-" validate:"-"`
}

//------------------------------ 设备 ------------------------------

// Device 设备
type Device struct {
	// 设备名称
	Name string `json:"name" validate:"required"`
	// 模型名称
	ModelName string `json:"-" validate:"-"`
	// 设备描述
	Description string `json:"description" validate:"required"`
	// 连接 Key
	ConnectionKey string `json:"connectionKey" validate:"required"`
	// 协议参数
	Protocol map[string]string `json:"protocol" validate:"-"`
}

// TimerTask 定时任务
type TimerTask struct {
	// 间隔
	Interval string `json:"interval" validate:"required"`
	// 点位名称
	Type string `json:"type" validate:"required"`
	// 定时动作
	Action []Action `json:"action" validate:"required"`
}

type Action struct {
	DeviceNames []string `json:"devices"`
	Points      []string `json:"points"`
}
