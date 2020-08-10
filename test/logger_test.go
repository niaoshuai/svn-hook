package test

import (
	"svn-hook/pkg/log"
	"testing"
)

func TestLogger(t *testing.T) {
	log.InitLog("tests.log")
	log.Info("testInfo")
	//log.Fatal(nil)
}
