package etcdclient

import (
	"context"
	"crypto/tls"
	"time"

	"go.etcd.io/etcd/clientv3"
)

type Config struct {
	Endpoints []string
	Username  string
	Password  string
	TLS       *tls.Config
}

type Client struct {
	*clientv3.Client
}

type KeyValue struct {
	Key   string
	Value []byte
}

func New(c Config) (*Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   c.Endpoints,
		Username:    c.Username,
		Password:    c.Password,
		TLS:         c.TLS,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return &Client{cli}, nil
}

func (c *Client) LsRecursive(dirname string) ([]KeyValue, error) {
	opts := []clientv3.OpOption{clientv3.WithPrefix()}
	resp, err := c.Client.Get(context.Background(), dirname, opts...)
	if err != nil {
		return nil, err
	}

	result := []KeyValue{}
	for _, respKv := range resp.Kvs {
		result = append(result, KeyValue{
			Key:   string(respKv.Key),
			Value: respKv.Value,
		})
	}
	return result, nil
}

func (c *Client) Put(kv KeyValue) error {
	if _, err := c.Client.Put(context.Background(), kv.Key, string(kv.Value)); err != nil {
		return err
	}
	return nil
}

func (c *Client) Close() error {
	return c.Client.Close()
}
