package main

import (
	"github.com/ceph/go-ceph/cephfs"
	"github.com/davecgh/go-spew/spew"
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

func main() {
	res, err := ListDir("/")
	if err != nil {
		logrus.Errorf("ListDir err:%v", err)
		return
	}
	spew.Dump(res)

	cephStat, err := GetFileAttr("/", "ftp")
	if err != nil {
		logrus.Errorf("GetFileAttr err:%v", err)
		return
	}
	spew.Dump(cephStat)

	return

}
