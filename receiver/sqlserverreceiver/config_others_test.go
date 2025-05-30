// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build !windows

package sqlserverreceiver

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/confmap/xconfmap"
	"go.opentelemetry.io/collector/scraper/scraperhelper"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sqlserverreceiver/internal/metadata"
)

func TestValidateOtherOS(t *testing.T) {
	testCases := []struct {
		desc            string
		cfg             *Config
		expectedSuccess bool
	}{
		{
			desc: "valid config",
			cfg: &Config{
				MetricsBuilderConfig: metadata.DefaultMetricsBuilderConfig(),
				ControllerConfig:     scraperhelper.NewDefaultControllerConfig(),
			},
			expectedSuccess: true,
		},
		{
			desc: "valid config with no metric settings",
			cfg: &Config{
				ControllerConfig: scraperhelper.NewDefaultControllerConfig(),
			},
			expectedSuccess: true,
		},
		{
			desc:            "default config is valid",
			cfg:             createDefaultConfig().(*Config),
			expectedSuccess: true,
		},
		{
			desc: "valid config with both names set",
			cfg: &Config{
				ControllerConfig: scraperhelper.NewDefaultControllerConfig(),
				ComputerName:     "ComputerName",
				InstanceName:     "InstanceName",
			},
			expectedSuccess: true,
		},
		{
			desc: "valid config with instance_name but not computer_name",
			cfg: &Config{
				ControllerConfig: scraperhelper.NewDefaultControllerConfig(),
				InstanceName:     "InstanceName",
			},
			expectedSuccess: true,
		},
		{
			desc: "valid config with computer_name but not instance_name",
			cfg: &Config{
				ControllerConfig: scraperhelper.NewDefaultControllerConfig(),
				ComputerName:     "ComputerName",
			},
			expectedSuccess: true,
		},
		{
			desc: "valid config with both instance and computer name",
			cfg: &Config{
				ControllerConfig: scraperhelper.NewDefaultControllerConfig(),
				ComputerName:     "ComputerName",
				InstanceName:     "InstanceName",
			},
			expectedSuccess: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			if tc.expectedSuccess {
				require.NoError(t, xconfmap.Validate(tc.cfg))
			} else {
				require.Error(t, xconfmap.Validate(tc.cfg))
			}
		})
	}
}
