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
ids are UUIDs.

The statistics server stores the date and time, version and instance id
everytime it receives a request.
