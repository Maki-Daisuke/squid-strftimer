Squid-strftimer
===============

A TCP daemon receiving Squid logs, replacing timestamps in the logs using `strftime` format
(cf. replacing `1568175398.603` with `2019-09-11T04:16:38.603000Z`), and then writing the result 
to a file or STDOUT


How to Install
--------------

```
$ got get github.com/Maki-Daisuke/squid-strftimer
```

How to Use
----------

Start squid-strftimer daemon:

```
$ squid-strftimer &
```

Add the following line in your `squid.conf`:

```
access_log tcp://localhost:36059 squid
```

Restart Squid server. Then, Squid start to send access log to squid-strftimer daemon.


Use with Docker & AWS CloudWatch Logs
-------------------------------------

You can easily send access logs to AWS CloudWatch Logs with using Docker.

For example, you can write docker-compose.yml like this:

```
# docker-compose.yml
version: "3.7"
services:
  logdaemon:
    build:
      image: squid-strftimer
    logging:
      driver: awslogs
      options:
        awslogs-region: ap-northeast-1
        awslogs-group: access.log
        awslogs-stream: mystream
        awslogs-datetime-format: "%Y-%m-%dT%H:%M:%S.%fZ"
  squid:
    image: sameersbn/squid:3.5.27
    volumes:
      - type: bind
        source: ./squid.conf
        target: /etc/squid/squid.conf
```

And, add the following line in your squid.conf:

```
access_log tcp://logdaemon:36059 squid
```

You need to configure AWS credentials.
See [docker docs](https://docs.docker.com/config/containers/logging/awslogs/) for more details.


Command-Line Options
--------------------

```
Usage: squid-strftimer [-h] [--format FORMAT] [--port PORT] [FILE]
```

- `FORMAT`
  - Strftime-style format
  - default: `"%Y-%m-%dT%H:%M:%S.%fZ"`
- `PORT`
  - TCP port to listen
  - default: `36059`
- `FILE`
  - File to output logs
  - Use STDOUT if `"-"` is specified
  - default: `"-"`


Author
------

Daisuke (yet another) Maki
