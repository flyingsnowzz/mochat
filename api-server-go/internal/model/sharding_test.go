package model

import (
	"strings"
	"testing"
)

func TestGetUnionidMappingTableName(t *testing.T) {
	tests := []struct {
		name    string
		unionid string
	}{
		{
			name:    "empty unionid",
			unionid: "",
		},
		{
			name:    "unionid starting with a",
			unionid: "abc123",
		},
		{
			name:    "unionid starting with z",
			unionid: "zxc456",
		},
		{
			name:    "unionid starting with 1",
			unionid: "123abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetUnionidMappingTableName(tt.unionid)
			// 验证返回的表名格式是否正确
			if len(result) == 0 {
				t.Error("GetUnionidMappingTableName() returned empty string")
			}
			// 验证表名前缀
			if !strings.HasPrefix(result, "mc_work_unionid_external_userid_mapping_") {
				t.Errorf("GetUnionidMappingTableName() returned invalid table name format: %v", result)
			}
		})
	}
}
