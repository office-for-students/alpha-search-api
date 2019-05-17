# alpha-search-api
==================
An app for accessing the Office for Students data (specifically institutions and courses)

### Installation

As the service are written in Go, make sure you have version 1.10.0 or greater installed.

Using [Homebrew](https://brew.sh/) to install go
* Run `brew install go` or `brew upgrade go`
* Set your `GOPATH` environment variable, this specifies the location of your workspace

#### Elasticsearch

It is expected that any minor version of elasticsearch 6 to be installed (e.g. 6.7.0)

* Run `brew install elasticsearch` - this will install latest version
* Run `brew services restart elasticsearch`

Set environment variable for elasticsearch uri either run the following code in terminal or add to `.bashrc`:

```
export ELASTIC_SEARCH_URI=<elasticsearch uri>
```
The elasticsearch uri should look something like: `localhost:9200`

#### Loading Data from scripts

You will need to install mongo db:

* Run `brew install mongodb`
* Run `brew services restart mongodb` 

Follow documentation in the [scripts repository](https://github.com/office-for-students/alpha-scripts)

#### Running Service

* Run `make debug`

#### Running tests

* Run `make test`

### Configuration

| Environment variable      | Default                | Description
| ------------------------- | ---------------------- | ----------------------------------------------------------------
| BIND_ADDR                 | :10100                 | The host and port to bind to
| DEFAULT_MAX_RESULTS       | 1000                   | The maximum number of results to be returned per page
| GRACEFUL_SHUTDOWN_TIMEOUT | 5s                     | The graceful shutdown timeout in seconds
| HOST_NAME                 | http://localhost       | The scheme and host name
| ES_DESTINATION_URL        | http://localhost:9200  | The address of the elasticsearch cluster
| ES_DESTINATION_INDEX      | courses                | The elasticsearch index in which the course data will be stored against
| ES_SHOW_SCORE             | false                  | A flag to return scores of course documents based on relevance. Should always be switched off in production environment


### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

See [LICENSE](LICENSE.md) for details.
