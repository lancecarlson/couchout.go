# Usage

This will stream all of the docs into Redis as Base64 format. See sister project, [couchin](https://github.com/lancecarlson/couchin.go) for importing documents back into couch.

```
couchout --url http://localhost:5984/db/_all_docs?include_docs=true | redis-cli
```

# Install

```
git clone git@github.com:lancecarlson/couchout.go.git
cd couchout.go
go build -o couchout # Builds a binary file for you. Put this in one of your PATH directories
```

# Why?

Redis works as a nice conduit for manipulating large datasets from couch. My approach is as follows:

* Define a view or use _all_docs to fetch all of the documents I want to modify
* Use couchout to export the data to Redis in Base64 format
* Make a node script (js has nice object merging syntax) that grabs the updated data set from somewhere
* Use node script to loop through new data set, call GET key in Redis for each updated doc, decode base64 value, apply merge changes, accumulate a bulk number of docs (maybe 100-1000), bulk save to couch, then DEL key from Redis

(optional steps for flushing documents that shouldn't be in your dataset anymore)

* Use node script to loop through the rest of the keys stored on Redis (KEYS *), accumulate a bulk number of docs, use the bulk save command and pass a json object back to couch {_id: idFromRedis, _deleted: true}

* Run FLUSHDB on redis.

# Future

I've had a lot of ideas for this. 

* Perhaps there could be other backends that are leveraged, but Redis is nice. 
* Couchin project 
* Port this to node js? 
* It might also be a good idea to just have Go store the data in memory and then hit a custom js file of your choosing and map over it or something. 

Anyway, feedback/patches appreciated.
