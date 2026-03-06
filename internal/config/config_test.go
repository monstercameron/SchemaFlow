package config

import "testing"

func TestGetTraceEnabledHonorsBothEnvNames(t *testing.T) {
	t.Setenv("SCHEMAFLOW_TRACE", "")
	t.Setenv("SCHEMAFLOW_ENABLE_TRACING", "")
	SetTraceEnabled(false)

	if GetTraceEnabled() {
		t.Fatal("expected tracing disabled when env vars are unset")
	}

	t.Setenv("SCHEMAFLOW_TRACE", "true")
	if !GetTraceEnabled() {
		t.Fatal("expected SCHEMAFLOW_TRACE=true to enable tracing")
	}

	t.Setenv("SCHEMAFLOW_TRACE", "")
	t.Setenv("SCHEMAFLOW_ENABLE_TRACING", "1")
	if !GetTraceEnabled() {
		t.Fatal("expected SCHEMAFLOW_ENABLE_TRACING=1 to enable tracing")
	}
}
