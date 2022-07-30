
# kvd

### Assignment

Exercise: Key-Value Server and Client

Write an API server that acts as an in-memory key-value store, and a command line
client to interact with the server. Keys may be any string; values may be any 
binary data.

Implement the server to allow as much concurrency as possible; operations on 
different keys should have minimal contention. You may assume that the overwhelming
majority of operations are to get values for existing keys.


### Deliverables: 
* In-memory KV Store API Service
* Command Line Client to interact with Service
* Tests
* README containing design and instructions on how to run

### Functional Requirements

* Keys can be any string, values any binary data.

* Highly concurrent

* Operations on different keys should have minimal contention.

* Optimize for read-heavy workflow

* Key Value system should be original

* External libraries can be used

#### CLI

* CLI can be used to get or set the value for a single key from the server.

* CLI supports setting multiple (key, values) in a single call
  * Support updating all keys or none if they can't all be updated (no partial updates)

* CLI can obtain values for multiple keys in a single call

* CLI can be used to delete one or more keys in a single call

* CLI can be used to obtain metrics data from the server, including:
    * Total number of keys stored
    * Total size of all values
    * Total number of get, set, and delete operations on keys

## Design

### Server

#### High speed REST API 
* Use Fiber
* Fiber boilerplate: https://github.com/gofiber/boilerplate

#### CI / Infrastructure
* Makefile for build

#### Server <server.go>
* Loads config
* Loads DB
* Initializes hash
* handles API requests

#### Transactions <tx.go>
* A simple transaction model will be created 
* Perhaps inspired by https://github.com/arriqaaq/flashdb/blob/main/txn.go
  * read-only
  * read-write

#### DB <db.go>
* In memory KV store / hash
* Interface for all db operations
* Set / Get / Del


### CLI

* Add daemonization? https://developpaper.com/start-and-stop-operations-of-golang-daemon/

```bash
kvcli delete key1,key2,...,keyn
```


```bash
kvcli get key1,key2,key3
```


```bash
kvcli set key1=val,key2=123,key3=xyz
```

### TODO

* Build basic stats / metrics subsystem
* Create command for reading stats
* Refactor DB subsystem into db.go
* Read config from a file, no http://localhost
* Add support for transactions?
* Pretty print metrics from server
* Modify cli get multi-key
* Modify cli put multi-key
* Modify cli delete multi-key
* Implement stats system
* Stats - Track Keys stored / loaded
  * Handle deletion case
* Stats - Track size of all values
  * On load track size of added? 
  * On delete remove size of added?
* Stats - Track total number of operations done

### Done
* Implement cli get 1 key
* Implement cli put 1 key
* Implement cli delete 1 key
* Rename put to set

### A theoretical business case for our KV Store

### Ideas for future improvement

* Horizontal scaling using partitioning
- Effort: S

* Building a distributed KV store
- Effort: XL

* Implement key versioning (update & rollback)

* Option to evict records older then TTL
- Effort: S