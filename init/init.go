package init

import (
	"github/go-xuan/ossx"

	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/utilx/errorx"
)

func init() {
	errorx.Panic(Init())
}

func Init() error {
	var err error
	if err = configx.LoadConfigurator(&ossx.Configs{}); err == nil && ossx.Initialized() {
		return nil
	} else if err = configx.LoadConfigurator(&ossx.Config{}); err == nil && ossx.Initialized() {
		return nil
	}
	return errorx.Wrap(err, "init oss failed")
}
