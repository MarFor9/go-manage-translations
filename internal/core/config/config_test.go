package config

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go-template/internal/log"
	"os"
	"testing"
)

func Test_Load(t *testing.T) {

	cases := []struct {
		name          string
		setupEnv      func()
		expectedError bool
	}{
		{
			name: "valid configuration",
			setupEnv: func() {
				setEnvironment()
			},
			expectedError: false,
		},
		{
			name: "missing SERVER_URL",
			setupEnv: func() {
				setEnvironment()
				os.Unsetenv("SERVER_URL")
			},
			expectedError: true,
		},
		{
			name: "missing SERVER_PORT",
			setupEnv: func() {
				setEnvironment()
				os.Unsetenv("SERVER_PORT")
			},
			expectedError: true,
		},
		{
			name: "missing LOG_LEVEL",
			setupEnv: func() {
				setEnvironment()
				os.Unsetenv("LOG_LEVEL")
			},
			expectedError: false,
		},
		{
			name: "wrong value for LOG_LEVEL",
			setupEnv: func() {
				setEnvironment()
				os.Setenv("LOG_LEVEL", "10")
			},
			expectedError: true,
		},
		{
			name: "missing LOG_MODE",
			setupEnv: func() {
				setEnvironment()
				os.Unsetenv("LOG_MODE")
			},
			expectedError: true,
		},
		{
			name: "wrong value for LOG_MODE",
			setupEnv: func() {
				setEnvironment()
				os.Setenv("LOG_MODE", "3")
			},
			expectedError: true,
		},
		{
			name: "missing TRANSLATION_SERVICE_API_KEY",
			setupEnv: func() {
				setEnvironment()
				os.Unsetenv("TRANSLATION_SERVICE_API_KEY")
			},
			expectedError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupEnv()
			_, err := Load()
			if tc.expectedError {
				log.Info(context.Background(), err.Error())
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func setEnvironment() {
	os.Setenv("SERVER_URL", "http://localhost")
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("LOG_LEVEL", "-4")
	os.Setenv("LOG_MODE", "2")
	os.Setenv("TRANSLATION_SERVICE_API_KEY", "secure-api-key")
}
