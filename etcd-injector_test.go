package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/ShotaKitazawa/etcd-injector/pkg/mock"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/clientv3"
)

var etcdEndpointsForTest []string

func Test(t *testing.T) {
	config := config{
		SrcEndpoints:  etcdEndpointsForTest,
		DstEndpoints:  etcdEndpointsForTest,
		RulesFilepath: "example/rules.yaml",
	}
	tests := []struct {
		name       string
		basePath   string
		targetPath string
	}{
		{"normal_1", "/test/src/1", "/test/dst/1"},
		{"normal_1", "/test/src/dir", "/test/dst/dir"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.SrcDirectory = tt.basePath
			config.DstDirectory = tt.targetPath
			err := Run(config)

			assert.NoError(t, err)
		})
	}
}

func TestMain(m *testing.M) {
	var err error

	// run etcd server (retry: 10s * 6)
	for i := 1; i <= 6; i++ {
		var errRecorver error
		if err := func() error {
			defer func() {
				var ok bool
				if errRecorver, ok = recover().(error); ok && errRecorver != nil {
					time.Sleep(time.Second * 10)
				}
			}()
			etcdEndpointsForTest, err = mock.StartEtcdServer()
			return err
		}(); err != nil || errRecorver != nil {
			continue
		}
		break
	}

	// put initialize value by go.etcd.io/etcd/clientv3
	if err := func() error {
		cli, err := clientv3.New(clientv3.Config{
			Endpoints: etcdEndpointsForTest,
		})
		if err != nil {
			return err
		}
		if _, err := cli.Put(context.Background(), "/test/src/1", `{"value":1}`); err != nil {
			return err
		}
		if _, err := cli.Put(context.Background(), "/test/src/dir/1", `{"value":"dir1"}`); err != nil {
			return err
		}
		if _, err := cli.Put(context.Background(), "/test/src/dir/2", `{"value":"dir2"}`); err != nil {
			return err
		}
		if _, err := cli.Put(context.Background(), "/test/src/dir/dir/1", `{"value":"dirdir1"}`); err != nil {
			return err
		}
		return nil
	}(); err != nil {
		panic(err)
	}

	// test
	status := m.Run()

	// for debug
	time.Sleep(time.Minute * 5)

	// exit
	os.Exit(status)
}
