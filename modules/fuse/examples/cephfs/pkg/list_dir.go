package pkg

import (
	"github.com/ceph/go-ceph/cephfs"
	"github.com/sirupsen/logrus"
)

func ListDir(path string) ([]cephfs.DirEntryPlus, error) {
	var (
		err error
		res []cephfs.DirEntryPlus
	)

	mount, err := cephfs.CreateMount()
	if err != nil {
		logrus.Errorf("cephfs.CreateMount err:%v", err)
		return res, err
	}
	defer func() {
		_ = mount.Release()
	}()

	err = mount.ReadConfigFile("/etc/ceph/ceph.conf")
	if err != nil {
		logrus.Errorf("mount.ReadConfigFile err:%v", err)
		return res, err
	}

	err = mount.Mount()
	if err != nil {
		logrus.Errorf("mount.Mount err:%v", err)
		return res, err
	}
	defer func() {
		_ = mount.Unmount()
	}()

	dir, err := mount.OpenDir(path)
	if err != nil {
		logrus.Errorf("mount.OpenDir err:%v", err)
		return res, err
	}
	for {
		entry, err := dir.ReadDirPlus(cephfs.StatxBasicStats, cephfs.AtSymlinkNofollow)
		if err != nil {
			logrus.Errorf("dir.ReadDirPlus err:%v", err)
			return res, err
		}

		if entry == nil {
			break
		}
		res = append(res, *entry)
		logrus.Infof("entry name:%s, inode:%d, type:%v", entry.Name(), entry.Inode(), entry.DType())
	}

	return res, err
}
