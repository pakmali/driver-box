package httpserver

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type connectorConfig struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
}

type connector struct {
	plugin *Plugin
	server *http.Server
}

// Release 释放资源
func (c *connector) Release() (err error) {
	return c.server.Shutdown(context.Background())
}

// Send 被动接收数据模式，无需实现
func (c *connector) Send(raw interface{}) (err error) {
	return nil
}

// startServer 启动服务
func (c *connector) startServer(opts connectorConfig) {
	// 启动服务
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	app := gin.Default()
	// 通用路由
	app.NoRoute(func(ctx *gin.Context) {
		// 取 body
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			c.plugin.logger.Error("http request read body error", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    -1,
				"message": err.Error(),
			})
			return
		}
		// 重组协议数据
		data := protoData{
			Path:   ctx.Request.URL.Path,
			Method: ctx.Request.Method,
			Body:   string(body),
		}
		// 调用回调函数
		if _, err = c.plugin.callback(c.plugin, data.ToJSON()); err != nil {
			c.plugin.logger.Error("http_server callback error", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    -1,
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "ok",
		})
		return
	})

	addr := fmt.Sprintf("%s:%d", opts.Host, opts.Port)
	c.server = &http.Server{
		Addr:    addr,
		Handler: app,
	}

	go func(addr string) {
		if err := c.server.ListenAndServe(); err != nil {
			c.plugin.logger.Error("start http server error", zap.Error(err))
		}
	}(addr)
}
