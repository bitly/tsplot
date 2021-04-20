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

// DefaultQuery the query filter used if no overrides are given.
const DefaultQueryFilter = "resource.type = \"global\" AND metric.type = \"%s\""

// MetricQuery is a type that encapsulates the various parts that are used to form a ListTimeSeriesRequest.
// If EndTime is not provided, it will be defaulted to the current time.
//
// Required Fields:
// Project
// MetricDescriptor
// StartTime
type MetricQuery struct {
	Project          string
	MetricDescriptor string
	StartTime        *time.Time
	EndTime          *time.Time

	queryFilter string
}

// setQueryFilter provides a hook to modify the metric query filter.
func (mq *MetricQuery) setQueryFilter(queryFilter string) {
	mq.queryFilter = queryFilter
}

// request builds and returns a *monitoringpb.ListTimeSeriesRequeset.
// If there is not enough information to build the request an error is returned.
func (mq *MetricQuery) request() (*monitoringpb.ListTimeSeriesRequest, error) {

	var tsreq monitoringpb.ListTimeSeriesRequest

	if mq.Project == "" {
		return nil, errors.New("MetricQuery missing GCE Project")
	}

	if mq.MetricDescriptor == "" {
		return nil, errors.New("MetricQuery missing MetricDescriptor")
	}

	if mq.StartTime == nil {
		return nil, errors.New("start time has not been provided")
	}

	now := time.Now()
	if mq.EndTime == nil {
		mq.EndTime = &now
	}

	if mq.queryFilter == "" {
		mq.queryFilter = fmt.Sprintf(DefaultQueryFilter, mq.MetricDescriptor)
	}

	tsreq = monitoringpb.ListTimeSeriesRequest{
		Name:   fmt.Sprintf("projects/%s", mq.Project),
		Filter: mq.queryFilter,
		Interval: &monitoringpb.TimeInterval{
			EndTime:   timestamppb.New(*mq.EndTime),
			StartTime: timestamppb.New(*mq.StartTime),
		},
		Aggregation: &monitoringpb.Aggregation{
			AlignmentPeriod:    durationpb.New(time.Minute * 1),
			PerSeriesAligner:   monitoringpb.Aggregation_ALIGN_RATE,
			CrossSeriesReducer: monitoringpb.Aggregation_REDUCE_MEAN,
		},
		View: monitoringpb.ListTimeSeriesRequest_FULL,
	}

	return &tsreq, nil
}

// PerformWithClient sends the MetricQuery.ListTimeSeriesRequest to the Google Cloud Monitoring API.
// If the request has not been built yet, i.e: BuildRequest() has not been called on the MetricQuery,
// an error will be returned. A Google Cloud Monitoring client is required to be passed in as a parameter
// if authentication has not been set up on the client, an error will result from the call.
func (mq *MetricQuery) PerformWithClient(client *monitoring.MetricClient) (*monitoring.TimeSeriesIterator, error) {
	request, err := mq.request()
	if err != nil {
		return nil, err
	}
	return client.ListTimeSeries(context.Background(), request), nil
}
