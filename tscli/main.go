package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"github.com/bitly/tsplot/tsplot"
	"github.com/spf13/cobra"
	"gonum.org/v1/plot/vg"
	"google.golang.org/api/option"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const GAP = "GOOGLE_APPLICATION_CREDENTIALS"

var GoogleCloudMonitoringClient *monitoring.MetricClient

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "err %v\n", err)
		os.Exit(1)
	}
}

var (
	rootCmd = &cobra.Command{
		Use:   "tscli",
		Short: "Plots time series data from GCM.",
		Long: `
_____________________________ .____    .___ 
\__    ___/   _____/\_   ___ \|    |   |   |
  |    |  \_____  \ /    \  \/|    |   |   |
  |    |  /        \\     \___|    |___|   |
  |____| /_______  / \______  /_______ \___|
                 \/         \/        \/    
A CLI front-end to the tsplot package which provides a method of plotting time series data taken
from Google Cloud Monitoring (formerly StackDriver).
`,
		PreRunE: auth,
		RunE:    executeQuery,
	}

	project   string
	query     string
	startTime string
	endTime   string
	outDir    string
	title     string
	groupBy   string
)

func init() {
	rootCmd.Flags().StringVarP(&project, "project", "p", "", "GCP Project.")
	rootCmd.Flags().StringVarP(&query, "query", "m", "", "The query filter.")
	rootCmd.Flags().StringVar(&startTime, "start", "", "Start time of window for which the query returns time series data for. Hours or minutes accepted, i.e: -5h or -5m.")
	rootCmd.Flags().StringVar(&endTime, "end", "now", "End of the time window for which the query returns time series data for. Hours or minutes accepted, i.e: -5h or -5m or now.")
	rootCmd.Flags().StringVarP(&outDir, "output", "o", "", "Specify output directory for resulting plot. Defaults to current working directory.")
	rootCmd.Flags().StringVarP(&title, "title", "t", "", "Specify title of graph.")
	rootCmd.Flags().StringVar(&groupBy, "group-by", "", "Key to group metric by when dealing with multiple time series.")
	rootCmd.MarkFlagRequired("project")
	rootCmd.MarkFlagRequired("query")
	rootCmd.MarkFlagRequired("start")
	rootCmd.MarkFlagRequired("title")
	rootCmd.MarkFlagRequired("output")
}

func auth(cmd *cobra.Command, args []string) error {
	var opts []option.ClientOption
	serviceAccountJsonPath := os.Getenv(GAP)
	if serviceAccountJsonPath != "" {
		opts = append(opts, option.WithCredentialsFile(serviceAccountJsonPath))
	}

	mc, err := monitoring.NewMetricClient(context.Background(), opts...)
	if err != nil {
		return err
	}
	GoogleCloudMonitoringClient = mc
	return nil
}

func executeQuery(cmd *cobra.Command, args []string) error {

	if !timeFormatOK(startTime) {
		return errors.New("err validating start time format")
	}

	if endTime != "" && !timeFormatOK(endTime) {
		return errors.New("err validating end time format")
	}

	if startTime == endTime {
		return errors.New("err invalid time frame")
	}

	st := parseTime(startTime)
	et := parseTime(endTime)

	request := &monitoringpb.ListTimeSeriesRequest{
		Name: fmt.Sprintf("projects/%s", project),
		//`resource.type = "global" AND metric.type = "custom.googleapis.com/opencensus/fishnet/queuereader_fishnet/messages_total"`
		Filter: query,
		Interval: &monitoringpb.TimeInterval{
			EndTime:   timestamppb.New(et),
			StartTime: timestamppb.New(st),
		},
		Aggregation: &monitoringpb.Aggregation{
			AlignmentPeriod: durationpb.New(time.Minute * 1),
			// todo: these need to be settable as they are not uniformly useful across all metric types.
			PerSeriesAligner:   monitoringpb.Aggregation_ALIGN_RATE,
			CrossSeriesReducer: monitoringpb.Aggregation_REDUCE_MEAN,
			GroupByFields:      []string{fmt.Sprintf("metric.labels.%s", groupBy)},
		},
		View: monitoringpb.ListTimeSeriesRequest_FULL,
	}

	tsi := GoogleCloudMonitoringClient.ListTimeSeries(context.Background(), request)
	plot, err := tsplot.NewPlotFromTimeSeriesIterator(tsi, groupBy, tsplot.WithXTimeTicks(time.Kitchen), tsplot.WithTitle(title), tsplot.WithXAxisName("UTC"))
	if err != nil {
		log.Fatal(err)
	}

	saveFile := fmt.Sprintf("%s/%s.png", outDir, title)
	return plot.Save(8*vg.Inch, 4*vg.Inch, saveFile)
}

func timeFormatOK(s string) bool {
	if s == "" {
		return false
	}
	if s == "now" {
		return true
	}
	b, _ := regexp.MatchString(`^-\d+[m,h]$`, s)
	if !b {
		return false
	}
	return true
}

func parseTime(s string) time.Time {
	if s == "now" {
		return time.Now()
	}
	parts := strings.Split(s, "")
	num, _ := strconv.Atoi(parts[1 : len(parts)-1][0])
	unit := parts[len(parts)-1]
	now := time.Now()
	var t time.Time
	switch unit {
	case "m":
		t = now.Add(-time.Minute * time.Duration(num))
	case "h":
		t = now.Add(-time.Hour * time.Duration(num))
	}
	return t
}
