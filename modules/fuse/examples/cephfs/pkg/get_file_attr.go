package pkg

import (
	"github.com/ceph/go-ceph/cephfs"
	"github.com/sirupsen/logrus"
)

func GetFileAttr(path string, name string) (*cephfs.CephStatx, error) {
	var err error

	mount, err := cephfs.CreateMount()
	if err != nil {
		logrus.Errorf("cephfs.CreateMount err:%v", err)
		return nil, err
	}
	defer func() {
		_ = mount.Release()
	}()

	err = mount.ReadConfigFile("/etc/ceph/ceph.conf")
	if err != nil {
		logrus.Errorf("mount.ReadConfigFile err:%v", err)
		return nil, err
	}

	err = mount.Mount()
	if err != nil {
		logrus.Errorf("mount.Mount err:%v", err)
		return nil, err
	}
	defer func() {
		_ = mount.Unmount()
	}()

	cephStat, err := mount.Statx(path+name, cephfs.StatxAllStats, cephfs.AtStatxDontSync)
	if err != nil {
		logrus.Errorf("mount.Statx err:%v", err)
		return nil, err
	}

	return cephStat, err
}
