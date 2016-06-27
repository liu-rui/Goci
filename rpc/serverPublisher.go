package rpc

import (
	"errors"
	"path"
	"strings"
	"time"

	"github.com/liu-rui/goci/log"

	"github.com/samuel/go-zookeeper/zk"
)

var (
	repairInterval = time.Second * 10
)

type serverPublisher struct {
	servers  []string
	path     string
	dataPath string
	conn     *zk.Conn
}

func (sc *serverPublisher) Register() {
	if err := sc.repairProcess(); err != nil {
		panic(err)
	}
}

func (sc *serverPublisher) repairProcess() error {
	if sc.conn != nil && sc.conn.State() != zk.StateConnected {
		log.Debug("关闭连接")
		sc.closeConn()
	}

	if sc.conn == nil {
		sc.createConn()
	}

	if err := sc.mkDirs(sc.path); err != nil {
		log.Error("zookeeper创建目录时出现异常，异常:", err)
		return errors.New("zookeeper创建目录时出现异常")
	}
	if _, err := sc.conn.Create(sc.dataPath, []byte{}, zk.FlagEphemeral, zk.WorldACL(zk.PermAll)); err != nil {
		log.Error("zookeeper创建数据时出现异常，异常：", err)
		return errors.New("zookeeper创建数据时出现异常")
	}

	return sc.validateExistPath()
}

func (sc *serverPublisher) validateExistPath() error {
	exist, _, watch, err := sc.conn.ExistsW(sc.dataPath)

	if !exist {
		return errors.New("节点不存在")
	}

	if err != nil {
		return err
	}

	go func() {
		event := <-watch

		log.Debugf("节点发生更改 %v", event)
		switch event.State {
		case zk.StateDisconnected, zk.StateExpired, zk.StateUnknown, zk.StateAuthFailed:
			sc.startRepair()
			return
		}

		if err := sc.validateExistPath(); err != nil {
			sc.startRepair()
		}
	}()

	return nil
}

func (sc *serverPublisher) startRepair() {
	go func() {
		log.Infof("与RPC服务中心%s断开连接，启动修复程序", sc.servers)
		for {
			if err := sc.repairProcess(); err != nil {
				log.Error(err)
				time.Sleep(repairInterval)
			} else {
				break
			}
		}
		log.Infof("已与RPC服务中心%s重新建立连接，完成修复", sc.servers)
	}()
}

func (sc *serverPublisher) createConn() {
	c, _, err := zk.Connect(sc.servers, time.Second*10)

	if err != nil {
		panic(err)
	}

	max := 5

	for c.State() != zk.StateConnected && max > 1 {
		time.Sleep(100 * time.Millisecond)
		max--
	}
	sc.conn = c
}

func (sc *serverPublisher) closeConn() {
	if sc.conn == nil {
		return
	}

	sc.conn.Close()
	sc.conn = nil
}

func (sc *serverPublisher) mkDirs(dir string) error {
	fpath := "/"

	for _, part := range strings.Split(dir, "/") {
		fpath = path.Join(fpath+"/", part)

		if _, err := sc.conn.Create(fpath, []byte{}, 0, zk.WorldACL(zk.PermAll)); err != nil {
			if err == zk.ErrNodeExists {
				continue
			} else {
				return err
			}
		}
	}
	return nil
}

func (sc *serverPublisher) Close() {
	if sc.conn != nil {
		sc.conn.Close()
	}
}

func newServerPublisher(servers []string, dir, data string) *serverPublisher {
	return &serverPublisher{servers: servers, path: dir, dataPath: path.Join(dir, data)}
}
