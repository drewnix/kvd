
## Basic Questions

* Is this meant to be a distributed KV store or just a single node key value store? (i.e.
  should keys and values be distributed across multiple stores)?

* Does the server need to support multiple clients or will all interaction come through
  CLI?
    * Any backend requirements, REST API, GRPC, GraphQL? etc

* Is the CLI expected to work across the network?

* Are there any constraints as far as the size of keys & values?

* Are keys required to be unique/distinct?

* Durability should we support persisting the key-value store from memory to disk? (and
  reloading from disk on startup?)

* What is projected volume, number of users, any desire number of requests/sec?

* Is there an expectation as far as response time?

* Any requirements around where this will be hosted (cloud, kubernetes)?
    * Should I provide a docker and or deployment scripts?

* Operations on different keys should have minimal contention
    * Only get lock if operation is on the same key?

## Advanced Questions

* Isolation of transactions?

* Overwhelming majority are reads-is running multiple instances w/ partitioning desired?

* Does the server need to support versioning (i.e. update and rollback of data stored in particular keys)

* Any requirements as far as handling failures
    * What if system is OOM

* Any tradeoffs as far as CAP theorem (consistency, availability, partitioning)
