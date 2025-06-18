package tracestore

import (
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseDynamicKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Multiple Keys",
			input:    "{{k8s.cluster.name}}-{{k8s.namespace.name}}",
			expected: []string{"k8s.cluster.name", "k8s.namespace.name"},
		},
		{
			name:     "Multiple same Keys",
			input:    "{{k8s.cluster.name}}-{{k8s.cluster.name}}",
			expected: []string{"k8s.cluster.name"},
		},
		{
			name:     "Single Key",
			input:    "{{k8s.cluster.name}}",
			expected: []string{"k8s.cluster.name"},
		},
		{
			name:     "No Keys",
			input:    "cluster-1",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseDynamicKeys(tt.input)
			require.Equal(t, tt.expected, got)
		})
	}
}

func TestBuildIndexSuffix(t *testing.T) {
	tests := []struct {
		name        string
		dynamicKeys []string
		template    string
		indexSuffix string
	}{
		{
			name:        "No Keys",
			dynamicKeys: nil,
			template:    "cluster",
			indexSuffix: "cluster",
		},
		{
			name:        "Single Key",
			dynamicKeys: []string{"k8s.cluster.name"},
			template:    "{{k8s.cluster.name}}",
			indexSuffix: "cluster-1",
		},
		{
			name:        "Multiple Keys",
			dynamicKeys: []string{"k8s.cluster.name", "k8s.namespace.name"},
			template:    "{{k8s.cluster.name}}-{{k8s.namespace.name}}",
			indexSuffix: "cluster-1-namespace-1",
		},
		{
			name:        "Multiple same Keys",
			dynamicKeys: []string{"k8s.cluster.name"},
			template:    "{{k8s.cluster.name}}-{{k8s.cluster.name}}",
			indexSuffix: "cluster-1-cluster-1",
		},
		{
			name:        "not exist Keys",
			dynamicKeys: []string{"not.exist.key1", "not.exist.key2"},
			template:    "${not.exist.key1}-${not.exist.key2}",
			indexSuffix: "unrouted-unrouted",
		},
		// TODO test sort,test attribute type not support
	}

	traces := ptrace.NewTraces()
	rSpans := traces.ResourceSpans().AppendEmpty()
	rSpans.Resource().Attributes().PutStr("k8s.cluster.name", "cluster-1")
	rSpans.Resource().Attributes().PutStr("k8s.namespace.name", "namespace-1")
	sSpans := rSpans.ScopeSpans().AppendEmpty()
	span := sSpans.Spans().AppendEmpty()

	spanID := pcommon.NewSpanIDEmpty()
	spanID[5] = 5 // 0000000000050000
	span.SetSpanID(spanID)

	traceID := pcommon.NewTraceIDEmpty()
	traceID[15] = 1 // 00000000000000000000000000000001
	span.SetTraceID(traceID)
	dbSpans := ToDBModel(traces)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := buildIndexSuffix(&dbSpans[0], tt.dynamicKeys, tt.template)
			assert.Equal(t, tt.indexSuffix, actual)
		})
	}

}
