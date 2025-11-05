package sb

import "testing"

func TestRawColumnProcessNilPanic(t *testing.T) {
	tests := []struct {
		name       string
		columnType string
		expectType string
		expectLen  string
		expectDec  string
	}{
		{
			name:       "normal case with parentheses",
			columnType: "VARCHAR(255)",
			expectType: "VARCHAR",
			expectLen:  "255",
			expectDec:  "",
		},
		{
			name:       "normal case with decimals",
			columnType: "DECIMAL(10,2)",
			expectType: "DECIMAL",
			expectLen:  "10",
			expectDec:  "2",
		},
		{
			name:       "edge case - only opening parenthesis",
			columnType: "VARCHAR(",
			expectType: "VARCHAR",
			expectLen:  "",
			expectDec:  "",
		},
		{
			name:       "edge case - empty string",
			columnType: "",
			expectType: "",
			expectLen:  "",
			expectDec:  "",
		},
		{
			name:       "edge case - no parentheses",
			columnType: "TEXT",
			expectType: "TEXT",
			expectLen:  "",
			expectDec:  "",
		},
		{
			name:       "edge case - malformed with comma but no second part",
			columnType: "DECIMAL(10,)",
			expectType: "DECIMAL",
			expectLen:  "10",
			expectDec:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, gotLen, gotDec := rawColumnProcess(tt.columnType)
			if gotType != tt.expectType {
				t.Errorf("rawColumnProcess() type = %v, want %v", gotType, tt.expectType)
			}
			if gotLen != tt.expectLen {
				t.Errorf("rawColumnProcess() length = %v, want %v", gotLen, tt.expectLen)
			}
			if gotDec != tt.expectDec {
				t.Errorf("rawColumnProcess() decimals = %v, want %v", gotDec, tt.expectDec)
			}
		})
	}
}
