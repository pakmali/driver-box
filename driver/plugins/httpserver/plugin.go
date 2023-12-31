package httpserver

import (
	"driver-box/core/config"
	"driver-box/core/contracts"
	"driver-box/core/helper"
	"driver-box/driver/common"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"sync"
)

type Plugin struct {
	logger   *zap.Logger                // 日志记录器
	config   config.Config              // 核心配置
	callback contracts.OnReceiveHandler // 回调函数
	adapter  contracts.ProtocolAdapter  // 协议适配器
	connPool []*connector               // 连接器
	ls       *lua.LState                // lua 虚拟机
}

func (p *Plugin) Initialize(logger *zap.Logger, c config.Config, handler contracts.OnReceiveHandler, ls *lua.LState) (err error) {
	p.logger = logger
	p.config = c
	p.callback = handler
	p.ls = ls

	// 初始化协议适配器
	p.adapter = &adapter{
		scriptDir: c.Key,
		ls:        ls,
		lock:      &sync.Mutex{},
	}

	// 初始化连接池
	if err = p.initConnPool(); err != nil {
		return
	}

	return nil
}

func (p *Plugin) ProtocolAdapter() contracts.ProtocolAdapter {
	return p.adapter
}

// Connector 此协议不支持获取连接器
func (p *Plugin) Connector(deviceName, pointName string) (connector contracts.Connector, err error) {
	return nil, common.NotSupportGetConnector
}

func (p *Plugin) Destroy() error {
	defer p.ls.Close()
	if len(p.connPool) > 0 {
		for i, _ := range p.connPool {
			if err := p.connPool[i].Release(); err != nil {
				return err
			}
		}
	}
	return nil
}

// initConnPool 初始化连接池
func (p *Plugin) initConnPool() (err error) {
	for key, _ := range p.config.Connections {
		var c connectorConfig
		if err = helper.Map2Struct(p.config.Connections[key], &c); err != nil {
			return
		}
		conn := &connector{
			plugin: p,
		}
		conn.startServer(c)
		p.connPool = append(p.connPool, conn)
	}
	return
}
