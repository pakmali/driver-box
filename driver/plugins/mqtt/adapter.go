package mqtt

import (
	"driver-box/core/contracts"
	"driver-box/core/helper"
	"driver-box/driver/common"
	"encoding/json"
	lua "github.com/yuin/gopher-lua"
)

// Adapter 协议适配器
type adapter struct {
	scriptDir string // 脚本目录名称
	ls        *lua.LState
}

func (a *adapter) Encode(deviceName string, mode contracts.EncodeMode, value contracts.PointData) (res interface{}, err error) {
	if mode == contracts.WriteMode {
		tmp, _ := json.Marshal(value)
		return helper.CallLuaEncodeConverter(a.ls, deviceName, string(tmp))
	}
	return nil, common.NotSupportEncode
}

// Decode 解析数据
func (a *adapter) Decode(raw interface{}) (res []contracts.DeviceData, err error) {
	return helper.CallLuaConverter(a.ls, "decode", raw)
}
