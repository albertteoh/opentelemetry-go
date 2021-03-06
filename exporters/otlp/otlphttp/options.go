// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package otlphttp

import (
	"crypto/tls"
	"time"
)

// Compression describes the compression used for payloads sent to the
// collector.
type Compression int

const (
	// NoCompression tells the driver to send payloads without
	// compression.
	NoCompression Compression = iota
	// GzipCompression tells the driver to send payloads after
	// compressing them with gzip.
	GzipCompression
)

const (
	// DefaultMaxAttempts describes how many times the driver
	// should retry the sending of the payload in case of a
	// retryable error.
	DefaultMaxAttempts int = 5
	// DefaultTracesPath is a default URL path for endpoint that
	// receives spans.
	DefaultTracesPath string = "/v1/traces"
	// DefaultMetricsPath is a default URL path for endpoint that
	// receives metrics.
	DefaultMetricsPath string = "/v1/metrics"
	// DefaultBackoff is a default base backoff time used in the
	// exponential backoff strategy.
	DefaultBackoff time.Duration = 300 * time.Millisecond
)

type config struct {
	endpoint       string
	compression    Compression
	tracesURLPath  string
	metricsURLPath string
	maxAttempts    int
	backoff        time.Duration
	tlsCfg         *tls.Config
	insecure       bool
	headers        map[string]string
}

// Option applies an option to the HTTP driver.
type Option interface {
	Apply(*config)
}

type endpointOption string

func (o endpointOption) Apply(cfg *config) {
	cfg.endpoint = (string)(o)
}

// WithEndpoint allows one to set the address of the collector
// endpoint that the driver will use to send metrics and spans. If
// unset, it will instead try to use
// DefaultCollectorHost:DefaultCollectorPort. Note that the endpoint
// must not contain any URL path.
func WithEndpoint(endpoint string) Option {
	return (endpointOption)(endpoint)
}

type compressionOption Compression

func (o compressionOption) Apply(cfg *config) {
	cfg.compression = (Compression)(o)
}

// WithCompression tells the driver to compress the sent data.
func WithCompression(compression Compression) Option {
	return (compressionOption)(compression)
}

type tracesURLPathOption string

func (o tracesURLPathOption) Apply(cfg *config) {
	cfg.tracesURLPath = (string)(o)
}

// WithTracesURLPath allows one to override the default URL path used
// for sending traces. If unset, DefaultTracesPath will be used.
func WithTracesURLPath(urlPath string) Option {
	return (tracesURLPathOption)(urlPath)
}

type metricsURLPathOption string

func (o metricsURLPathOption) Apply(cfg *config) {
	cfg.metricsURLPath = (string)(o)
}

// WithMetricsURLPath allows one to override the default URL path used
// for sending metrics. If unset, DefaultMetricsPath will be used.
func WithMetricsURLPath(urlPath string) Option {
	return (metricsURLPathOption)(urlPath)
}

type maxAttemptsOption int

func (o maxAttemptsOption) Apply(cfg *config) {
	cfg.maxAttempts = (int)(o)
}

// WithMaxAttempts allows one to override how many times the driver
// will try to send the payload in case of retryable errors. If unset,
// DefaultMaxAttempts will be used.
func WithMaxAttempts(maxAttempts int) Option {
	return maxAttemptsOption(maxAttempts)
}

type backoffOption time.Duration

func (o backoffOption) Apply(cfg *config) {
	cfg.backoff = (time.Duration)(o)
}

// WithBackoff tells the driver to use the duration as a base of the
// exponential backoff strategy. If unset, DefaultBackoff will be
// used.
func WithBackoff(duration time.Duration) Option {
	return (backoffOption)(duration)
}

type tlsClientConfigOption tls.Config

func (o *tlsClientConfigOption) Apply(cfg *config) {
	cfg.tlsCfg = (*tls.Config)(o)
}

// WithTLSClientConfig can be used to set up a custom TLS
// configuration for the client used to send payloads to the
// collector. Use it if you want to use a custom certificate.
func WithTLSClientConfig(tlsCfg *tls.Config) Option {
	return (*tlsClientConfigOption)(tlsCfg)
}

type insecureOption struct{}

func (insecureOption) Apply(cfg *config) {
	cfg.insecure = true
}

// WithInsecure tells the driver to connect to the collector using the
// HTTP scheme, instead of HTTPS.
func WithInsecure() Option {
	return insecureOption{}
}

type headersOption map[string]string

func (o headersOption) Apply(cfg *config) {
	cfg.headers = (map[string]string)(o)
}

// WithHeaders allows one to tell the driver to send additional HTTP
// headers with the payloads. Specifying headers like Content-Length,
// Content-Encoding and Content-Type may result in a broken driver.
func WithHeaders(headers map[string]string) Option {
	return (headersOption)(headers)
}
