package mock

import (
	"time"

	"go.etcd.io/etcd/embed"
	"go.etcd.io/etcd/pkg/osutil"
)

const (
	retryCount = 6
)

var (
	etcdEndpointsForMock = []string{"http://localhost:2379"}
)

func StartEtcdServer() (etcdEndpoints []string, err error) {
	for i := 1; i <= retryCount; i++ {
		var errRecorver error
		if err := func() error {
			defer func() {
				var ok bool
				if errRecorver, ok = recover().(error); ok && errRecorver != nil {
					time.Sleep(time.Second * 10)
				}
			}()
			etcdEndpoints, err = startEtcdServer()
			return err
		}(); err != nil || errRecorver != nil {
			continue
		}
		break
	}
	return etcdEndpoints, nil
}

func startEtcdServer() ([]string, error) {
	// run etcd server
	// TODO: undisplay log
	if err := func() error {
		cfg := embed.NewConfig()
		cfg.Dir = "/tmp/.etcd-injector.test.etcd"
		e, err := embed.StartEtcd(cfg)
		if err != nil {
			panic(err)
		}
		osutil.RegisterInterruptHandler(e.Close)
		select {
		case <-e.Server.ReadyNotify(): // wait for e.Server to join the cluster
			return nil
		case <-e.Server.StopNotify(): // publish aborted from 'ErrStopped'
			return <-e.Err()
		}
	}(); err != nil {
		return nil, err
	}

	return etcdEndpointsForMock, nil
}
