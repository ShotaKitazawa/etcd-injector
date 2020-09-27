package rulesource

type Rule struct {
	JSONPath string      `json:"jsonpath"`
	Repl     interface{} `json:"repl"`
}
