package pkg

import (
	"github.com/ceph/go-ceph/cephfs"
	"github.com/sirupsen/logrus"
	"os"
)

func ReadFile(path string) ([]byte, error) {
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

	file, err := mount.Open(path, os.O_RDONLY, 0)
	if err != nil {
		logrus.Errorf("mount.Open err:%v", err)
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	bts := make([]byte, 100)
	_, err = file.Read(bts)
	if err != nil {
		logrus.Errorf("file.Read err:%v", err)
		return nil, err
	}

	return bts, err
}
