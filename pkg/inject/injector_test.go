package inject

import (
	"os"
	"testing"

	"github.com/ShotaKitazawa/etcd-injector/pkg/etcdclient"
	"github.com/ShotaKitazawa/etcd-injector/pkg/rulesource"
	"github.com/stretchr/testify/assert"
)

func TestInjector(t *testing.T) {
	i := NewInjector(false)
	tests := []struct {
		name      string
		keyValues []etcdclient.KeyValue
		rules     []rulesource.Rule
		results   []etcdclient.KeyValue
	}{
		{
			"normal_1",
			[]etcdclient.KeyValue{
				{Key: "/test/src/1", Value: []byte(`{"hoge":"ooo"}`)},
			},
			[]rulesource.Rule{
				{JSONPath: ".hoge", Repl: "replaced"},
			},
			[]etcdclient.KeyValue{
				{Key: "/test/src/1", Value: []byte(`{"hoge":"replaced"}`)},
			},
		},
		{
			"normal_2",
			[]etcdclient.KeyValue{
				{Key: "/test/src/1", Value: []byte(`{"hoge":"ooo"}`)},
				{Key: "/test/src/2", Value: []byte(`{"hoge":"xxx"}`)},
			},
			[]rulesource.Rule{
				{JSONPath: ".hoge", Repl: "replaced"},
			},
			[]etcdclient.KeyValue{
				{Key: "/test/src/1", Value: []byte(`{"hoge":"replaced"}`)},
				{Key: "/test/src/2", Value: []byte(`{"hoge":"replaced"}`)},
			},
		},
		{
			"normal_2",
			[]etcdclient.KeyValue{
				{Key: "/test/src/1", Value: []byte(`{"hoge":"ooo"}`)},
			},
			[]rulesource.Rule{
				{JSONPath: ".injected", Repl: "value"},
			},
			[]etcdclient.KeyValue{
				{Key: "/test/src/1", Value: []byte(`{"hoge":"ooo","injected":"value"}`)},
			},
		},
		{
			"normal_3",
			[]etcdclient.KeyValue{
				{Key: "/test/src/1", Value: []byte(`{"fuga":"xxx","hoge":"ooo"}`)},
			},
			[]rulesource.Rule{
				{JSONPath: ".hoge", Repl: "replaced_hoge"},
				{JSONPath: ".fuga", Repl: "replaced_fuga"},
			},
			[]etcdclient.KeyValue{
				{Key: "/test/src/1", Value: []byte(`{"fuga":"replaced_fuga","hoge":"replaced_hoge"}`)},
			},
		},
		{
			"normal_3",
			[]etcdclient.KeyValue{
				{Key: "/test/src/1", Value: []byte(`{"fuga":"xxx","hoge":"ooo"}`)},
			},
			[]rulesource.Rule{},
			[]etcdclient.KeyValue{
				{Key: "/test/src/1", Value: []byte(`{"fuga":"xxx","hoge":"ooo"}`)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := i.Inject(tt.keyValues, tt.rules)

			assert.NoError(t, err)
			assert.Equal(t, tt.results, results)
		})
	}

}

func Test_injectOne(t *testing.T) {
	i := NewInjector(false)
	tests := []struct {
		name     string
		input    []byte
		jsonPath string
		repl     interface{}
		output   []byte
	}{
		{"normal_1", []byte(`{"key":"value"}`), ".key", "replaced", []byte(`{"key":"replaced"}`)},
		{"normal_2", []byte(`{"key":"value"}`), ".key", 1, []byte(`{"key":1}`)},
		{"normal_3", []byte(`[{"key":"value"}]`), ".[0].key", "replaced", []byte(`[{"key":"replaced"}]`)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := i.injectOne(tt.input, tt.jsonPath, tt.repl)

			assert.NoError(t, err)
			assert.Equal(t, tt.output, output)
		})
	}
}

func TestMain(m *testing.M) {
	// test
	status := m.Run()

	// exit
	os.Exit(status)
}
