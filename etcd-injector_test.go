package main

import (
	"context"
	"os"
	"testing"

	"github.com/ShotaKitazawa/etcd-injector/pkg/mock"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/clientv3"
)

var etcdEndpointsForTest []string

func Test(t *testing.T) {
	tests := []struct {
		name   string
		config config
	}{
		{"normal_1", // 1keyのコピー
			config{
				SrcEndpoints:  etcdEndpointsForTest,
				DstEndpoints:  etcdEndpointsForTest,
				SrcDirectory:  "/test/src/1",
				DstDirectory:  "/test/dst/1",
				RulesFilepath: "example/rules.yaml",
			},
		},
		{"normal_2", // 複数keyのコピーが (srcにdirectoryを指定)
			config{
				SrcEndpoints:  etcdEndpointsForTest,
				DstEndpoints:  etcdEndpointsForTest,
				SrcDirectory:  "/test/src/dir",
				DstDirectory:  "/test/dst/dir",
				RulesFilepath: "example/rules.yaml",
			},
		},
		{"normal_3", // --ignore-keysオプションの有効化
			config{
				SrcEndpoints:  etcdEndpointsForTest,
				DstEndpoints:  etcdEndpointsForTest,
				SrcDirectory:  "/test/src/dir",
				DstDirectory:  "/test/dst/dir",
				RulesFilepath: "example/rules.yaml",
				IgnoreKeys:    []string{"/test/src/dir/1"},
			},
		},
		{"normal_4", // --deleteオプションの有効化
			config{
				SrcEndpoints:  etcdEndpointsForTest,
				DstEndpoints:  etcdEndpointsForTest,
				SrcDirectory:  "/test/src/dir",
				DstDirectory:  "/test/dst/dir",
				RulesFilepath: "example/rules.yaml",
				DeleteEnable:  true,
			},
		},
		{"normal_5", // --verboseオプションの有効化
			config{
				SrcEndpoints:  etcdEndpointsForTest,
				DstEndpoints:  etcdEndpointsForTest,
				SrcDirectory:  "/test/src/dir",
				DstDirectory:  "/test/dst/dir",
				RulesFilepath: "example/rules.yaml",
				LoggingEnable: true,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := Run(tt.config)

			assert.NoError(t, err)
		})
	}
}

func TestMain(m *testing.M) {
	var err error

	// run etcd server (retry: 10s * 6)
	etcdEndpointsForTest, err = mock.StartEtcdServer()
	if err != nil {
		panic(err)
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

	// exit
	os.Exit(status)
}
