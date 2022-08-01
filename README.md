
# KVD

KVD is a in-memory key-value store and CLI client.

Major features:
* REST API offering single key or multiple key operations following a 
  straightforward API.
* CLI `kv` built using Cobra CLI framework offering the ability to set, 
  get, or delete one or many keys in the key store.
* Metrics subsystem presents: keys stored, # set ops, # get ops, # delete 
  ops, and size of total data stored in the key, value store.

To Build:

```bash
make
```

Here is an example of basic usage

```bash
$ ./kv serve &
[1] 79789
Kvd started. Press ctrl-c to stop.   

$ ./kv set cake=🎂 horns=😈 smiley=😁
Keys set

$ ./kv get cake horns smiley         
cake: 🎂
horns: 😈
smiley: 😁

$ ./kv metrics
Keys Stored: 0
Set Operations: 3
Get Operations: 3
Delete Operations: 3

```

Tests pass, but currently service needs to be running.

```bash
$ ./kv serve &
$ make test
```

### Deliverables:
* In-memory KV Store API Service
* Command Line Client to interact with Service
* Tests
* README containing design and instructions on how to run

### Functional Requirements

##### Satisfied:

* [X] Keys can be any string, values any binary data.
* [X] Key Value code should be original
* [X] CLI can be used to get or set the value for a single key from the server.
* [X] CLI supports setting multiple (key, values) in a single call
* [X] CLI can obtain values for multiple keys in a single call
* [X] CLI can be used to delete one or more keys in a single call
* [X] CLI can be used to obtain metrics data from the server, including:
  * [X] Total number of keys stored
  * [X] Total size of all values
  * [X] Total number of get, set, and delete operations on keys

##### Incomplete:

* [ ] Support updating all keys or none if they can't all be updated (no partial updates)
* [ ] Highly concurrent
* [ ] Operations on different keys should have minimal contention.
* [ ] Optimize for read-heavy workflow

### TODO

* [ ] Add accounting for deletion of values
* [ ] If deletion on invalid key return ErrInvalidKey
* [ ] Stats - Track size of all values
  * [ ] On load track size of added?
  * [ ] On delete remove size of added?
* [ ] Fix metrics 'Keys Stored' issue, currently not updating
* [ ] Fix delete issue and add test case
* [ ] Add command helpers for easier testing
* [ ] Start error handling, for example get when key is not there.
* [ ] Add better error handling
* [ ] Add better error handling in = split code in set
* [ ] Add support for transactions
* [ ] When serve is called, check if service is already running on port, if it is, don't try and start a new one
* [ ] Stats - Track Keys stored / loaded
  * [ ] Handle deletion case
* [ ] Stats - Track total number of operations done

### Test Cases to Implement
* [ ] Test invalid input to set
* [ ] Delete non-existent key
* [ ] Get non-existent key
* [ ] Set non-existent key
* [ ] Set existing (overwrite) key
* [ ] API get single key
* [ ] API get multiple keys
* [ ] API set single key
* [ ] API set multiple keys
* [ ] API del single key
* [ ] API del multiple keys
* [ ] Basic Metrics test
* [ ] Metrics add key / delete key test
* [ ] Serve test
* [ ] Serve daemonize test
* [ ] Serve shutdown test

### Done
* [X] Write README.md containing basic instructions on how to run <-- required
* [X] Modify cli get multi-key
* [X] Modify cli put multi-key
* [X] Supper multi-key in service
* [X] Implement cli get 1 key
* [X] Implement cli put 1 key
* [X] Implement cli delete 1 key
* [X] Build basic stats / metrics subsystem
* [X] Create command for reading stats
* [X] Refactor DB subsystem into db.go
* [X] Create Deletion Test
* [X] Cleanup cli delete multi-key
* [X] Cleanup cli set multi-key
* [X] Cleanup cli get multi-key
* [X] make output from 'kv metrics' nicer - also stop printing http code
* [X] Implement metrics system
* [X] Refactor clean up of set, and create get and get testcase
* [X] Clean up all the extra printing from set command
* [X] Stats command: Pretty print metrics from server


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

### A theoretical business case for our KV Store

### Ideas for future improvement

* Horizontal scaling using partitioning
- Effort: S

* Building a distributed KV store
- Effort: XL

* Implement key versioning (update & rollback)

* Option to evict records older then TTL
- Effort: S