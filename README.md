# Mailu statistics server

Mail statistics are generated while avoiding all uses of intrusive trackers,
especially to protect users privacy and avoid collecting information that
could help build a comprehensive list of all mailu setups.

To accomplish this, Mailu uses DNS and relies on the fact that the vast
majority of servers do not run a local resolver but rely on recursors instead,
thus masking the true originating address from a query.

Unless the admin disables statisics upload in the configuration, Mailu
instances run a DNS query (try and download a HTTP URI) for a name composed
as follows:

```
<version>.<instanceid>.stats.mailu.io
```

The version is the Mailu version string. The instanceid is a unique id
persisted per instance, stored in a file next to the main database. instance
ids are UUIDv4.

The statistics server stores the timestamp, version and instance id
everytime it receives a request.

### Running

#### Docker

If you use Docker, getting it running to serve `stats.mailu.io` is very easy:

```
docker build -t mailu-stats --rm .
docker run -p 53:53/udp -it --rm --name mailu-stats mailu-stats
```

All environment-variables listed below can also be set in Docker via the `-e/--env` flag:
```
docker run -p 53:53/udp -it -e DOMAIN=stats.mydomain.com --rm --name mailu-stats mailu-stats
```

#### Manual

If you want to build the service manually you should set some environment-variables, here is a list of all, their defaults and what's their purpose:

|variable|default value|purpose|
|---|---|---|
|LOGPATH|`/output.log`|where to store the log|
|HOST|`0.0.0.0`|which IP or Hostname to listen on (supports IPv4 and IPv6)|
|PORT|`53`|which port to listen on|
|DOMAIN|`stats.mailu.io.`|which domain we serve. MUST end with a dot!|
|VALUECOUNT|`2`|how many subdomains (=values) we want|

And then you can simply build and run it with:

```
go build -v
./Stats
```
