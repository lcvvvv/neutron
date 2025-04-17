package templates

import (
	"encoding/json"
	"errors"

	"github.com/chainreactors/neutron/protocols"
	"github.com/chainreactors/neutron/protocols/executer"
	"github.com/chainreactors/neutron/protocols/http"
	"github.com/chainreactors/neutron/protocols/network"
	"gopkg.in/yaml.v3"
)

type StringOrSlice struct {
	value []string
}

func (m *StringOrSlice) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if err := node.Decode(&s); err == nil {
		m.value = []string{s}
		return nil
	}
	var ss []string
	if err := node.Decode(&ss); err == nil {
		m.value = ss
		return nil
	}
	return errors.New("failed to unmarshal StringOrSlice")
}

func (m *StringOrSlice) MarshalYAML() (interface{}, error) {
	return m.value, nil
}

// 自定义JSON序列化
func (m *StringOrSlice) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.value)
}

func (m *StringOrSlice) UnmarshalJSON(data []byte) error {
	var ss []string
	if err := json.Unmarshal(data, &ss); err == nil {
		m.value = ss
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		m.value = []string{s}
		return nil
	}
	return errors.New("failed to unmarshal StringOrSlice")
}

type Template struct {
	Id    string `json:"id" yaml:"id"`
	Opsec bool   `json:"opsec" yaml:"opsec"`
	Info  struct {
		Name        string                 `json:"name" yaml:"name"`
		Author      string                 `json:"author" yaml:"author"`
		Severity    string                 `json:"severity" yaml:"severity"`
		Description string                 `json:"description" yaml:"description"`
		Reference   *StringOrSlice         `json:"reference" yaml:"reference"`
		Vendor      string                 `json:"vendor" yaml:"vendor"`
		Tags        string                 `json:"tags" yaml:"tags"`
		Zombie      string                 `json:"zombie" yaml:"zombie"`
		Metadata    map[string]interface{} `json:"metadata" yaml:"metadata"`
	} `json:"info" yaml:"info"`

	Variables protocols.Variable `yaml:"variables,omitempty" json:"variables,omitempty"`

	RequestsHTTP []*http.Request `json:"http" yaml:"http"`
	// 适配部分较为老的PoC
	Requests        []*http.Request    `json:"requests" yaml:"requests"`
	RequestsNetwork []*network.Request `json:"network" yaml:"network"`

	// TotalRequests is the total number of requests for the template.
	TotalRequests int `yaml:"-" json:"-"`
	// Executor is the actual template executor for running template requests
	Executor *executer.Executer `yaml:"-" json:"-"`
}

func (t *Template) GetRequests() []*http.Request {
	if len(t.RequestsHTTP) > 0 {
		return t.RequestsHTTP
	}
	if len(t.Requests) > 0 {
		return t.Requests
	}
	return nil
}
