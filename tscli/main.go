package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"tsplot/tsplot"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"github.com/spf13/cobra"
	"gonum.org/v1/plot/vg"
	"google.golang.org/api/option"
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
	app       string
	service   string
	metric    string
	startTime string
	endTime   string
	justPrint bool
)

func init() {
	rootCmd.Flags().StringVarP(&project, "project", "p", "", "GCP Project.")
	rootCmd.Flags().StringVarP(&app, "app", "a", "", "The (Bitly) application. Usually top level directory")
	rootCmd.Flags().StringVarP(&service, "service", "s", "", "The (Bitly) service. Service directory found under application directory.")
	rootCmd.Flags().StringVarP(&metric, "metric", "m", "", "The metric.")
	rootCmd.Flags().StringVar(&startTime, "start", "", "Start time of window for which the query returns time series data for. Hours or minutes accepted, i.e: -5h or -5m.")
	rootCmd.Flags().StringVar(&endTime, "end", "now", "End of the time window for which the query returns time series data for. Hours or minutes accepted, i.e: -5h or -5m or now.")
	rootCmd.Flags().BoolVar(&justPrint, "print-raw", false, "Only print time series data and exit.")
	rootCmd.MarkFlagRequired("project")
	rootCmd.MarkFlagRequired("app")
	rootCmd.MarkFlagRequired("service")
	rootCmd.MarkFlagRequired("metric")
	rootCmd.MarkFlagRequired("start")
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

	query := tsplot.MetricQuery{
		Project:    project,
		TimeSeries: fmt.Sprintf("custom.googleapis.com/opencensus/%s/%s/%s", app, service, metric),
		StartTime:  st,
		EndTime:    et,
	}

	if err := query.BuildRequest(); err != nil {
		return err
	}
	tsi, err := query.Perform(GoogleCloudMonitoringClient)
	if err != nil {
		return err
	}

	// query.Perform() will only return a single time series.
	// no need to loop
	timeSeries, err := tsi.Next()
	if err != nil {
		return err
	}

	if justPrint {
		fmt.Printf("%v\n\n", timeSeries)
		return nil
	}

	// TODO: update to support multiple time series
	// tsplot pkg has been updated with support, now just need to do CLI
	ts := tsplot.TimeSeries{}
	ts[metric] = timeSeries.GetPoints()

	plot := tsplot.TimeSeriesPlot{
		Name:        timeSeries.GetMetric().GetType(),
		XAxisName:   "",
		YAxisName:   "",
		Description: "",
		TimeSeries:  ts,
	}

	timeSeriesPlot, err := plot.Create()
	if err != nil {
		return err
	}

	saveFile := fmt.Sprintf("%s-%s.png", service, metric)
	timeSeriesPlot.Save(8*vg.Inch, 4*vg.Inch, saveFile)

	return nil
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
