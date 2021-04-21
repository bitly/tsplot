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
	"golang.org/x/image/colornames"
	"gonum.org/v1/plot/vg"
	"google.golang.org/api/iterator"
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

	project       string
	app           string
	service       string
	metric        string
	startTime     string
	endTime       string
	queryOverride string
	reduce        bool
	justPrint     bool
)

func init() {
	rootCmd.Flags().StringVarP(&project, "project", "p", "", "GCP Project.")
	rootCmd.Flags().StringVarP(&app, "app", "a", "", "The (Bitly) application. Usually top level directory")
	rootCmd.Flags().StringVarP(&service, "service", "s", "", "The (Bitly) service. Service directory found under application directory.")
	rootCmd.Flags().StringVarP(&metric, "metric", "m", "", "The metric.")
	rootCmd.Flags().StringVar(&startTime, "start", "", "Start time of window for which the query returns time series data for. Hours or minutes accepted, i.e: -5h or -5m.")
	rootCmd.Flags().StringVar(&endTime, "end", "now", "End of the time window for which the query returns time series data for. Hours or minutes accepted, i.e: -5h or -5m or now.")
	rootCmd.Flags().BoolVar(&justPrint, "print-raw", false, "Only print time series data and exit.")
	rootCmd.Flags().BoolVar(&reduce, "reduce", false, "Use a time series reducer to return a single averaged result.")
	rootCmd.Flags().StringVar(&queryOverride, "query-override", "", "Override the default query. Must be a full valid query. Metric flag is not used.")
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

	if metric != "" && queryOverride != "" {
		fmt.Println("warn: both --metric and --query-override flag used. Favoring --query-override.")
	}

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
		Project:          project,
		MetricDescriptor: fmt.Sprintf("custom.googleapis.com/opencensus/%s/%s/%s", app, service, metric),
		StartTime:        &st,
		EndTime:          &et,
	}

	if queryOverride != "" {
		query.SetQueryFilter(queryOverride)
	}

	query.SetReduce(reduce)

	tsi, err := query.PerformWithClient(GoogleCloudMonitoringClient)
	if err != nil {
		return err
	}

	ts := tsplot.TimeSeries{}
	for {
		timeSeries, err := tsi.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return err
		}

		// todo: implement "just-print" mode for multiple time series
		//if justPrint {
		//	fmt.Printf("%v", timeSeries)
		//	return nil
		//}

		// key helps to fill out legend.
		// Here we are grabbing the pod name.
		key := timeSeries.GetMetric().GetLabels()["opencensus_task"]
		if key == "" {
			// Labels we want to use don't necessarily exist when a cross series reducer has been used.
			// So we can just use "mean" in the legend.
			key = "mean"
		}
		ts[key] = timeSeries.GetPoints()
	}

	p, err := ts.Plot([]tsplot.PlotOption{tsplot.WithXAxisName("UTC"),
		tsplot.WithXTimeTicks(time.Kitchen),
		tsplot.WithGrid(colornames.Darkgrey),
		tsplot.WithTitle(metric)}...)
	if err != nil {
		return err
	}

	saveFile := fmt.Sprintf("%s-%s.png", service, metric)
	p.Save(8*vg.Inch, 4*vg.Inch, saveFile)

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
