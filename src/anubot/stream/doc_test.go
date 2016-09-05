package stream

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/fluffle/goirc/logging/golog"
)

func init() {
	golog.Init()
	if !testing.Verbose() {
		log.SetOutput(ioutil.Discard)
	}
}
