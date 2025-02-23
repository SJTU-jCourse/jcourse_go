package converter

import (
	"testing"

	"jcourse_go/model/po"
	"jcourse_go/model/types"
)

func TestGetSettingFromPO(t *testing.T) {
	tests := []struct {
		name        string
		input       po.SettingPO
		expectedKey string
		expectedVal any
		expectedErr bool
	}{
		{
			name: "Valid String Setting",
			input: po.SettingPO{
				Key:   "site_name",
				Value: "JCourse",
				Type:  string(types.SettingTypeString),
			},
			expectedKey: "site_name",
			expectedVal: "JCourse",
			expectedErr: false,
		},
		{
			name: "Valid Int Setting",
			input: po.SettingPO{
				Key:   "max_users",
				Value: "1000",
				Type:  string(types.SettingTypeInt),
			},
			expectedKey: "max_users",
			expectedVal: int64(1000),
			expectedErr: false,
		},
		{
			name: "Valid Bool Setting",
			input: po.SettingPO{
				Key:   "enable_feature",
				Value: "true",
				Type:  string(types.SettingTypeBool),
			},
			expectedKey: "enable_feature",
			expectedVal: true,
			expectedErr: false,
		},
		{
			name: "Invalid Int Value",
			input: po.SettingPO{
				Key:   "max_users",
				Value: "invalid_int",
				Type:  string(types.SettingTypeInt),
			},
			expectedKey: "max_users",
			expectedVal: nil,
			expectedErr: true,
		},
		{
			name: "Unknown Setting Type",
			input: po.SettingPO{
				Key:   "unknown",
				Value: "value",
				Type:  "unknown_type",
			},
			expectedKey: "unknown",
			expectedVal: nil,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setting, err := GetSettingFromPO(tt.input)
			if tt.expectedErr {
				if err == nil {
					t.Errorf("预期错误，但得到正常结果")
				}
				return
			}
			if err != nil {
				t.Errorf("未预期的错误: %v", err)
				return
			}
			if setting.GetKey() != tt.expectedKey {
				t.Errorf("期望Key: %s, 实际Key: %s", tt.expectedKey, setting.GetKey())
			}
			if setting.GetValue() != tt.expectedVal {
				t.Errorf("期望Value: %v, 实际Value: %v", tt.expectedVal, setting.GetValue())
			}
		})
	}
}
