package conf

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	type Level2 struct {
		Count int    `default:"3"`
		Name  string `default:"pidgeon"`
	}
	type Nest struct {
		Egg    string `default:"chicken"`
		Level2 Level2
	}
	cfg := struct {
		ApiPort          int     `default:"8080"`
		Offer            float64 `default:"1.99"`
		Amount           float64 `default:"10.01"`
		ServiceEnv       string  `default:"local"`
		Enabled          bool    `default:"true"`
		Disabled         bool
		CountryPrefixMap map[string]int    `default:"Argentina:54,USA:1,Spain:34"`
		CountryCodeMap   map[string]string `default:"Italy:it"`
		FibonacciSlice   []int             `default:"0,1,1,2,3,5,8"`
		CountSlice       []string          `default:"ichi,ni,san"`
		Nested           Nest
	}{}

	os.Setenv("MYAPP_API_PORT", "9090")
	os.Setenv("MYAPP_NESTED_EGG", "iguana")
	os.Setenv("MYAPP_NESTED_LEVEL2_COUNT", "5")
	os.Setenv("MYAPP_DISABLED", "true")
	os.Setenv("MYAPP_COUNTRY_CODE_MAP", "Argentina:ar,Spain:es,France:fr")
	os.Setenv("MYAPP_COUNT_SLICE", "one,two,three")
	os.Setenv("MYAPP_AMOUNT", "8.88")

	err := Unmarshal(&cfg, "MYAPP", ",", ":")
	require.NoError(t, err)

	assert.Equal(t, 9090, cfg.ApiPort)
	assert.Equal(t, "local", cfg.ServiceEnv)
	assert.Equal(t, "iguana", cfg.Nested.Egg)
	assert.Equal(t, 5, cfg.Nested.Level2.Count)
	assert.Equal(t, "pidgeon", cfg.Nested.Level2.Name)
	assert.Equal(t, float64(8.88), cfg.Amount)
	assert.Equal(t, float64(1.99), cfg.Offer)
	code, ok := cfg.CountryCodeMap["Argentina"]
	assert.Equal(t, true, ok)
	assert.Equal(t, "ar", code)
	code, ok = cfg.CountryCodeMap["Spain"]
	assert.Equal(t, true, ok)
	assert.Equal(t, "es", code)
	code, ok = cfg.CountryCodeMap["France"]
	assert.Equal(t, true, ok)
	assert.Equal(t, "fr", code)
	code, ok = cfg.CountryCodeMap["Italy"]
	assert.Equal(t, false, ok)
	prefix, ok := cfg.CountryPrefixMap["Argentina"]
	assert.Equal(t, true, ok)
	assert.Equal(t, 54, prefix)
	prefix, ok = cfg.CountryPrefixMap["USA"]
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, prefix)
	prefix, ok = cfg.CountryPrefixMap["Spain"]
	assert.Equal(t, true, ok)
	assert.Equal(t, 34, prefix)
	for i, v := range cfg.FibonacciSlice {
		if i < 2 {
			continue
		}
		assert.Equal(t, cfg.FibonacciSlice[i-2]+cfg.FibonacciSlice[i-1], v)
	}
	assert.Equal(t, "one", cfg.CountSlice[0])
	assert.Equal(t, "two", cfg.CountSlice[1])
	assert.Equal(t, "three", cfg.CountSlice[2])
}
