[![GoDoc](https://godoc.org/github.com/romanyx/jwalk?status.svg)](https://godoc.org/github.com/romanyx/jwalk)
[![Go Report Card](https://goreportcard.com/badge/github.com/romanyx/jwalk)](https://goreportcard.com/report/github.com/romanyx/jwalk)
[![Build Status](https://travis-ci.org/romanyx/jwalk.svg?branch=master)](https://travis-ci.org/romanyx/jwalk)

# jwalk

Jwalk is builded un top of easyjson/jlexer, and allows to easely unmarshal any JSON input with arbitrary key names by walking through it.

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/romanyx/jwalk"
)

const input = `{
	"key1": [{
		"id": 1
	}, {
		"id": 2
	}],
	"key2": [1,2],
	"key3": {
		"id": 1
		"name": null,
	}
}`

func main() {
	i, err := jwalk.Parse([]byte(input))
	if err != nil {
		log.Fatal(err)
	}

	switch v := i.(type) {
	case jwalk.ObjectWalker:
		v.Walk(func(key string, value interface{}) {
			fmt.Println(key + ":")
			switch v := value.(type) {
			case jwalk.ObjectsWalker:
				v.Walk(func(obj jwalk.ObjectWalker) {
					fmt.Println("\t-")
					obj.Walk(func(key string, value interface{}) {
						if v, ok := value.(jwalk.Value); ok {
							fmt.Println("\t", key+":", v.Interface())
						}
					})
				})
			case jwalk.Value:
				fmt.Println("\t", v.Interface())
			case jwalk.ObjectWalker:
				v.Walk(func(key string, value interface{}) {
					if v, ok := value.(jwalk.Value); ok {
						fmt.Println("\t", key+":", v.Interface())
					}
				})
			}
		})
	}
}
```

## Testing

```bash
go test
```

## Contributing

Please feel free to submit issues, fork the repository and send pull requests!