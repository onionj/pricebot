package price

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestPrice_Refresh(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Return mock response
		mockResponse := `{
			"current": {
				"price_dollar_rl": {"p": "500000", "ts": "2024-03-20 12:00:00", "t": "12:00"},
				"price_eur": {"p": "550000", "ts": "2024-03-20 12:00:00", "t": "12:00"},
				"price_gbp": {"p": "600000", "ts": "2024-03-20 12:00:00", "t": "12:00"},
				"price_cad": {"p": "400000", "ts": "2024-03-20 12:00:00", "t": "12:00"},
				"price_aud": {"p": "350000", "ts": "2024-03-20 12:00:00", "t": "12:00"},
				"price_aed": {"p": "140000", "ts": "2024-03-20 12:00:00", "t": "12:00"},
				"price_try": {"p": "20000", "ts": "2024-03-20 12:00:00", "t": "12:00"},
				"price_sek": {"p": "50000", "ts": "2024-03-20 12:00:00", "t": "12:00"},
				"price_cny": {"p": "80000", "ts": "2024-03-20 12:00:00", "t": "12:00"},
				"price_sar": {"p": "130000", "ts": "2024-03-20 12:00:00", "t": "12:00"},
				"price_iqd": {"p": "400", "ts": "2024-03-20 12:00:00", "t": "12:00"},
				"crypto-tether-irr": {"p": "510000", "ts": "2024-03-20 12:00:00", "t": "12:00"},
				"crypto-bitcoin": {"p": "65000", "ts": "2024-03-20 12:00:00", "t": "12:00"},
				"crypto-ethereum": {"p": "3500", "ts": "2024-03-20 12:00:00", "t": "12:00"},
				"sekeb": {"p": "30000000", "ts": "2024-03-20 12:00:00", "t": "12:00"},
				"sekee": {"p": "31000000", "ts": "2024-03-20 12:00:00", "t": "12:00"},
				"nim": {"p": "15000000", "ts": "2024-03-20 12:00:00", "t": "12:00"},
				"rob": {"p": "9000000", "ts": "2024-03-20 12:00:00", "t": "12:00"},
				"rob_down": {"p": "6000000", "ts": "2024-03-20 12:00:00", "t": "12:00"},
				"geram18": {"p": "3000000", "ts": "2024-03-20 12:00:00", "t": "12:00"},
				"mesghal": {"p": "13000000", "ts": "2024-03-20 12:00:00", "t": "12:00"},
				"ons": {"p": "2200", "ts": "2024-03-20 12:00:00", "t": "12:00"}
			}
		}`
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Create price instance
	p := NewPrice()

	// Replace the HTTP client with our test client
	httpClient = server.Client()

	// Replace the API URL with our mock server URL
	baseURL = server.URL

	// Test Refresh
	err := p.Refresh()
	if err != nil {
		t.Errorf("Refresh failed: %v", err)
	}

	// Verify some values
	if p.Current.Dollar.Price != "500000" {
		t.Errorf("Expected Dollar price '500000', got '%s'", p.Current.Dollar.Price)
	}
	if p.Current.BitCoin.Price != "65000" {
		t.Errorf("Expected Bitcoin price '65000', got '%s'", p.Current.BitCoin.Price)
	}
}

func TestPrice_String(t *testing.T) {
	p := &Price{
		Current: CurrentData{
			Dollar:  Detail{Price: "500000", Time: "12:00"},
			Eur:     Detail{Price: "550000", Time: "12:00"},
			BitCoin: Detail{Price: "65000", Time: "12:00"},
		},
		LastRefresh: time.Date(2024, 3, 20, 12, 0, 0, 0, time.UTC),
	}

	result := p.String()

	// Test that the output contains expected values
	expectedStrings := []string{
		"دلار امریکا (12:00) ⬅️ *50,000* تومان",
		"یورو اروپا (12:00) ⬅️ *55,000* تومان",
		"بیتکوین (12:00) ⬅️ *65000* دلار",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("Expected string to contain '%s', but it didn't", expected)
		}
	}
}

func TestPrice_PrettyNumber(t *testing.T) {
	p := NewPrice()
	testCases := []struct {
		input    int
		expected string
	}{
		{1234, "1,234"},
		{1234567, "1,234,567"},
		{1000000, "1,000,000"},
		{123, "123"},
		{0, "0"},
	}

	for _, tc := range testCases {
		result := p.prettyNumber(tc.input)
		if result != tc.expected {
			t.Errorf("prettyNumber(%d) = %s; want %s", tc.input, result, tc.expected)
		}
	}
}

func TestPrice_ToToman(t *testing.T) {
	p := NewPrice()
	testCases := []struct {
		input    string
		expected string
	}{
		{"500000", "50,000"},
		{"1,000,000", "100,000"},
		{"invalid", "0"},
		{"", "0"},
	}

	for _, tc := range testCases {
		result := p.toToman(tc.input)
		if result != tc.expected {
			t.Errorf("toToman(%s) = %s; want %s", tc.input, result, tc.expected)
		}
	}
}
