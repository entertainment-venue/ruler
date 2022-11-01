package ruler

import (
	"encoding/json"
	"encoding/xml"

	"gopkg.in/yaml.v2"
)

// DataDecoder 数据解析
type DataDecoder func(bts []byte, m *map[string]interface{}) error

func JsonDecoder(bts []byte, m *map[string]interface{}) error {
	return json.Unmarshal(bts, m)
}

func XmlDecoder(bts []byte, m *map[string]interface{}) error {
	return xml.Unmarshal(bts, m)
}

func YamlDecoder(bts []byte, m *map[string]interface{}) error {
	return yaml.Unmarshal(bts, m)
}
