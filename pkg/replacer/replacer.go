package replacer

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/itchyny/gojq"

	"github.com/ShotaKitazawa/etcd-replacer/pkg/etcdclient"
	"github.com/ShotaKitazawa/etcd-replacer/pkg/rulesource"
)

func Replace(keyValues []etcdclient.KeyValue, rules []rulesource.Rule) (result []etcdclient.KeyValue, err error) {
	for _, rule := range rules {
		for _, kv := range keyValues {
			replaced, err := replace(kv.Value, rule.JSONPath, rule.Repl)
			if err != nil {
				return nil, err
			}
			result = append(result, etcdclient.KeyValue{
				Key:   kv.Key,
				Value: replaced,
			})
		}
	}
	return result, nil
}

func replace(input []byte, jsonPath, repl string) ([]byte, error) {
	var m []interface{}
	if err := json.Unmarshal(wrap(input), &m); err != nil {
		return nil, err
	}
	query, err := gojq.Parse(fmt.Sprintf(`.[0]%s|=%s | .[0]`, jsonPath, repl))
	if err != nil {
		log.Fatalln(err)
	}
	v, ok := query.Run(m).Next()
	if !ok {
		return nil, fmt.Errorf("gojq iterator cannot )ext()")
	}
	if err, ok := v.(error); ok {
		return nil, err
	}

	output, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func wrap(input []byte) []byte {
	return []byte("[" + string(input) + "]")
}
