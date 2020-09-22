package etcdclient

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/clientv3"

	"github.com/ShotaKitazawa/etcd-replacer/pkg/mock"
)

var etcdEndpointsForTest []string

func TestLsRecursive(t *testing.T) {
	t.Parallel()
	client, err := New(Config{
		Endpoints: etcdEndpointsForTest,
	})
	if err != nil {
		panic(err)
	}
	tests := []struct {
		name    string
		dirname string
		kvs     []KeyValue
	}{
		{"normal_1", "/test/ls/1", []KeyValue{
			{
				Key:   "/test/ls/1",
				Value: []byte("hogehoge"),
			},
		}},
		{"normal_2", "/test/ls/dir", []KeyValue{
			{
				Key:   "/test/ls/dir/1",
				Value: []byte("hogehoge"),
			},
			{
				Key:   "/test/ls/dir/2",
				Value: []byte("fugafuga"),
			},
		}},
		{"empty_1", "/test/ls/empty", []KeyValue{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kvs, err := client.LsRecursive(tt.dirname)

			assert.NoError(t, err)
			assert.Equal(t, tt.kvs, kvs)
		})
	}
}

func TestPut(t *testing.T) {
	t.Parallel()
	client, err := New(Config{
		Endpoints: etcdEndpointsForTest,
	})
	if err != nil {
		panic(err)
	}
	tests := []struct {
		name string
		kv   KeyValue
	}{
		{"normal_1", KeyValue{
			Key:   "/test/put/1",
			Value: []byte("hogehoge"),
		}},
		{"normal_2", KeyValue{
			Key:   "/test/put/1",
			Value: []byte("overwrited"),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Put(tt.kv)
			assert.NoError(t, err)

			resp, err := client.Client.Get(context.Background(), tt.kv.Key)
			if err != nil {
				panic(err)
			}
			assert.Equal(t, tt.kv.Value, resp.Kvs[0].Value)
		})
	}
}

func TestMain(m *testing.M) {
	var err error

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
		if _, err := cli.Put(context.Background(), "/test/ls/1", "hogehoge"); err != nil {
			return err
		}
		if _, err := cli.Put(context.Background(), "/test/ls/dir/1", "hogehoge"); err != nil {
			return err
		}
		if _, err := cli.Put(context.Background(), "/test/ls/dir/2", "fugafuga"); err != nil {
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
