package gusher

import (
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("gusher")
var format = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
)

/*
log.Info("info")
log.Notice("notice")
log.Warning("warning")
log.Error("err")
log.Critical("crit")
*/
