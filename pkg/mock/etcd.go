package mock

import (
	"go.etcd.io/etcd/embed"
	"go.etcd.io/etcd/pkg/osutil"
)

var etcdEndpointsForMock = []string{"http://localhost:2379"}

func StartEtcdServer() ([]string, error) {
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
