# DevLog

A super simple non-persistent development log collector.

## Motivation

There are already some high quality log collectors out there: Elastic Search, GrayLog, Splunk, etc. However, there are many cases where you just want something simple to parse and watch your json logs so you can run it locally on your machine. This project aims to be just that: a bare-bones log collector that you can send your logs to. No databases, no persistence, no volumes. Just a docker container and a TCP port.

## Usage

To get started, you'll need docker, and then you can just run this command:
```
docker run --name devlog -p 9090-9091:9090-9091/tcp monodop/devlog
```

Next, configure your logger(s) of choice to send json logs to `devlog:9090` if you're running your code inside docker, or `localhost:9090` from your host OS.

To view logs, simply navigate your browser to `localhost:9091` and logs will begin streaming in as they are received. DevLog has no persistence, so you will only see a limited amount of messages whenever you refresh the page. Everything else is temporarily stored in your browser window.

## Filter syntax

In the DevLog UI, you can perform some basic types of filters to narrow in on what you want to see in your logs. You can combine any of the following filters; each one must match a log message for it to appear.

### Simple Keyword Filter

Simple keyword filters can be added to a query by just typing the keyword, surrounded by spaces, or by putting the search in double quotes. All keywords must match anywhere in the provided json. Example
```
keyword1 keyword2 "multiple length keyword"
```

### Key/Value Filter

A Key/Value filter works by specifying both the key in the json object, along with the value you want to search for. If you want to do a partial match, you can use `*` as a wildcard at the start or end of your search. Example:

```
app=myWebApplication logger=restapi.users.* message=*offline*
```

## Configuration

Configuration is done with environment variables on the container. Here's the full list of options:

(psych! there's nothing yet. coming soon)

## Contributing

If there's any features you'd like to see, or if you find any bugs, feel free to open an issue on the [GitHub issue tracker](https://github.com/monodop/devlog/issues), or open a pull request!