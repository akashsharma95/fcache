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

## Benchmark
```
goos: darwin
goarch: amd64
pkg: inmemcache/pkg/cache
cpu: Intel(R) Core(TM) i5-8279U CPU @ 2.40GHz
BenchmarkInMemoryCache_Set
BenchmarkInMemoryCache_Set-8      	     249	   4220088 ns/op	  15.53 MB/s	  332945 B/op	   65546 allocs/op
BenchmarkInMemoryCache_Get
BenchmarkInMemoryCache_Get-8      	     405	   2892615 ns/op	  22.66 MB/s	  306411 B/op	   65704 allocs/op
BenchmarkInMemoryCache_SetGet
BenchmarkInMemoryCache_SetGet-8   	     127	  10741089 ns/op	  12.20 MB/s	  663429 B/op	  131093 allocs/op
BenchmarkStdMap_Set
BenchmarkStdMap_Set-8             	     100	  10856491 ns/op	   6.04 MB/s	  363718 B/op	   65558 allocs/op
BenchmarkStdMap_Get
BenchmarkStdMap_Get-8             	     433	   2396273 ns/op	  27.35 MB/s	   24023 B/op	     156 allocs/op
BenchmarkStdMap_SetGet
BenchmarkStdMap_SetGet-8          	      84	  51293764 ns/op	   2.56 MB/s	  383134 B/op	   65565 allocs/op
BenchmarkSyncMap_Set
BenchmarkSyncMap_Set-8            	      61	  21072146 ns/op	   3.11 MB/s	 3560814 B/op	  264332 allocs/op
BenchmarkSyncMap_Get
BenchmarkSyncMap_Get-8            	    1376	    852513 ns/op	  76.87 MB/s	    9245 B/op	     287 allocs/op
BenchmarkSyncMap_SetGet
BenchmarkSyncMap_SetGet-8         	     187	   5710918 ns/op	  22.95 MB/s	 3457697 B/op	  262857 allocs/op
PASS
```