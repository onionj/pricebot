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

		mockResponse := `{
			"current": {
				"price_dollar_rl": {"p": "500000", "ts": "2024-03-20 12:00:00", "t": "12:00", "dp": 2.45, "dt": "high"},
				"price_eur": {"p": "550000", "ts": "2024-03-20 12:00:00", "t": "12:00", "dp": 1.23, "dt": "low"},
				"price_gbp": {"p": "600000", "ts": "2024-03-20 12:00:00", "t": "12:00", "dp": 0, "dt": ""},
				"price_cad": {"p": "400000", "ts": "2024-03-20 12:00:00", "t": "12:00", "dp": 1.5, "dt": "high"},
				"price_aud": {"p": "350000", "ts": "2024-03-20 12:00:00", "t": "12:00", "dp": 0.5, "dt": "low"},
				"price_aed": {"p": "140000", "ts": "2024-03-20 12:00:00", "t": "12:00", "dp": 0.75, "dt": "high"},
				"price_try": {"p": "20000", "ts": "2024-03-20 12:00:00", "t": "12:00", "dp": 2.1, "dt": "low"},
				"price_sek": {"p": "50000", "ts": "2024-03-20 12:00:00", "t": "12:00", "dp": 0, "dt": ""},
				"price_cny": {"p": "80000", "ts": "2024-03-20 12:00:00", "t": "12:00", "dp": 1.1, "dt": "high"},
				"price_sar": {"p": "130000", "ts": "2024-03-20 12:00:00", "t": "12:00", "dp": 0.8, "dt": "low"},
				"price_iqd": {"p": "400", "ts": "2024-03-20 12:00:00", "t": "12:00", "dp": 0.3, "dt": "high"},
				"crypto-tether-irr": {"p": "510000", "ts": "2024-03-20 12:00:00", "t": "12:00", "dp": 0.1, "dt": "high"},
				"crypto-bitcoin": {"p": "65000", "ts": "2024-03-20 12:00:00", "t": "12:00", "dp": 5.2, "dt": "high"},
				"crypto-ethereum": {"p": "3500", "ts": "2024-03-20 12:00:00", "t": "12:00", "dp": 3.1, "dt": "low"},
				"sekeb": {"p": "30000000", "ts": "2024-03-20 12:00:00", "t": "12:00", "dp": 1.8, "dt": "high"},
				"sekee": {"p": "31000000", "ts": "2024-03-20 12:00:00", "t": "12:00", "dp": 2.0, "dt": "high"},
				"nim": {"p": "15000000", "ts": "2024-03-20 12:00:00", "t": "12:00", "dp": 0.9, "dt": "low"},
				"rob": {"p": "9000000", "ts": "2024-03-20 12:00:00", "t": "12:00", "dp": 1.5, "dt": "low"},
				"rob_down": {"p": "6000000", "ts": "2024-03-20 12:00:00", "t": "12:00", "dp": 1.2, "dt": "low"},
				"geram18": {"p": "3000000", "ts": "2024-03-20 12:00:00", "t": "12:00", "dp": 0.7, "dt": "high"},
				"mesghal": {"p": "13000000", "ts": "2024-03-20 12:00:00", "t": "12:00", "dp": 0.5, "dt": "high"},
				"ons": {"p": "2200", "ts": "2024-03-20 12:00:00", "t": "12:00", "dp": 0.4, "dt": "low"}
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

	// Test cases for price values and changes
	testCases := []struct {
		name          string
		price         string
		changePercent float64
		changeDir     string
		wantFormat    string
	}{
		{"Dollar", p.Current.Dollar.Price, p.Current.Dollar.ChangePercentage, p.Current.Dollar.ChangeDirection, "(2.45%+🟢)"},
		{"Euro", p.Current.Eur.Price, p.Current.Eur.ChangePercentage, p.Current.Eur.ChangeDirection, "(1.23%-🔴)"},
		{"GBP", p.Current.GBP.Price, p.Current.GBP.ChangePercentage, p.Current.GBP.ChangeDirection, "⬅️"},
		{"Bitcoin", p.Current.BitCoin.Price, p.Current.BitCoin.ChangePercentage, p.Current.BitCoin.ChangeDirection, "(5.20%+🟢)"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			detail := Detail{
				Price:            tc.price,
				ChangePercentage: tc.changePercent,
				ChangeDirection:  tc.changeDir,
			}
			got := detail.FormatChange()
			if got != tc.wantFormat {
				t.Errorf("FormatChange() = %v, want %v", got, tc.wantFormat)
			}
		})
	}

	// Additional verification for specific values
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
			Dollar:  Detail{Price: "500000", Time: "12:00", ChangePercentage: 2.45, ChangeDirection: "high"},
			Eur:     Detail{Price: "550000", Time: "12:00", ChangePercentage: 1.23, ChangeDirection: "low"},
			BitCoin: Detail{Price: "65000", Time: "12:00", ChangePercentage: 5.20, ChangeDirection: "high"},
			GBP:     Detail{Price: "600000", Time: "12:00", ChangePercentage: 0, ChangeDirection: ""},
		},
		LastRefresh: time.Date(2024, 3, 20, 12, 0, 0, 0, time.UTC),
	}

	result := p.String()

	// Test that the output contains expected values
	expectedStrings := []string{
		"دلار امریکا (12:00) (2.45%+🟢) *50,000* تومان",
		"یورو اروپا (12:00) (1.23%-🔴) *55,000* تومان",
		"بیتکوین (12:00) (5.20%+🟢) *65000* دلار",
		"پوند انگلیس (12:00) ⬅️ *60,000* تومان",
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
