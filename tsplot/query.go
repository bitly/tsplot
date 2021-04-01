package tsplot

import (
	"context"
	"errors"
	"fmt"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MetricQuery struct {
	Project    string
	TimeSeries string // todo: add ability to compare multiple time series
	StartTime  time.Time
	EndTime    time.Time

	ListTimeSeriesRequest *monitoringpb.ListTimeSeriesRequest
}

func (mq *MetricQuery) BuildRequest() error {

	var tsreq monitoringpb.ListTimeSeriesRequest

	if mq.Project == "" {
		return errors.New("MetricQuery missing GCE Project")
	}

	if mq.TimeSeries == "" {
		return errors.New("MetricQuery missing TimeSeries")
	}

	tsreq = monitoringpb.ListTimeSeriesRequest{
		Name:   fmt.Sprintf("projects/%s", mq.Project),
		Filter: fmt.Sprintf("resource.type = \"global\" AND metric.type = \"%s\"", mq.TimeSeries),
		Interval: &monitoringpb.TimeInterval{
			EndTime:   timestamppb.New(mq.EndTime),
			StartTime: timestamppb.New(mq.StartTime),
		},
		Aggregation: &monitoringpb.Aggregation{
			AlignmentPeriod:    durationpb.New(time.Minute * 1),
			PerSeriesAligner:   monitoringpb.Aggregation_ALIGN_RATE,
			CrossSeriesReducer: monitoringpb.Aggregation_REDUCE_MEAN,
		},
		View: monitoringpb.ListTimeSeriesRequest_FULL,
	}

	mq.ListTimeSeriesRequest = &tsreq
	return nil
}

func (mq *MetricQuery) Perform(client *monitoring.MetricClient) (*monitoring.TimeSeriesIterator, error) {
	if mq.ListTimeSeriesRequest == nil {
		return nil, errors.New("attempted to call Perform() with nil request")
	}
	return client.ListTimeSeries(context.Background(), mq.ListTimeSeriesRequest), nil
}
