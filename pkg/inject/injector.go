package inject

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/itchyny/gojq"

	"github.com/ShotaKitazawa/etcd-injector/pkg/etcdclient"
	"github.com/ShotaKitazawa/etcd-injector/pkg/rulesource"
)

type Injector struct {
	loggingEnable bool
}

func NewInjector(loggingEnable bool) *Injector {
	return &Injector{loggingEnable}
}

func (i *Injector) Inject(keyValues []etcdclient.KeyValue, rules []rulesource.Rule) (results []etcdclient.KeyValue, err error) {
	results = keyValues
	for _, rule := range rules {
		keyValues, results = results, []etcdclient.KeyValue{}
		for _, kv := range keyValues {
			result, err := i.injectOne(kv.Value, rule.JSONPath, rule.Repl)
			if err != nil {
				return nil, err
			}
			results = append(results, etcdclient.KeyValue{
				Key:   kv.Key,
				Value: result,
			})
			i.printf("key: %s, based_value: %s, replaced_value: %s\n", kv.Key, kv.Value, result)
		}
	}
	return results, nil
}

func (i *Injector) injectOne(input []byte, jsonPath string, replInterface interface{}) ([]byte, error) {
	repl, err := i.parseRepl(replInterface)
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

func (i *Injector) parseRepl(replInterface interface{}) (string, error) {
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

func (i *Injector) printf(format string, a ...interface{}) {
	if i.loggingEnable {
		fmt.Printf(format, a...)
	}
}
