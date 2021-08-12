package tsplot

import (
	"testing"
	"time"

	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
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
				MetricDescriptor: "compute.googleapis.com/instance/cpu/usage_time",
			},
			out:       &MetricQuery{},
			expectErr: true,
		},
		{
			desc: "fail to  build request, missing MetricDescriptor",
			in: &MetricQuery{
				Project: "my-project",
			},
			out:       &MetricQuery{},
			expectErr: true,
		},
	}

	for _, test := range tests {
		_, err := test.in.request()
		if err != nil && !test.expectErr {
			t.Errorf("got unexpected err: %v\n", err)
		}
	}
}

func TestMetricQuery_SetQueryFilter(t *testing.T) {
	expectedFilter := "some advanced query"
	st := time.Now().Add(-1 * time.Hour)
	et := time.Now()
	query := NewMetricQuery("bitly-gcp-prod", "", &st, &et)
	query.SetQueryFilter(expectedFilter)

	tsr, err := query.request()
	if err != nil {
		t.Error(err)
	}

	filter := tsr.GetFilter()
	if filter != expectedFilter {
		t.Fatalf("query filter not overriden. got: %s, expected: %s", filter, expectedFilter)
	}
}

func TestMetricQuery_SetAlignmentPeriod(t *testing.T) {
	expectedAlignmentPeriod := time.Minute * 10
	st := time.Now().Add(-1 * time.Hour)
	et := time.Now()
	query := NewMetricQuery("bitly-gcp-prod", "some metric", &st, &et)
	query.SetAlignmentPeriod(expectedAlignmentPeriod)

	req, err := query.request()
	if err != nil {
		t.Error(err)
	}

	alignmentPeriod := req.GetAggregation().GetAlignmentPeriod().GetSeconds()
	if req.GetAggregation().GetAlignmentPeriod().GetSeconds() != int64(expectedAlignmentPeriod.Seconds()) {
		t.Fatalf("alignment period not overriden. got %d, expected: %s", alignmentPeriod, expectedAlignmentPeriod)
	}
}

func TestMetricQuery_AggregationOptions(t *testing.T) {
	st := time.Now().Add(-1 * time.Hour)
	et := time.Now()
	query := NewMetricQuery("bitly-gcp-prod", "some metric", &st, &et)
	query.Set_ALIGN_NONE()
	query.Set_REDUCE_NONE()

	req, err := query.request()
	if err != nil {
		t.Error(err)
	}

	aggregation := req.GetAggregation()
	if aggregation.GetPerSeriesAligner() != monitoringpb.Aggregation_ALIGN_NONE {
		t.Fatal("aligner not overridden")
	}
	if aggregation.GetCrossSeriesReducer() != monitoringpb.Aggregation_REDUCE_NONE {
		t.Fatal("reducer not overridden")
	}
}
