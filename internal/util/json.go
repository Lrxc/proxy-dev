package util

import "encoding/json"

func PrettyJSON(raw any) string {
	if s, ok := raw.(string); ok {
		json.Unmarshal([]byte(s), &raw)
	}

	pretty, _ := json.MarshalIndent(raw, "", "  ")
	return string(pretty)
}
