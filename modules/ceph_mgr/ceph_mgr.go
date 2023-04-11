package ceph_mgr

import (
	"github.com/liucxer/ceph-fuse/modules/ceph/cephfs"
	"github.com/sirupsen/logrus"
	"os"
)

type CephMgr struct {
	mountInfo *cephfs.MountInfo
}

func NewCephMgr() (*CephMgr, error) {
	var (
		cephMgr   CephMgr
		err       error
		mountInfo *cephfs.MountInfo
	)
	mountInfo, err = cephfs.CreateMount()
	if err != nil {
		logrus.Errorf("cephfs.CreateMount err:%v", err)
		return nil, err
	}

	err = mountInfo.ReadConfigFile("/etc/ceph/ceph.conf")
	if err != nil {
		logrus.Errorf("mountInfo.ReadConfigFile err:%v", err)
		return nil, err
	}

	err = mountInfo.Mount()
	if err != nil {
		logrus.Errorf("mountInfo.Mount err:%v", err)
		return nil, err
	}

	cephMgr.mountInfo = mountInfo
	return &cephMgr, err
}

func (mgr *CephMgr) Close() error {
	err := mgr.mountInfo.Unmount()
	if err != nil {
		logrus.Errorf("mgr.mountInfo.Unmount err:%v", err)
		return err
	}

	err = mgr.mountInfo.Release()
	if err != nil {
		logrus.Errorf("mgr.mountInfo.Release err:%v", err)
		return err
	}
	return nil
}

func (mgr *CephMgr) ListDir(path string) ([]cephfs.DirEntryPlus, error) {
	var (
		err error
		res []cephfs.DirEntryPlus
	)

	dir, err := mgr.mountInfo.OpenDir(path)
	if err != nil {
		logrus.Errorf("mgr.mountInfo.OpenDir err:%v,path:%s", err, path)
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

func (mgr *CephMgr) GetFileAttr(path string, name string) (*cephfs.CephStatx, error) {
	var err error

	cephStat, err := mgr.mountInfo.Statx(path+name, cephfs.StatxAllStats, cephfs.AtStatxDontSync)
	if err != nil {
		logrus.Errorf("mgr.mountInfo.Statx err:%v", err)
		return nil, err
	}

	return cephStat, err
}

func (mgr *CephMgr) ReadFile(path string) ([]byte, error) {
	var err error

	file, err := mgr.mountInfo.Open(path, os.O_RDONLY, 0)
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
		logrus.Errorf("file.Read err:%v,path:%s", err, path)
		return nil, err
	}

	return bts, err
}
