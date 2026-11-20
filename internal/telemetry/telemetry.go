// Copyright Splunk, Inc.
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

package telemetry

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"ktlo/internal/config"
)

const (
	defaultInsecureEnv          = "OTEL_EXPORTER_OTLP_INSECURE"
	defaultTraceEndpointEnv     = "OTEL_EXPORTER_OTLP_ENDPOINT"
	emptyString                 = ""
	errorCreateMetricInstrument = "could not create OpenTelemetry operation duration metric: %w"
	errorCreateMetricExporter   = "could not create OpenTelemetry metric exporter: %w"
	errorCreateTraceExporter    = "could not create OpenTelemetry trace exporter: %w"
	errorSetEnvironmentVariable = "could not set OpenTelemetry environment variable %q: %w"
	metricDurationDescription   = "Duration of a KTLO operation."
	metricDurationName          = "ktlo.operation.duration"
	metricOperationAttribute    = "ktlo.operation"
	metricUnitSeconds           = "s"
	metricsEndpointEnv          = "OTEL_EXPORTER_OTLP_METRICS_ENDPOINT"
	metricsInsecureEnv          = "OTEL_EXPORTER_OTLP_METRICS_INSECURE"
	serviceName                 = "ktlo"
	serviceNameAttribute        = "service.name"
	traceEndpointEnv            = "OTEL_EXPORTER_OTLP_TRACES_ENDPOINT"
	traceInsecureEnv            = "OTEL_EXPORTER_OTLP_TRACES_INSECURE"
)

var operationDuration metric.Float64Histogram

// Setup initializes OTLP tracing and metrics when enabled in the YAML configuration.
func Setup(ctx context.Context, instrumentationConfig config.OTelGoInstrumentationConfig) (func(context.Context) error, error) {
	if !instrumentationConfig.Enabled {
		return func(context.Context) error { return nil }, nil
	}
	if err := setConfiguredEnvironmentVariables(instrumentationConfig); err != nil {
		return nil, err
	}

	defaultEndpoint := strings.TrimSpace(os.Getenv(defaultTraceEndpointEnv))
	traceEnabled := defaultEndpoint != emptyString || strings.TrimSpace(os.Getenv(traceEndpointEnv)) != emptyString
	metricsEnabled := defaultEndpoint != emptyString || strings.TrimSpace(os.Getenv(metricsEndpointEnv)) != emptyString
	if !traceEnabled && !metricsEnabled {
		return func(context.Context) error { return nil }, nil
	}

	var traceExporter sdktrace.SpanExporter
	if traceEnabled {
		exporter, err := otlptracegrpc.New(ctx)
		if err != nil {
			return nil, fmt.Errorf(errorCreateTraceExporter, err)
		}
		traceExporter = exporter
	}

	var metricExporter *otlpmetricgrpc.Exporter
	if metricsEnabled {
		var err error
		metricExporter, err = otlpmetricgrpc.New(ctx)
		if err != nil {
			if traceExporter != nil {
				_ = traceExporter.Shutdown(ctx)
			}
			return nil, fmt.Errorf(errorCreateMetricExporter, err)
		}
	}

	otelResource := resource.NewSchemaless(attribute.String(serviceNameAttribute, serviceName))
	var traceProvider *sdktrace.TracerProvider
	if traceExporter != nil {
		traceProvider = sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(traceExporter),
			sdktrace.WithResource(otelResource),
		)
	}
	var meterProvider *sdkmetric.MeterProvider
	if metricExporter != nil {
		meterProvider = sdkmetric.NewMeterProvider(
			sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
			sdkmetric.WithResource(otelResource),
		)
		var err error
		operationDuration, err = meterProvider.Meter(serviceName).Float64Histogram(
			metricDurationName,
			metric.WithDescription(metricDurationDescription),
			metric.WithUnit(metricUnitSeconds),
		)
		if err != nil {
			_ = meterProvider.Shutdown(ctx)
			if traceProvider != nil {
				_ = traceProvider.Shutdown(ctx)
			}
			return nil, fmt.Errorf(errorCreateMetricInstrument, err)
		}
		otel.SetMeterProvider(meterProvider)
	}
	if traceProvider != nil {
		otel.SetTracerProvider(traceProvider)
	}

	return func(ctx context.Context) error {
		var shutdownErrors []error
		if traceProvider != nil {
			shutdownErrors = append(shutdownErrors, traceProvider.Shutdown(ctx))
		}
		if meterProvider != nil {
			shutdownErrors = append(shutdownErrors, meterProvider.Shutdown(ctx))
		}
		return errors.Join(shutdownErrors...)
	}, nil
}

func setConfiguredEnvironmentVariables(instrumentationConfig config.OTelGoInstrumentationConfig) error {
	settings := [...]struct {
		name  string
		value string
	}{
		{name: defaultTraceEndpointEnv, value: instrumentationConfig.OTELExporterOTLPEndpoint},
		{name: defaultInsecureEnv, value: instrumentationConfig.OTELExporterOTLPInsecure},
		{name: metricsEndpointEnv, value: instrumentationConfig.OTELExporterOTLPMetricsEndpoint},
		{name: metricsInsecureEnv, value: instrumentationConfig.OTELExporterOTLPMetricsInsecure},
		{name: traceEndpointEnv, value: instrumentationConfig.OTELExporterOTLPTracesEndpoint},
		{name: traceInsecureEnv, value: instrumentationConfig.OTELExporterOTLPTracesInsecure},
	}
	for _, setting := range settings {
		value := strings.TrimSpace(setting.value)
		if value == emptyString {
			continue
		}
		if err := os.Setenv(setting.name, value); err != nil {
			return fmt.Errorf(errorSetEnvironmentVariable, setting.name, err)
		}
	}
	return nil
}

func RecordDuration(ctx context.Context, operation string, startedAt time.Time) {
	if operationDuration == nil {
		return
	}
	operationDuration.Record(ctx, time.Since(startedAt).Seconds(), metric.WithAttributes(attribute.String(metricOperationAttribute, operation)))
}

func StartSpan(ctx context.Context, instrumentationName string, spanName string, attributes ...attribute.KeyValue) (context.Context, trace.Span) {
	return otel.Tracer(instrumentationName).Start(ctx, spanName, trace.WithAttributes(attributes...))
}

func FinishSpan(ctx context.Context, span trace.Span, operation string, startedAt time.Time, failed bool) {
	if failed {
		span.SetStatus(codes.Error, emptyString)
	}
	RecordDuration(ctx, operation, startedAt)
	span.End()
}
