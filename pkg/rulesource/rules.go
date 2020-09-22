package rulesource

type Rule struct {
	JSONPath string `json:"jsonpath"`
	Repl     string `json:"repl"`
}
