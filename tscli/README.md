# TSCLI
A CLI front-end to the libbitly/tsplot package.

```
$ tscli -h

_____________________________ .____    .___ 
\__    ___/   _____/\_   ___ \|    |   |   |
  |    |  \_____  \ /    \  \/|    |   |   |
  |    |  /        \\     \___|    |___|   |
  |____| /_______  / \______  /_______ \___|
                 \/         \/        \/    
A CLI front-end to the tsplot package which provides a method of plotting time series data taken
from Google Cloud Monitoring (formerly StackDriver).

Usage:
  tscli [flags]

Flags:
  -a, --app string              The (Bitly) application. Usually top level directory
      --end string              End of the time window for which the query returns time series data for. Hours or minutes accepted, i.e: -5h or -5m or now. (default "now")
  -h, --help                    help for tscli
  -m, --metric string           The metric.
  -o, --output string           Specify output directory for resulting plot. Defaults to current working directory.
      --print-raw               Only print time series data and exit.
  -p, --project string          GCP Project.
      --query-override string   Override the default query. Must be a full valid query. Metric flag is not used.
      --reduce                  Use a time series reducer to return a single averaged result.
  -s, --service string          The (Bitly) service. Service directory found under application directory.
      --start string            Start time of window for which the query returns time series data for. Hours or minutes accepted, i.e: -5h or -5m.

```

## Authentication
Authentication is done via `GOOGLE_APPLICATION_CREDENTIALS`. Ensure this variable is exported
in your environment and that it points to a valid GCP Service Account JSON file.

Otherwise, you will encounter this error:
```
err creating metric client google: could not find default credentials. See https://developers.google.com/accounts/docs/application-default-credentials for more information.
exit status 1
```
