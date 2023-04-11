// Hellofs implements a simple "hello world" file system.
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/liucxer/ceph-fuse/modules/fuse"
	"github.com/liucxer/ceph-fuse/modules/fuse/examples/cephfs/pkg"
	"github.com/liucxer/ceph-fuse/modules/fuse/fs"
	_ "github.com/liucxer/ceph-fuse/modules/fuse/fs/fstestutil"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"
)

// Dir implements both Node and Handle for the root directory.
type Dir struct {
	Path string
}

func (Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 1
	a.Mode = os.ModeDir | 0o555
	return nil
}

func (dir Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	return File{
		Path: dir.Path,
		Name: name,
	}, nil
}

var dirDirs = []fuse.Dirent{
	{Inode: 2, Name: "hello", Type: fuse.DT_File},
	{Inode: 3, Name: "hello_dir", Type: fuse.DT_Dir},
}

func (dir Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	dirs, err := pkg.ListDir(dir.Path)
	if err != nil {
		return nil, err
	}

	var res []fuse.Dirent
	for _, dir := range dirs {
		res = append(res, fuse.Dirent{
			Inode: uint64(dir.Inode()),
			Type:  fuse.DirentType(dir.DType()),
			Name:  dir.Name(),
		})
	}

	return res, nil
}

// File implements both Node and Handle for the hello file.
type File struct {
	Path string
	Name string
}

const greeting = "hello, world\n"

func (file File) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	dirs, err := pkg.ListDir(file.Path + "/" + file.Name)
	if err != nil {
		return nil, err
	}

	var res []fuse.Dirent
	for _, dir := range dirs {
		res = append(res, fuse.Dirent{
			Inode: uint64(dir.Inode()),
			Type:  fuse.DirentType(dir.DType()),
			Name:  dir.Name(),
		})
	}

	return res, nil
}

func (file File) Attr(ctx context.Context, a *fuse.Attr) error {
	cephAttr, err := pkg.GetFileAttr(file.Path, file.Name)
	if err != nil {
		return err
	}

	a.Inode = uint64(cephAttr.Inode)
	a.Size = cephAttr.Size
	a.Blocks = cephAttr.Blocks
	a.Atime = time.Unix(cephAttr.Atime.Sec, 0)
	a.Mtime = time.Unix(cephAttr.Mtime.Sec, 0)
	a.Ctime = time.Unix(cephAttr.Ctime.Sec, 0)

	a.Mode = os.FileMode(uint32(cephAttr.Mode) & (1<<10 - 1))
	logrus.Debugf("cephAttr.Mode: %b", cephAttr.Mode)
	logrus.Debugf("a.Mode: %b", a.Mode)
	logrus.Debugf("mode: %d", cephAttr.Mode>>14)
	if cephAttr.Mode>>14 == 1 {
		a.Mode = os.ModeDir | a.Mode
	}

	a.Nlink = cephAttr.Nlink
	a.Uid = cephAttr.Uid
	a.Gid = cephAttr.Gid
	a.Rdev = uint32(cephAttr.Rdev)
	a.BlockSize = cephAttr.Blksize
	/*
		none Valid time.Duration // how long Attr can be cached
		none  Flags     AttrFlags
	*/
	return nil
}

func (file File) Lookup(ctx context.Context, name string) (fs.Node, error) {
	return File{
		Path: file.Path + file.Name + "/",
		Name: name,
	}, nil
}

func (file File) ReadAll(ctx context.Context) ([]byte, error) {
	return pkg.ReadFile(file.Path + file.Name)
}

// FS implements the hello world file system.
type FS struct{}

func (FS) Root() (fs.Node, error) {
	return Dir{
		Path: "/",
	}, nil
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s MOUNTPOINT\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
		os.Exit(2)
	}
	mountpoint := flag.Arg(0)

	fuse.Unmount(mountpoint)
	c, err := fuse.Mount(
		mountpoint,
		fuse.FSName("helloworld"),
		fuse.Subtype("hellofs"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	err = fs.Serve(c, FS{})
	if err != nil {
		log.Fatal(err)
	}
}
