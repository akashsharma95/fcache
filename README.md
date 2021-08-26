Inmemcache
====
simple in-memory cache with an HTTP interface

## Build
* Run docker build
```bash
docker build -t cache .
```

* Run the built image
```bash
docker run -p 4000:4000 --rm -it cache
```

## Usage
* Store key value in cache
```zsh
curl --location --request POST 'http://localhost:4000/key/key1' \
--data-raw 'value'
```

* Retrieve value by key from cache
```zsh
curl --location --request GET 'http://localhost:4000/key/key1'
```