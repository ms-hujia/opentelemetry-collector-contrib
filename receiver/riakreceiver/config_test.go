// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package riakreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/riakreceiver"

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/scraper/scraperhelper"
	"go.uber.org/multierr"
)

func TestValidate(t *testing.T) {
	clientConfigInvalid := confighttp.NewDefaultClientConfig()
	clientConfigInvalid.Endpoint = "invalid://endpoint:  12efg"

	clientConfig := confighttp.NewDefaultClientConfig()
	clientConfig.Endpoint = defaultEndpoint

	testCases := []struct {
		desc        string
		cfg         *Config
		expectedErr error
	}{
		{
			desc: "missing username, password, and invalid endpoint",
			cfg: &Config{
				ClientConfig:     clientConfigInvalid,
				ControllerConfig: scraperhelper.NewDefaultControllerConfig(),
			},
			expectedErr: multierr.Combine(
				errMissingUsername,
				errMissingPassword,
				fmt.Errorf("%w: %s", errInvalidEndpoint, `parse "invalid://endpoint:  12efg": invalid port ":  12efg" after host`),
			),
		},
		{
			desc: "missing password and invalid endpoint",
			cfg: &Config{
				Username:         "otelu",
				ClientConfig:     clientConfigInvalid,
				ControllerConfig: scraperhelper.NewDefaultControllerConfig(),
			},
			expectedErr: multierr.Combine(
				errMissingPassword,
				fmt.Errorf("%w: %s", errInvalidEndpoint, `parse "invalid://endpoint:  12efg": invalid port ":  12efg" after host`),
			),
		},
		{
			desc: "missing username and invalid endpoint",
			cfg: &Config{
				Password:         "otelp",
				ClientConfig:     clientConfigInvalid,
				ControllerConfig: scraperhelper.NewDefaultControllerConfig(),
			},
			expectedErr: multierr.Combine(
				errMissingUsername,
				fmt.Errorf("%w: %s", errInvalidEndpoint, `parse "invalid://endpoint:  12efg": invalid port ":  12efg" after host`),
			),
		},
		{
			desc: "invalid endpoint",
			cfg: &Config{
				Username:         "otelu",
				Password:         "otelp",
				ClientConfig:     clientConfigInvalid,
				ControllerConfig: scraperhelper.NewDefaultControllerConfig(),
			},
			expectedErr: multierr.Combine(
				fmt.Errorf("%w: %s", errInvalidEndpoint, `parse "invalid://endpoint:  12efg": invalid port ":  12efg" after host`),
			),
		},
		{
			desc: "valid config",
			cfg: &Config{
				Username:         "otelu",
				Password:         "otelp",
				ClientConfig:     clientConfig,
				ControllerConfig: scraperhelper.NewDefaultControllerConfig(),
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actualErr := tc.cfg.Validate()
			if tc.expectedErr != nil {
				require.EqualError(t, actualErr, tc.expectedErr.Error())
			} else {
				require.NoError(t, actualErr)
			}
		})
	}
}
