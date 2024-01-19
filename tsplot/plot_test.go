package tsplot

import (
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
	"testing"
)

func Test_findMaxFromFloat64Data(t *testing.T) {
	type args []*monitoringpb.Point
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "find max value between two points",
			args: []*monitoringpb.Point{
				{
					Value: &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_DoubleValue{DoubleValue: float64(0)}},
				},
				{
					Value: &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_DoubleValue{DoubleValue: float64(1)}},
				},
			},
			want: float64(1),
		},
		{
			name: "find max value between multiple points",
			args: []*monitoringpb.Point{
				{
					Value: &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_DoubleValue{DoubleValue: float64(0)}},
				},
				{
					Value: &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_DoubleValue{DoubleValue: float64(1)}},
				},
				{
					Value: &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_DoubleValue{DoubleValue: float64(3)}},
				},
				{
					Value: &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_DoubleValue{DoubleValue: float64(4)}},
				},
			},
			want: float64(4),
		},
		{
			name: "find max value between multiple non whole numbers",
			args: []*monitoringpb.Point{
				{
					Value: &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_DoubleValue{DoubleValue: float64(0.1)}},
				},
				{
					Value: &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_DoubleValue{DoubleValue: float64(1.5)}},
				},
				{
					Value: &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_DoubleValue{DoubleValue: float64(3.8)}},
				},
				{
					Value: &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_DoubleValue{DoubleValue: float64(0.4)}},
				},
			},
			want: float64(3.8),
		},
		{
			name: "find max value between multiple points of the same value",
			args: []*monitoringpb.Point{
				{
					Value: &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_DoubleValue{DoubleValue: float64(1)}},
				},
				{
					Value: &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_DoubleValue{DoubleValue: float64(1)}},
				},
				{
					Value: &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_DoubleValue{DoubleValue: float64(1)}},
				},
				{
					Value: &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_DoubleValue{DoubleValue: float64(1)}},
				},
			},
			want: float64(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findMaxFromFloat64Data(tt.args); got != tt.want {
				t.Errorf("findMaxFromFloat64Data() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findMaxFromInt64Data(t *testing.T) {
	type args []*monitoringpb.Point
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "find max value between two points",
			args: []*monitoringpb.Point{
				{
					Value: &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_Int64Value{Int64Value: int64(0)}},
				},
				{
					Value: &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_Int64Value{Int64Value: int64(1)}},
				},
			},
			want: float64(1),
		},
		{
			name: "find max value between multiple points",
			args: []*monitoringpb.Point{
				{
					Value: &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_Int64Value{Int64Value: int64(0)}},
				},
				{
					Value: &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_Int64Value{Int64Value: int64(1)}},
				},
				{
					Value: &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_Int64Value{Int64Value: int64(3)}},
				},
				{
					Value: &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_Int64Value{Int64Value: int64(4)}},
				},
			},
			want: float64(4),
		},
		{
			name: "find max value between multiple points of the same value",
			args: []*monitoringpb.Point{
				{
					Value: &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_Int64Value{Int64Value: int64(1)}},
				},
				{
					Value: &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_Int64Value{Int64Value: int64(1)}},
				},
				{
					Value: &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_Int64Value{Int64Value: int64(1)}},
				},
				{
					Value: &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_Int64Value{Int64Value: int64(1)}},
				},
			},
			want: float64(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findMaxFromInt64Data(tt.args); got != tt.want {
				t.Errorf("findMaxFromInt64Data() = %v, want %v", got, tt.want)
			}
		})
	}
}
