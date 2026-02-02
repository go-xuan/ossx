package ossx

import (
	"github.com/go-xuan/configx"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterClientBuilder("minio", MinioClientBuilder) // 注册minio客户端构建器
	Init()                                             // 初始化 oss
}

func Init() {
	logger := log.WithField("package", "ossx")
	if err := configx.LoadConfigurator(&Configs{}); err == nil && Initialized() {
		logger.Info("initialized success")
		return
	}
	if err := configx.LoadConfigurator(&Config{}); err == nil && Initialized() {
		logger.Info("initialized success")
		return
	}
	logger.Warn("initialized failed")
}
