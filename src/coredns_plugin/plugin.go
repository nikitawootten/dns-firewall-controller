package squawker

import (
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
)

var log = clog.NewWithPlugin("squawker")

type Squawker struct {
	Next plugin.Handler
}
