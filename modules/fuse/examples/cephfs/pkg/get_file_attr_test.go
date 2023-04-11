package pkg_test

import (
	"bazil.org/fuse/examples/cephfs/pkg"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetFileAttr(t *testing.T) {
	res, err := pkg.GetFileAttr("/", "111")
	require.NoError(t, err)
	logrus.Infof("111: %b", res.Mode)

	res, err = pkg.GetFileAttr("/", "222")
	require.NoError(t, err)
	logrus.Infof("222: %b", res.Mode)

	res, err = pkg.GetFileAttr("/", "333")
	require.NoError(t, err)
	logrus.Infof("333: %b", res.Mode)

	res, err = pkg.GetFileAttr("/", "444")
	require.NoError(t, err)
	logrus.Infof("444: %b", res.Mode)

	res, err = pkg.GetFileAttr("/", "555")
	require.NoError(t, err)
	logrus.Infof("555: %b", res.Mode)

	res, err = pkg.GetFileAttr("/", "666")
	require.NoError(t, err)
	logrus.Infof("666: %b", res.Mode)

	res, err = pkg.GetFileAttr("/", "777")
	require.NoError(t, err)
	logrus.Infof("777: %b", res.Mode)

	res, err = pkg.GetFileAttr("/", "888")
	require.NoError(t, err)
	logrus.Infof("888: %b", res.Mode)

	res, err = pkg.GetFileAttr("/", "999")
	require.NoError(t, err)
	logrus.Infof("999: %b", res.Mode)
}

/*
time="2023-04-10T18:13:49+08:00" level=info msg="111: 1000000001001001"
time="2023-04-10T18:13:50+08:00" level=info msg="222: 1000000010010010"
time="2023-04-10T18:13:50+08:00" level=info msg="333: 1000000011011011"
time="2023-04-10T18:13:50+08:00" level=info msg="444: 1000000100100100"
time="2023-04-10T18:13:50+08:00" level=info msg="555: 1000000101101101"
time="2023-04-10T18:13:50+08:00" level=info msg="666: 1000000110110110"
time="2023-04-10T18:13:50+08:00" level=info msg="777: 1000000111111111"
time="2023-04-10T18:13:50+08:00" level=info msg="888:  100000111111111"
time="2023-04-10T18:13:50+08:00" level=info msg="999:  100000111111111"
*/
