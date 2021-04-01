package tsplot

import (
	"testing"
)

func TestMetricQuery_BuildRequest(t *testing.T) {
	tests := []struct {
		desc      string
		in        *MetricQuery
		out       *MetricQuery
		expectErr bool
	}{
		{
			desc: "fail to  build request, missing Project",
			in: &MetricQuery{
				TimeSeries: "compute.googleapis.com/instance/cpu/usage_time",
			},
			out:       &MetricQuery{},
			expectErr: true,
		},
		{
			desc: "fail to  build request, missing TimeSeries",
			in: &MetricQuery{
				Project: "my-project",
			},
			out:       &MetricQuery{},
			expectErr: true,
		},
	}

	for _, test := range tests {
		err := test.in.BuildRequest()
		if err != nil && !test.expectErr {
			t.Errorf("got unexpected err: %v\n", err)
		}
	}
}
