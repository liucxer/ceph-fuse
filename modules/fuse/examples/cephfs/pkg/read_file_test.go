package pkg_test

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/liucxer/ceph-fuse/modules/fuse/examples/cephfs/pkg"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReadFile(t *testing.T) {
	bts, err := pkg.ReadFile("/cifs/111/222")
	require.NoError(t, err)
	spew.Dump(string(bts))
}
