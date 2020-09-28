package etcdclient

import (
	"context"
	"crypto/tls"
	"time"

	"go.etcd.io/etcd/clientv3"
)

type Config struct {
	Endpoints     []string
	Username      string
	Password      string
	TLS           *tls.Config
	LoggingEnable bool
}

type Client struct {
	clientv3.Client
	loggingEnable bool
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
	return &Client{*cli, c.LoggingEnable}, nil
}

func (c *Client) LsRecursive(dirname string) ([]KeyValue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	opts := []clientv3.OpOption{clientv3.WithPrefix()}
	resp, err := c.Client.Get(ctx, dirname, opts...)
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if _, err := c.Client.Put(ctx, kv.Key, string(kv.Value)); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteRecursive(dirname string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	opts := []clientv3.OpOption{clientv3.WithPrefix()}
	_, err := c.Client.Delete(ctx, dirname, opts...)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Close() error {
	return c.Client.Close()
}
