package main

import (
	"strings"

	"github.com/ShotaKitazawa/etcd-injector/pkg/etcdclient"
	"github.com/ShotaKitazawa/etcd-injector/pkg/inject"
	"github.com/ShotaKitazawa/etcd-injector/pkg/rulesource/file"
)

func Run(c config) error {
	// generate structs
	srcClient, err := etcdclient.New(etcdclient.Config{
		Endpoints:     c.SrcEndpoints,
		LoggingEnable: c.LoggingEnable,
	})
	if err != nil {
		return err
	}
	defer srcClient.Close()
	dstClient, err := etcdclient.New(etcdclient.Config{
		Endpoints:     c.DstEndpoints,
		LoggingEnable: c.LoggingEnable,
	})
	if err != nil {
		return err
	}
	defer dstClient.Close()
	injector := inject.NewInjector(c.LoggingEnable).WithIgnoreKeys(c.IgnoreKeys...)

	// load rules
	rules, err := file.GetRules(c.RulesFilepath)
	if err != nil {
		return err
	}

	// get values from source etcd
	keyValues, err := srcClient.LsRecursive(c.SrcDirectory)
	if err != nil {
		return err
	}

	// inject (or replace) values based rule
	results, err := injector.Inject(keyValues, rules)
	if err != nil {
		return err
	}

	// delete all keys injected previous if --delete=true
	if c.DeleteEnable {
		dstClient.DeleteRecursive(c.DstDirectory)
	}

	// replace keys & set values to destination etcd
	for _, kv := range results {
		kv.Key = strings.Replace(kv.Key, c.SrcDirectory, c.DstDirectory, 1)

		err := dstClient.Put(kv)
		if err != nil {
			return err
		}
	}

	return nil
}
