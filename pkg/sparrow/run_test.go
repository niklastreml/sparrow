// sparrow
// (C) 2023, Deutsche Telekom IT GmbH
//
// Deutsche Telekom IT GmbH and all other contributors /
// copyright owners license this file to you under the Apache
// License, Version 2.0 (the "License"); you may not use this
// file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package sparrow

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/caas-team/sparrow/pkg/sparrow/targets"

	"github.com/caas-team/sparrow/pkg/api"
	"github.com/caas-team/sparrow/pkg/checks/dns"

	"github.com/caas-team/sparrow/pkg/checks/runtime"

	"github.com/caas-team/sparrow/pkg/checks"
	"github.com/caas-team/sparrow/pkg/checks/health"
	"github.com/caas-team/sparrow/pkg/checks/latency"

	gitlabmock "github.com/caas-team/sparrow/pkg/sparrow/targets/test"

	"github.com/stretchr/testify/assert"

	"github.com/caas-team/sparrow/internal/logger"
	"github.com/caas-team/sparrow/pkg/config"
)

func TestSparrow_ReconcileChecks(t *testing.T) {
	ctx, cancel := logger.NewContextWithLogger(context.Background())
	defer cancel()

	rtcfg := &runtime.Config{}
	tests := []struct {
		name             string
		checks           map[string]checks.Check
		newRuntimeConfig runtime.Config
	}{
		{
			name: "no checks registered yet but register one",

			checks: map[string]checks.Check{},

			newRuntimeConfig: runtime.Config{Health: &health.Config{
				Targets: []string{"https://gitlab.com"},
			}},
		},
		{
			name:   "no checks registered, register multiple new ones",
			checks: map[string]checks.Check{},
			newRuntimeConfig: runtime.Config{
				Health: &health.Config{
					Targets: []string{"https://gitlab.com"},
				},
				Latency: &latency.Config{
					Targets: []string{"https://gitlab.com"},
				},
				Dns: &dns.Config{
					Targets: []string{"gitlab.com"},
				},
			},
		},
		{
			name: "one healtcheck registered, register latency check",

			checks: map[string]checks.Check{
				health.CheckName: health.NewCheck(),
			},

			newRuntimeConfig: runtime.Config{
				Latency: &latency.Config{
					Targets: []string{"https://gitlab.com"},
				},
				Health: &health.Config{
					Targets: []string{"https://gitlab.com"},
				},
			},
		},
		{
			name: "no checks registered but unregister all",

			checks: map[string]checks.Check{
				health.CheckName: health.NewCheck(),
			},

			newRuntimeConfig: *rtcfg,
		},
		{
			name: "one health check registered, register latency and unregister health",

			checks: map[string]checks.Check{
				health.CheckName: health.NewCheck(),
			},

			newRuntimeConfig: runtime.Config{
				Latency: &latency.Config{
					Targets: []string{"https://gitlab.com"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sparrow{
				config: &config.Config{SparrowName: "sparrow.com"},
				tarMan: &gitlabmock.MockTargetManager{
					Targets: []checks.GlobalTarget{
						{
							Url: "https://gitlab.com",
						},
					},
				},
				metrics:      NewMetrics(),
				errorHandler: errorHandler{},
				checkCoordinator: checkCoordinator{
					checks:      tt.checks,
					resultFanIn: make(map[string]chan checks.Result),
					cRuntime:    make(chan runtime.Config, 1),
					cResult:     make(chan checks.ResultDTO),
				},
			}

			s.ReconcileChecks(ctx, tt.newRuntimeConfig)

			// iterate of the sparrow's checks and check if they are configured
			for _, c := range s.checks {
				cfg := c.GetConfig()
				assert.NotNil(t, cfg)
				if cfg.For() == health.CheckName {
					assert.Equal(t, tt.newRuntimeConfig.Health, cfg)
				}
				if cfg.For() == latency.CheckName {
					assert.Equal(t, tt.newRuntimeConfig.Latency, cfg)
				}
				if cfg.For() == dns.CheckName {
					assert.Equal(t, tt.newRuntimeConfig.Dns, cfg)
				}
			}
		})
	}
}

func Test_fanInResults(t *testing.T) {
	checkChan := make(chan checks.Result, 1)
	cResult := make(chan checks.ResultDTO, 1)
	name := "check"
	go fanInResults(checkChan, cResult, name)

	result := checks.Result{
		Timestamp: time.Time{},
		Err:       "",
		Data:      0,
	}

	checkChan <- result
	output := <-cResult

	want := checks.ResultDTO{
		Name:   name,
		Result: &result,
	}

	if !reflect.DeepEqual(output, want) {
		t.Errorf("fanInResults() = %v, want %v", output, want)
	}

	close(checkChan)
}

// TestSparrow_Run_FullComponentStart tests that the Run method starts the API,
// loader and a targetManager all start.
func TestSparrow_Run_FullComponentStart(t *testing.T) {
	// create simple file loader config
	c := &config.Config{
		Api: api.Config{ListeningAddress: ":9090"},
		Loader: config.LoaderConfig{
			Type:     "file",
			File:     config.FileLoaderConfig{Path: "../config/testdata/config.yaml"},
			Interval: time.Second * 1,
		},
		TargetManager: targets.TargetManagerConfig{
			Config: targets.Config{
				CheckInterval:        time.Second * 1,
				RegistrationInterval: time.Second * 1,
				UnhealthyThreshold:   time.Second * 1,
			},
			Gitlab: targets.GitlabTargetManagerConfig{
				BaseURL:   "https://gitlab.com",
				Token:     "my-cool-token",
				ProjectID: 42,
			},
		},
	}

	// start sparrow
	s := New(c)
	ctx := context.Background()
	go func() {
		err := s.Run(ctx)
		if err != nil {
			t.Errorf("Sparrow.Run() error = %v", err)
		}
	}()

	t.Log("Running sparrow for 10ms")
	time.Sleep(time.Millisecond * 10)
}

// TestSparrow_Run_ContextCancel tests that after a context cancels the Run method
// will return an error and all started components will be shut down.
func TestSparrow_Run_ContextCancel(t *testing.T) {
	// create simple file loader config
	c := &config.Config{
		Api: api.Config{ListeningAddress: ":9090"},
		Loader: config.LoaderConfig{
			Type:     "file",
			File:     config.FileLoaderConfig{Path: "../config/testdata/config.yaml"},
			Interval: time.Second * 1,
		},
	}

	// start sparrow
	s := New(c)
	s.tarMan = &gitlabmock.MockTargetManager{}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err := s.Run(ctx)
		t.Logf("Sparrow exited with error: %v", err)
		if err == nil {
			t.Error("Sparrow.Run() should have errored out, no error received")
		}
	}()

	t.Log("Running sparrow for 10ms")
	time.Sleep(time.Millisecond * 10)

	t.Log("Canceling context and waiting for shutdown")
	cancel()
	time.Sleep(time.Millisecond * 30)
}

// TestSparrow_enrichTargets tests that the enrichTargets method
// updates the targets of the configured checks.
func TestSparrow_enrichTargets(t *testing.T) {
	now := time.Now()
	testTarget := "https://localhost.de"
	gt := []checks.GlobalTarget{
		{
			Url:      testTarget,
			LastSeen: now,
		},
	}
	tests := []struct {
		name          string
		config        runtime.Config
		globalTargets []checks.GlobalTarget
		expected      runtime.Config
	}{
		{
			name:          "no config",
			config:        runtime.Config{},
			globalTargets: gt,
			expected:      runtime.Config{},
		},
		{
			name: "config with no targets",
			config: runtime.Config{
				Health: &health.Config{
					Targets: nil,
				},
				Latency: &latency.Config{
					Targets: nil,
				},
			},
			globalTargets: gt,
			expected: runtime.Config{
				Health: &health.Config{
					Targets: []string{testTarget},
				},
				Latency: &latency.Config{
					Targets: []string{testTarget},
				},
			},
		},
		{
			name: "config with empty targets",
			config: runtime.Config{
				Health: &health.Config{
					Targets: nil,
				},
				Latency: &latency.Config{
					Targets: nil,
				},
			},
			globalTargets: gt,
			expected: runtime.Config{
				Health: &health.Config{
					Targets: []string{testTarget},
				},
				Latency: &latency.Config{
					Targets: []string{testTarget},
				},
			},
		},
		{
			name: "config with targets (health + latency)",
			config: runtime.Config{
				Health: &health.Config{
					Targets: []string{"https://gitlab.com"},
				},
				Latency: &latency.Config{
					Targets: []string{"https://gitlab.com"},
				},
			},
			globalTargets: gt,
			expected: runtime.Config{
				Health: &health.Config{
					Targets: []string{"https://gitlab.com", testTarget},
				},
				Latency: &latency.Config{
					Targets: []string{"https://gitlab.com", testTarget},
				},
			},
		},
		{
			name: "config with targets (dns)",
			config: runtime.Config{
				Dns: &dns.Config{
					Targets: []string{"gitlab.com"},
				},
			},
			globalTargets: gt,
			expected: runtime.Config{
				Dns: &dns.Config{
					Targets: []string{"gitlab.com", "localhost.de"},
				},
			},
		},
		{
			name: "config has a target already present in global targets - no duplicates",
			config: runtime.Config{
				Health: &health.Config{
					Targets: []string{testTarget},
				},
			},
			globalTargets: gt,
			expected: runtime.Config{
				Health: &health.Config{
					Targets: []string{testTarget},
				},
			},
		},
		{
			name: "global targets contains self - do not add to config",
			config: runtime.Config{
				Health: &health.Config{
					Targets: []string{testTarget},
				},
			},
			globalTargets: append(gt, checks.GlobalTarget{
				Url:      "https://sparrow.com",
				LastSeen: now,
			}),
			expected: runtime.Config{
				Health: &health.Config{
					Targets: []string{testTarget},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sparrow{
				tarMan: &gitlabmock.MockTargetManager{
					Targets: tt.globalTargets,
				},
				config: &config.Config{
					SparrowName: "sparrow.com",
				},
			}
			got := s.enrichTargets(tt.config)
			assert.Equal(t, tt.expected, got)
		})
	}
}
