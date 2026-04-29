package utils

import "encoding/json"

// ParseJSONArray 将 JSON 字符串解析为 []interface{}，空字符串返回空数组
func ParseJSONArray(raw string) []interface{} {
	if raw == "" || raw == "[]" || raw == "null" {
		return []interface{}{}
	}
	var arr []interface{}
	if err := json.Unmarshal([]byte(raw), &arr); err != nil {
		return []interface{}{}
	}
	return arr
}

// ParseJSONStringMap 将 JSON 字符串解析为 map[string]interface{}，空字符串返回空 map
func ParseJSONStringMap(raw string) map[string]interface{} {
	if raw == "" || raw == "{}" || raw == "null" {
		return map[string]interface{}{}
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return map[string]interface{}{}
	}
	return m
}

// ToJSON 将任意值序列化为 JSON 字符串，失败返回 "{}"
func ToJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(data)
}

// ToJSONIndent 将任意值序列化为带缩进的 JSON 字符串
func ToJSONIndent(v interface{}) string {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(data)
}
