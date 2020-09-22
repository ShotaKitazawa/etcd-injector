package file

import (
	"io/ioutil"

	"github.com/ShotaKitazawa/etcd-injector/pkg/rulesource"
	"gopkg.in/yaml.v2"
)

func GetRules(filepath string) ([]rulesource.Rule, error) {
	// open policy file
	buf, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	// if file is empty
	if len(buf) == 0 {
		return nil, nil
	}

	// unmarshal
	var rules []rulesource.Rule
	if err = yaml.Unmarshal(buf, &rules); err != nil {
		return nil, err
	}

	return rules, nil
}
