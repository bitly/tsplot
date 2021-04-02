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

// MetricQuery is a type that encapsulates the various parts that are used to form a ListTimeSeriesRequest.
// If EndTime is not provided, it will be defaulted to the current time.
//
// Required Fields:
// Project
// TimeSeries
// StartTime
//
// When the request is built, it will be stored in the ListTimeSeries request field.
type MetricQuery struct {
	Project    string
	TimeSeries string
	StartTime  *time.Time
	EndTime    *time.Time

	ListTimeSeriesRequest *monitoringpb.ListTimeSeriesRequest
}

// BuildRequest will build a ListTimeSeriesRequest from the fields in MetricQuery.
// The resulting request, when performed, will result in either zero or one time series
// being returned to the caller. This is due to the Aggregation settings aligning on
// rate and reducing on mean. An error will be returned if MetricQuery is missing key data
// required to build the request.
//
// Required data:
// MetricQuery.Project
// MetricQuery.TimeSeries
// MetricQuery.StartTime
//
// If MetricQuery.EndTime has not been provided, it will default to time.Now()
//
func (mq *MetricQuery) BuildRequest() error {

	var tsreq monitoringpb.ListTimeSeriesRequest

	if mq.Project == "" {
		return errors.New("MetricQuery missing GCE Project")
	}

	if mq.TimeSeries == "" {
		return errors.New("MetricQuery missing TimeSeries")
	}

	if mq.StartTime == nil {
		return errors.New("start time has not been provided")
	}

	now := time.Now()
	if mq.EndTime == nil {
		mq.EndTime = &now
	}

	tsreq = monitoringpb.ListTimeSeriesRequest{
		Name:   fmt.Sprintf("projects/%s", mq.Project),
		Filter: fmt.Sprintf("resource.type = \"global\" AND metric.type = \"%s\"", mq.TimeSeries),
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

	mq.ListTimeSeriesRequest = &tsreq
	return nil
}

// Perform sends the MetricQuery.ListTimeSeriesRequest to the Google Cloud Monitoring API.
// If the request has not been built yet, i.e: BuildRequest() has not been called on the MetricQuery,
// an error will be returned. A Google Cloud Monitoring client is required to be passed in as a parameter
// if authentication has not been set up on the client, an error will result from the call.
func (mq *MetricQuery) Perform(client *monitoring.MetricClient) (*monitoring.TimeSeriesIterator, error) {
	if mq.ListTimeSeriesRequest == nil {
		return nil, errors.New("attempted to call Perform() with nil request")
	}
	return client.ListTimeSeries(context.Background(), mq.ListTimeSeriesRequest), nil
}
