package bench_test

import (
	"flag"
	"os"
	"testing"

	"github.com/liucxer/ceph-fuse/modules/fuse/fs/fstestutil/spawntest"
)

var helpers spawntest.Registry

func TestMain(m *testing.M) {
	helpers.AddFlag(flag.CommandLine)
	flag.Parse()
	helpers.RunIfNeeded()
	os.Exit(m.Run())
}
