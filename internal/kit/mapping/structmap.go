package mapping

import (
	"encoding/json"
	"fmt"
)

func StructToStringMap(v any) (map[string]string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var temp map[string]any
	if err := json.Unmarshal(b, &temp); err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for k, val := range temp {
		if val == nil {
			continue
		}
		if s, ok := val.(string); ok && s == "" {
			continue
		}
		result[k] = fmt.Sprintf("%v", val)
	}

	return result, nil
}
