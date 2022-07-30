
# KVStore

## Assignment

Exercise: Key-Value Server and Client

Write an API server that acts as an in-memory key-value store, and a command line
client to interact with the server. Keys may be any string; values may be any 
binary data.

Implement the server to allow as much concurrency as possible; operations on 
different keys should have minimal contention. You may assume that the overwhelming
majority of operations are to get values for existing keys.

### Functional Requirements

A user should be able to perform the following actions:

* The command line client (cli) can be used to get or set the value for a single key 
  from the server.
* The cli can be used to set the values for multiple keys in a single call, such that 
  either all of the keys are updated or none of them are.
    * Needs support for basic transactions?
* The cli can be used to obtain the values for multiple keys in a single call.
* The cli can be used to delete one or more keys in a single call.
* The cli can be used to obtain metrics data from the server, including:
    * the total number of keys stored
    * the total size of all values
    * the total number of get, set, and delete operations on keys
* You should write the code to manage the keys & values yourself, but you may use any 
  other libraries or packages.

### Non-Functional Requirements

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
## Design

* High speed REST API 
** Use Fiber
*** Fiber boilerplate: https://github.com/gofiber/boilerplate

* Makefile for build instructions

* Server <server.go>
** Loads config
** Loads DB
** Initializes hash
** handles API requests

* Transactions <tx.go>
** A simple transaction model will be created 
** Perhaps inspired by https://github.com/arriqaaq/flashdb/blob/main/txn.go
*** read-only
*** read-write

* DB <db.go>
** In memory KV store / hash
** Interface for all db operations
** Set / Get / Del

## Plan


## A theoretical business case for our KV Store

## Ideas for future improvement

* 

* Horizontal scaling using partitioning
- Effort: S

* Building a distributed KV store
- Effort: XL

* Implement key versioning (update & rollback)

* Option to evict records older then TTL
- Effort: S