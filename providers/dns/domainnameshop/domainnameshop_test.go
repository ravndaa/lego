package domainnameshop

import (
	"testing"

	"github.com/go-acme/lego/v4/platform/tester"
	"github.com/stretchr/testify/require"
)

const envDomain = envNamespace + "DOMAIN"

var envTest = tester.NewEnvTest(
	EnvToken,
	EnvSecret).
	WithDomain(envDomain)

	/*
		func setup() (*DNSProvider, *http.ServeMux, func()) {
			handler := http.NewServeMux()
			server := httptest.NewServer(handler)

			config := NewDefaultConfig()
			config.Token = "TOKEN"
			config.Secret = "SECRET"

			provider, err := NewDNSProviderConfig(config)
			if err != nil {
				panic(err)
			}

			return provider, handler, server.Close
		}
	*/

func TestNewDNSProvider(t *testing.T) {
	testCases := []struct {
		desc     string
		envVars  map[string]string
		expected string
	}{
		{
			desc: "success",
			envVars: map[string]string{
				EnvToken:  "A",
				EnvSecret: "B",
			},
		},
		{
			desc: "missing credentials",
			envVars: map[string]string{
				EnvToken:  "",
				EnvSecret: "",
			},
			expected: "domainnameshop: some credentials information are missing: DOMAINNAMESHOP_TOKEN,DOMAINNAMESHOP_SECRET",
		},
		{
			desc: "missing secret",
			envVars: map[string]string{
				EnvToken:  "A",
				EnvSecret: "",
			},
			expected: "domainnameshop: some credentials information are missing: DOMAINNAMESHOP_SECRET",
		},
		{
			desc: "missing token",
			envVars: map[string]string{
				EnvToken:  "",
				EnvSecret: "B",
			},
			expected: "domainnameshop: some credentials information are missing: DOMAINNAMESHOP_TOKEN",
		},
	}
	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			defer envTest.RestoreEnv()
			envTest.ClearEnv()

			envTest.Apply(test.envVars)

			p, err := NewDNSProvider()

			if len(test.expected) == 0 {
				require.NoError(t, err)
				require.NotNil(t, p)
				require.NotNil(t, p.config)
			} else {
				require.EqualError(t, err, test.expected)
			}
		})
	}
}

func TestNewDNSProviderConfig(t *testing.T) {
	testCases := []struct {
		desc     string
		token    string
		secret   string
		expected string
	}{
		{
			desc:   "success",
			token:  "A",
			secret: "B",
		},
		{
			desc:     "missing credentials",
			token:    "",
			secret:   "",
			expected: "domainnameshop: some credentials information are missing: DOMAINNAMESHOP_TOKEN,DOMAINNAMESHOP_SECRET",
		},
		{
			desc:     "missing token",
			token:    "",
			secret:   "B",
			expected: "domainnameshop: some credentials information are missing: DOMAINNAMESHOP_TOKEN",
		},
		{
			desc:     "missing secret",
			token:    "A",
			secret:   "",
			expected: "domainnameshop: some credentials information are missing: DOMAINNAMESHOP_SECRET",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			config := NewDefaultConfig()
			config.Token = test.token
			config.Secret = test.secret

			p, err := NewDNSProviderConfig(config)

			if len(test.expected) == 0 {
				require.NoError(t, err)
				require.NotNil(t, p)
				require.NotNil(t, p.config)
			} else {
				require.EqualError(t, err, test.expected)
			}
		})
	}
}
