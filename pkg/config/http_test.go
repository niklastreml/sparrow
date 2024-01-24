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

package config

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/caas-team/sparrow/internal/helper"

	"github.com/caas-team/sparrow/internal/logger"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestHttpLoader_GetRuntimeConfig(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	type httpResponder struct {
		statusCode int
		response   string
	}
	tests := []struct {
		name          string
		cfg           *Config
		httpResponder httpResponder
		want          *RuntimeConfig
		wantErr       bool
	}{
		{
			name: "Get runtime configuration",
			cfg: &Config{
				Loader: LoaderConfig{
					Type:     "http",
					Interval: time.Second,
				},
			},
			httpResponder: httpResponder{
				statusCode: 200,
				response:   httpmock.File("testdata/config.yaml").String(),
			},
			want: &RuntimeConfig{
				Checks: map[string]any{
					"testCheck1": map[string]any{
						"enabled": true,
					},
				},
			},
		},
		{
			name: "Get runtime configuration with auth",
			cfg: &Config{
				Loader: LoaderConfig{
					Type:     "http",
					Interval: time.Second,
					Http: HttpLoaderConfig{
						Token: "SECRET",
					},
				},
			},
			httpResponder: httpResponder{
				statusCode: 200,
				response:   httpmock.File("testdata/config.yaml").String(),
			},
			want: &RuntimeConfig{
				Checks: map[string]any{
					"testCheck1": map[string]any{
						"enabled": true,
					},
				},
			},
		},
		{
			name: "Get runtime configuration with statuscode 400",
			cfg: &Config{
				Loader: LoaderConfig{
					Type:     "http",
					Interval: time.Second,
				},
			},
			httpResponder: httpResponder{
				statusCode: 400,
				response:   httpmock.File("testdata/config.yaml").String(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Get runtime configuration payload not yaml",
			cfg: &Config{
				Loader: LoaderConfig{
					Type:     "http",
					Interval: time.Second,
				},
			},
			httpResponder: httpResponder{
				statusCode: 200,
				response:   `this is not yaml`,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			endpoint := "https://api.test.com/test"
			httpmock.RegisterResponder("GET", endpoint,
				func(req *http.Request) (*http.Response, error) {
					if tt.cfg.Loader.Http.Token != "" {
						require.Equal(t, req.Header.Get("Authorization"), fmt.Sprintf("Bearer %s", tt.cfg.Loader.Http.Token))
						fmt.Println("TOKEN tested")
					}
					resp, _ := httpmock.NewStringResponder(tt.httpResponder.statusCode, tt.httpResponder.response)(req)
					return resp, nil
				},
			)

			handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
			ctx := logger.IntoContext(context.Background(), logger.NewLogger(handler).WithGroup("httpLoader-test"))
			ctx, cancel := context.WithTimeout(ctx, time.Second*10)
			defer cancel()

			gl := &HttpLoader{
				cfg:        tt.cfg,
				cCfgChecks: make(chan<- map[string]any, 1),
				client: &http.Client{
					Timeout: tt.cfg.Loader.Http.Timeout,
				},
			}
			gl.cfg.Loader.Http.Url = endpoint

			got, err := gl.GetRuntimeConfig(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("HttpLoader.GetRuntimeConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HttpLoader.GetRuntimeConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestHttpLoader_Run tests the Run method of the HttpLoader
// The test runs the Run method for a while
// and then shuts it down via a goroutine
func TestHttpLoader_Run(t *testing.T) {
	tests := []struct {
		name     string
		response *RuntimeConfig
		code     int
		wantErr  bool
	}{
		{
			name:     "non-200 response",
			response: nil,
			code:     500,
			wantErr:  false,
		},
		{
			name: "nil checks",
			response: &RuntimeConfig{
				Checks: nil,
			},
			code: 200,
		},
		{
			name: "empty config",
			response: &RuntimeConfig{
				Checks: map[string]any{},
			},
			code: 200,
		},
		{
			name: "config with checks",
			response: &RuntimeConfig{
				Checks: map[string]any{
					"testCheck1": map[string]any{
						"enabled": true,
					},
				},
			},
			code: 200,
		},
	}

	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	for _, tt := range tests {
		body, err := yaml.Marshal(tt.response)
		if err != nil {
			t.Fatalf("Failed marshaling response to bytes: %v", err)
		}
		resp := httpmock.NewBytesResponder(tt.code, body)
		httpmock.RegisterResponder("GET", "https://api.test.com/test", resp)

		t.Run(tt.name, func(t *testing.T) {
			hl := &HttpLoader{
				cfg: &Config{
					Loader: LoaderConfig{
						Type:     "http",
						Interval: time.Millisecond * 500,
						Http: HttpLoaderConfig{
							Url: "https://api.test.com/test",
							RetryCfg: helper.RetryConfig{
								Count: 3,
								Delay: 100 * time.Millisecond,
							},
						},
					},
				},
				cCfgChecks: make(chan<- map[string]any, 1),
				client:     http.DefaultClient,
				done:       make(chan struct{}, 1),
			}

			// shutdown routine
			ctx := context.Background()
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				time.Sleep(time.Millisecond * 600)
				t.Log("Shutting down the Run method")
				hl.Shutdown(ctx)
			}()

			err := hl.Run(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("HttpLoader.Run() error = %v, wantErr %v", err, tt.wantErr)
			}

			httpmock.Reset()
		})
	}
}

func TestHttpLoader_Shutdown(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "shutdown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hl := &HttpLoader{
				done: make(chan struct{}, 1),
			}
			hl.Shutdown(context.Background())

			// check if the signal is sent
			select {
			case <-hl.done:
				t.Log("Shutdown signal received")
			default:
				t.Error("Shutdown signal not received")
			}
		})
	}
}

// TestHttpLoader_Run_config_sent_to_channel tests if the config is sent to the channel
// when the Run method is called and the remote endpoint returns a valid response
func TestHttpLoader_Run_config_sent_to_channel(t *testing.T) {
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	expected := map[string]any{
		"testCheck1": map[string]any{
			"enabled": true,
		},
	}
	body, err := yaml.Marshal(RuntimeConfig{Checks: expected})
	if err != nil {
		t.Fatalf("Failed marshaling yaml: %v", err)
	}
	resp := httpmock.NewBytesResponder(200, body)
	httpmock.RegisterResponder("GET", "https://api.test.com/test", resp)

	cCfgChecks := make(chan map[string]any, 1)

	hl := &HttpLoader{
		cfg: &Config{
			Loader: LoaderConfig{
				Type:     "http",
				Interval: time.Millisecond * 500,
				Http: HttpLoaderConfig{
					Url: "https://api.test.com/test",
					RetryCfg: helper.RetryConfig{
						Count: 2,
						Delay: 100 * time.Millisecond,
					},
				},
			},
		},
		cCfgChecks: cCfgChecks,
		client:     http.DefaultClient,
		done:       make(chan struct{}, 1),
	}

	ctx := context.Background()
	go func() {
		err := hl.Run(ctx)
		if err != nil {
			t.Errorf("HttpLoader.Run() error = %v", err)
		}
	}()

	// check if the config is sent to the channel
	select {
	case <-time.After(time.Second):
		t.Error("Config not sent to channel")
	case c := <-cCfgChecks:
		if !reflect.DeepEqual(c, expected) {
			t.Errorf("Config sent to channel is not equal to expected config: got %v, want %v", c, expected)
		}
	}

	hl.Shutdown(ctx)
}

// TestHttpLoader_Run_config_not_sent_to_channel_500 tests if the config is not sent to the channel
// when the Run method is called
// and the remote endpoint returns a non-200 response
func TestHttpLoader_Run_config_not_sent_to_channel_500(t *testing.T) {
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	resp, err := httpmock.NewJsonResponder(500, nil)
	if err != nil {
		t.Fatalf("Failed creating json responder: %v", err)
	}

	httpmock.RegisterResponder("GET", "https://api.test.com/test", resp)

	cCfgChecks := make(chan map[string]any, 1)

	hl := &HttpLoader{
		cfg: &Config{
			Loader: LoaderConfig{
				Type:     "http",
				Interval: time.Millisecond * 500,
				Http: HttpLoaderConfig{
					Url: "https://api.test.com/test",
					RetryCfg: helper.RetryConfig{
						Count: 2,
						Delay: 100 * time.Millisecond,
					},
				},
			},
		},
		cCfgChecks: cCfgChecks,
		client:     http.DefaultClient,
		done:       make(chan struct{}, 1),
	}

	ctx := context.Background()
	go func() {
		err := hl.Run(ctx)
		if err != nil {
			t.Errorf("HttpLoader.Run() error = %v", err)
		}
	}()

	// check if the config is sent to the channel
	select {
	// make sure you wait for at least an interval
	case <-time.After(time.Second):
		t.Log("Config not sent to channel")
	case c := <-cCfgChecks:
		t.Errorf("Config sent to channel: %v", c)
	}

	hl.Shutdown(ctx)
}

// TestHttpLoader_Run_config_not_sent_to_channel_client_error tests if the config is not sent to the channel
// when the Run method is called
// and the client can't execute the requests
func TestHttpLoader_Run_config_not_sent_to_channel_client_error(t *testing.T) {
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	resp := httpmock.NewErrorResponder(fmt.Errorf("client error"))
	httpmock.RegisterResponder("GET", "https://api.test.com/test", resp)

	cCfgChecks := make(chan map[string]any, 1)

	hl := &HttpLoader{
		cfg: &Config{
			Loader: LoaderConfig{
				Type:     "http",
				Interval: time.Millisecond * 500,
				Http: HttpLoaderConfig{
					Url: "https://api.test.com/test",
					RetryCfg: helper.RetryConfig{
						Count: 2,
						Delay: 100 * time.Millisecond,
					},
				},
			},
		},
		cCfgChecks: cCfgChecks,
		client:     http.DefaultClient,
		done:       make(chan struct{}, 1),
	}

	ctx := context.Background()
	go func() {
		err := hl.Run(ctx)
		if err != nil {
			t.Errorf("HttpLoader.Run() error = %v", err)
		}
	}()

	// check if the config is sent to the channel
	select {
	// make sure you wait for at least an interval
	case <-time.After(time.Second):
		t.Log("Config not sent to channel")
	case c := <-cCfgChecks:
		t.Errorf("Config sent to channel: %v", c)
	}

	hl.Shutdown(ctx)
}
