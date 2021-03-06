# Usage

This will stream all of the docs into Redis as Base64 format. See sister project, [couchin](https://github.com/lancecarlson/couchin.go) for importing documents back into couch.

```
couchout http://localhost:5984/db/_all_docs?include_docs=true | redis-cli
```

Full usage with your own custom node script and couchin might look like this:

```
DB="http://localhost:5984/db"
redis-cli flushdb && couchout $DB/_all_docs?include_docs=true | redis-cli && node update.js && couchin $DB/_bulk_docs
```

# Install

```
git clone https://github.com/lancecarlson/couchout.go
cd couchout.go
go build -o couchout # Builds a binary file for you. Put this in one of your PATH directories
```

# Why?

Redis works as a nice conduit for manipulating large datasets from couch. By sticking your data in Redis, you get the ability to make language agnostic scripts and also separate out data manipulation into separate processes. Your only limitation is the memory on the box. 

My approach is as follows:

1. Define a view or use _all_docs to fetch all of the documents I want to modify. Make sure to set include_docs=true and if you have a dataset that is huge, you should consider specifying a limit to the query string. My experience has been that Redis is pretty efficient even with several million docs.
2. Use couchout to export the data to Redis in Base64 format
3. Make a node script (js has nice object merging syntax) that grabs the updated data set from somewhere
4. Use node script to loop through new data set, call GET key in Redis for each updated doc, decode base64 value, apply merge changes
5. Accumulate a bulk number of docs (maybe 100-1000), bulk save to couch, then DEL key from Redis (or FLUSHDB at the end) OR you can use [couchin](https://github.com/lancecarlson/couchin.go)

# Future

I've had a lot of ideas for this. 

* Perhaps there could be other backends that are leveraged, but Redis is nice. 
* Port this to node js?
* It might also be a good idea to just have Go store the data in memory and then hit a custom js file of your choosing and map over it or something. 

Anyway, feedback/patches appreciated.
