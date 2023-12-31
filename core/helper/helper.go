// 核心工具助手文件

package helper

import (
	"driver-box/core/contracts"
	"driver-box/core/helper/shadow"
	"encoding/json"
	sdkModels "github.com/edgexfoundry/device-sdk-go/v2/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
)

var MessageBus chan<- *sdkModels.AsyncValues // 消息总线

var RunningPlugin contracts.Plugin // todo 待删除

var DeviceShadow shadow.DeviceShadow // 本地设备影子

var PluginCacheMap = &sync.Map{} // 插件通用缓存

// Map2Struct map 转 struct，用于解析连接器配置
// m：map[string]interface
// v：&struct{}
func Map2Struct(m interface{}, v interface{}) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

// PointValueType2EdgeX 点位值类型转换为 EdgeX 数据类型
// int => Int64、float => Float64、string => String
func PointValueType2EdgeX(valueType string) string {
	switch strings.ToLower(valueType) {
	case "int":
		return common.ValueTypeInt64
	case "float":
		return common.ValueTypeFloat64
	case "string":
		return common.ValueTypeString
	default:
		return valueType
	}
}

// GetChildDir 获取指定路径下所有子目录
func GetChildDir(path string) (list []string, err error) {
	err = filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			list = append(list, path)
		}
		return nil
	})
	return
}
