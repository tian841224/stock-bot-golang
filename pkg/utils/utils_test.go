package utils

import (
	"testing"
)

func TestToString(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{123, "123"},
		{123.45, "123.45"},
		{"hello", "hello"},
		{" test ", "test"},
		{nil, "<nil>"},
	}

	for _, test := range tests {
		result := ToString(test.input)
		if result != test.expected {
			t.Errorf("ToString(%v) = %s; expected %s", test.input, result, test.expected)
		}
	}
}

func TestToFloat(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected float64
	}{
		{"123.45", 123.45},
		{"1,234.56", 1234.56},
		{"123％", 123},
		{"--", 0},
		{"", 0},
		{"+", 0},
		{"-", 0},
		{123, 123},
	}

	for _, test := range tests {
		result := ToFloat(test.input)
		if result != test.expected {
			t.Errorf("ToFloat(%v) = %f; expected %f", test.input, result, test.expected)
		}
	}
}

func TestExtractUpDownSign(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"+5.00", "+"},
		{"-3.20", "-"},
		{"＋2.50", "+"},
		{"－1.80", "-"},
		{"0.00", ""},
		{"", ""},
		{"123", ""},
	}

	for _, test := range tests {
		result := ExtractUpDownSign(test.input)
		if result != test.expected {
			t.Errorf("ExtractUpDownSign(%s) = %s; expected %s", test.input, result, test.expected)
		}
	}
}

func TestFormatNumberWithCommas(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{123, "123"},
		{1234, "1,234"},
		{1234567, "1,234,567"},
		{1234567890, "1,234,567,890"},
	}

	for _, test := range tests {
		result := FormatNumberWithCommas(test.input)
		if result != test.expected {
			t.Errorf("FormatNumberWithCommas(%d) = %s; expected %s", test.input, result, test.expected)
		}
	}
}

func TestFormatAmount(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{1234, "1,234.00"},
		{12345, "1.23萬"},
		{123456789, "1.23億"},
		{1234567890123, "1.23兆"},
	}

	for _, test := range tests {
		result := FormatAmount(test.input)
		if result != test.expected {
			t.Errorf("FormatAmount(%f) = %s; expected %s", test.input, result, test.expected)
		}
	}
}

func TestIsValidStockID(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"2330", true},
		{"0050", true},
		{"123", false},
		{"12345", false},
		{"AAPL", false},
		{"", false},
	}

	for _, test := range tests {
		result := IsValidStockID(test.input)
		if result != test.expected {
			t.Errorf("IsValidStockID(%s) = %t; expected %t", test.input, result, test.expected)
		}
	}
}

func TestIsValidDate(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"2023-12-01", true},
		{"2023-1-1", false},
		{"23-12-01", false},
		{"2023/12/01", false},
		{"", false},
	}

	for _, test := range tests {
		result := IsValidDate(test.input)
		if result != test.expected {
			t.Errorf("IsValidDate(%s) = %t; expected %t", test.input, result, test.expected)
		}
	}
}
