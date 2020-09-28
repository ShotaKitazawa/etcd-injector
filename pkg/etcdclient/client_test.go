package etcdclient

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/clientv3"

	"github.com/ShotaKitazawa/etcd-injector/pkg/mock"
)

var etcdEndpointsForTest []string

func TestLsRecursive(t *testing.T) {
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
				Value: []byte("hogehoge"),
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
			Value: []byte("overwrited"),
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

func TestDeleteRecursive(t *testing.T) {
	client, err := New(Config{
		Endpoints: etcdEndpointsForTest,
	})
	if err != nil {
		panic(err)
	}
	tests := []struct {
		name    string
		dirname string
		keyOne  string
	}{
		{"normal_1", "/test/ls/1", "/test/ls/1"},
		{"normal_2", "/test/ls/dir", "/test/ls/dir/1"},
		{"empty_1", "/test/ls/empty", "/test/ls/empty"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.DeleteRecursive(tt.dirname)
			assert.NoError(t, err)

			resp, err := client.Client.Get(context.Background(), tt.keyOne)
			if err != nil {
				panic(err)
			}
			assert.Empty(t, resp.Kvs)
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
		for _, val := range []struct {
			key   string
			value string
		}{
			{"/test/ls/1", "hogehoge"},
			{"/test/ls/dir/1", "hogehoge"},
			{"/test/ls/dir/2", "hogehoge"},
			{"/test/del/1", "hogehoge"},
			{"/test/del/dir/1", "hogehoge"},
			{"/test/del/dir/2", "hogehoge"},
		} {
			if _, err := cli.Put(context.Background(), val.key, val.value); err != nil {
				return err
			}
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
