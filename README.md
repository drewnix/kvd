
# KVD

KVD is a in-memory key-value store and CLI client.

Major features:
* REST API offering single key access or multiple key operations following 
  a straightforward API.
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
[1] 88547
Kvd started. Press ctrl-c to stop.                                                                                                                                                                                                      

$ ./kv set cake=🎂 horns=😈 smiley=😁
Keys set

$ ./kv get cake horns smiley
cake: 🎂
horns: 😈
smiley: 😁

$ ./kv metrics
Keys Stored: 3
Bytes Stored (Values): 60
Set Operations: 3
Get Operations: 3
Delete Operations: 0

$ ./kv del horns

$ ./kv metrics
Keys Stored: 2
Bytes Stored (Values): 40
Set Operations: 3
Get Operations: 3
Delete Operations: 1

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
* [ ] Operations on different keys should have minimal contention.
  * Partially complete - A BulkSet (which is default on multiple key update) will share a
    single mutex lock.
* [ ] Support updating all keys or none if they can't all be updated (no partial updates)
  * Could use rollback (described under "Thoughts for future improvement") 
* [ ] Highly concurrent
  * Could use fasthttp based (described under "Thoughts for future improvement") 
* [ ] Optimize for read-heavy workflow

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

## Thoughts for future improvement

* More Cleanup and Test Cases:
  * Effort: S
  * Better error handling
    * Print deleted keys
    * Error handling showing key
  * Add better logging on server
  * Finished Planned Test Cases

* For performance improvement: 
  * Effort: M 
  * Create benchmark to drive performance improvements using `go test -bench`
  * Switch API to using fasthttp based HTTP implementation (perhaps fiber framework)
    * Should improve on "highly concurrent" requirement

* Transaction support:
  * Effort: M 
  * Perform mutex locking by transaction - write transaction vs. read transaction
  * Keep audit log of old values and support rollback functionality on failure
  * Store commit records and do basic validation that commits can be done before committing

* Infrastructure and misc:
  * Effort: S
  * Support reading from file based config using Viper, remove references to localhost
  * Docker and kubernetes deployment?
  * Support maximum size limit to 

* TTL LRU caching
  * Effort: S
  * Add TTL tracking and eviction for records older then TTL