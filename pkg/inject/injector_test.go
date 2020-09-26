package inject

import (
	"os"
	"testing"

	"github.com/ShotaKitazawa/etcd-injector/pkg/etcdclient"
	"github.com/ShotaKitazawa/etcd-injector/pkg/rulesource"
	"github.com/stretchr/testify/assert"
)

func TestInjector(t *testing.T) {
	tests := []struct {
		name      string
		injector  *Injector
		keyValues []etcdclient.KeyValue
		rules     []rulesource.Rule
		results   []etcdclient.KeyValue
	}{
		{
			"normal_1", // 置換されることのテスト
			NewInjector(false),
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
			"normal_2", // 置換対象が複数存在する際にどちらも置換されることのテスト
			NewInjector(false),
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
			"normal_3", // 存在しないkeyに対するルールの場合、keyを挿入することのテスト
			NewInjector(false),
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
			"normal_4", // ルールが複数存在する際にどちらも動作することのテスト
			NewInjector(false),
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
			"normal_5", // ルールが存在しない場合そのままコピーすることのテスト
			NewInjector(false),
			[]etcdclient.KeyValue{
				{Key: "/test/src/1", Value: []byte(`{"fuga":"xxx","hoge":"ooo"}`)},
			},
			[]rulesource.Rule{},
			[]etcdclient.KeyValue{
				{Key: "/test/src/1", Value: []byte(`{"fuga":"xxx","hoge":"ooo"}`)},
			},
		},
		{
			"normal_6", // ignore が指定されたときに該当keyがコピー対象から除外されることのテスト
			NewInjector(false).WithIgnoreKeys(`/test/src/1`),
			[]etcdclient.KeyValue{
				{Key: "/test/src/1", Value: []byte(`{"fuga":"xxx","hoge":"ooo"}`)},
				{Key: "/test/src/2", Value: []byte(`{"fuga":"xxx","hoge":"ooo"}`)},
			},
			[]rulesource.Rule{},
			[]etcdclient.KeyValue{
				{Key: "/test/src/2", Value: []byte(`{"fuga":"xxx","hoge":"ooo"}`)},
			},
		},
		{
			"normal_7", // ignore が指定されたときに該当directory以下の複数keyがコピー対象から除外されることのテスト
			NewInjector(false).WithIgnoreKeys(`/test/src/dir`),
			[]etcdclient.KeyValue{
				{Key: "/test/src/1", Value: []byte(`{"fuga":"xxx","hoge":"ooo"}`)},
				{Key: "/test/src/dir/1", Value: []byte(`{"fuga":"xxx","hoge":"ooo"}`)},
				{Key: "/test/src/dir/2", Value: []byte(`{"fuga":"xxx","hoge":"ooo"}`)},
			},
			[]rulesource.Rule{},
			[]etcdclient.KeyValue{
				{Key: "/test/src/1", Value: []byte(`{"fuga":"xxx","hoge":"ooo"}`)},
			},
		},
		{
			"normal_8", // ignore が複数指定されたときに該当keyがコピー対象から除外されることのテスト
			NewInjector(false).WithIgnoreKeys(`/test/src/1`, `/test/src/2`),
			[]etcdclient.KeyValue{
				{Key: "/test/src/1", Value: []byte(`{"fuga":"xxx","hoge":"ooo"}`)},
				{Key: "/test/src/2", Value: []byte(`{"fuga":"xxx","hoge":"ooo"}`)},
				{Key: "/test/src/3", Value: []byte(`{"fuga":"xxx","hoge":"ooo"}`)},
			},
			[]rulesource.Rule{},
			[]etcdclient.KeyValue{
				{Key: "/test/src/3", Value: []byte(`{"fuga":"xxx","hoge":"ooo"}`)},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			results, err := tt.injector.Inject(tt.keyValues, tt.rules)

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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
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
