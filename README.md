# gomcketama [![Build Status](https://travis-ci.org/thatguystone/gomcketama.svg?branch=master)](https://travis-ci.org/thatguystone/gomcketama) [![GoDoc](https://godoc.org/github.com/thatguystone/gomcketama?status.svg)](https://godoc.org/github.com/thatguystone/gomcketama)

This package implements a ServerSelector for [gomemcache](https://github.com/bradfitz/gomemcache) that provides ketama hashing that's compatible with SpyMemcached's ketama hashing.

## Usage

```go
import "github.com/thatguystone/gomcketama"

// Create a memcache client using ketama as the server selector
mc := gomcketama.New("server1:11211", "server2:11211")
```
