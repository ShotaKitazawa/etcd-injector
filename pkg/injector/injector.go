package injector

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/itchyny/gojq"

	"github.com/ShotaKitazawa/etcd-injector/pkg/etcdclient"
	"github.com/ShotaKitazawa/etcd-injector/pkg/rulesource"
)

func Inject(keyValues []etcdclient.KeyValue, rules []rulesource.Rule) (results []etcdclient.KeyValue, err error) {
	for _, rule := range rules {
		for _, kv := range keyValues {
			result, err := injectOne(kv.Value, rule.JSONPath, rule.Repl)
			if err != nil {
				return nil, err
			}
			results = append(results, etcdclient.KeyValue{
				Key:   kv.Key,
				Value: result,
			})
		}
	}
	return results, nil
}

func injectOne(input []byte, jsonPath string, replInterface interface{}) ([]byte, error) {
	repl, err := parseRepl(replInterface)
	if err != nil {
		return nil, err
	}

	var m []interface{}
	if err := json.Unmarshal([]byte("["+string(input)+"]"), &m); err != nil {
		return nil, err
	}

	q := fmt.Sprintf(`.[0]%s|=%s | .[0]`, jsonPath, repl)
	query, err := gojq.Parse(q)
	if err != nil {
		return nil, err
	}
	v, ok := query.Run(m).Next()
	if !ok {
		return nil, fmt.Errorf("gojq iterator cannot Next()")
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

func parseRepl(replInterface interface{}) (string, error) {
	switch t := replInterface.(type) {
	case string:
		return `"` + t + `"`, nil
	case int:
		return strconv.Itoa(t), nil
	case int64:
		return strconv.Itoa(int(t)), nil
	default:
		return "", fmt.Errorf(`repl is unsupported type: %v`, t)
	}
}
