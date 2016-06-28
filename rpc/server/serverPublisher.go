package server

import (
	"errors"
	"path"
	"strings"

	"github.com/liu-rui/goci/log"
	"github.com/liu-rui/goci/rpc"
	"github.com/samuel/go-zookeeper/zk"
)

var (
	errorMkDir          = errors.New("zookeeper创建目录时出现异常")
	errorCreateDataPath = errors.New("zookeeper创建数据时出现异常")
	errorNodeNotExist   = errors.New("节点不存在")
)

type serverPublisher struct {
	rpc.ZKObject
	path     string
	dataPath string
}

func (sc *serverPublisher) Process() error {
	if err := sc.mkDirs(sc.path); err != nil {
		log.Error("zookeeper创建目录时出现异常，异常:", err)
		return errorMkDir
	}
	if _, err := sc.Conn.Create(sc.dataPath, []byte{}, zk.FlagEphemeral, zk.WorldACL(zk.PermAll)); err != nil {
		log.Error("zookeeper创建数据时出现异常，异常：", err)
		return errorCreateDataPath
	}
	return sc.ProcessWithEvent()
}

func (sc *serverPublisher) ProcessWithEvent() error {
	exist, _, watch, err := sc.Conn.ExistsW(sc.dataPath)

	if !exist {
		return errorNodeNotExist
	}

	if err != nil {
		return err
	}

	sc.Listen(watch)
	return nil
}

func (sc *serverPublisher) mkDirs(dir string) error {
	fpath := "/"

	for _, part := range strings.Split(dir, "/") {
		fpath = path.Join(fpath+"/", part)

		if _, err := sc.Conn.Create(fpath, []byte{}, 0, zk.WorldACL(zk.PermAll)); err != nil {
			if err == zk.ErrNodeExists {
				continue
			} else {
				return err
			}
		}
	}
	return nil
}

func newServerPublisher(servers []string, dir, data string) *serverPublisher {
	ret := &serverPublisher{ZKObject: rpc.ZKObject{Servers: servers}, path: dir, dataPath: path.Join(dir, data)}
	ret.Processer = ret
	return ret
}
