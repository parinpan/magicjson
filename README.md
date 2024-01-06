This Readme is generated by ChatGPT

---

# MagicJSON

MagicJSON is a Golang package designed to provide marshaling of private struct fields. This library allows you to overcome the limitations of traditional JSON marshaling by employing reflection and custom parsing strategies. Unlock the full potential of your data structures while maintaining encapsulation.

## Why MagicJSON?

Golang's JSON marshaling typically requires struct fields to be exported, revealing implementation details. MagicJSON breaks this limitation, enabling you to marshal private fields without compromising the integrity of your struct's encapsulation. This is particularly useful when you want to hide implementation details.

## Features

- **Private Field Marshaling:** Marshal private struct fields without exposing them.
- **Reflection-based Parsing:** Leverage reflection to dynamically explore and parse complex structures.
- **Customizable Parsing:** Implement custom parsing strategies for specific types or structures.

## Disclaimer

**Experimental and Learning Purposes Only:** MagicJSON is designed for experimentation and learning purposes. While it showcases innovative techniques for private field marshaling, it may not be suitable for production use. Use it with caution and understand the potential risks before incorporating it into critical systems.

## Limitations

- **Does Not Read JSON Struct Tags for Private Fields:** MagicJSON relies on reflection and does not read JSON struct tags for private fields. This means that custom JSON struct tags on private fields will not be considered during marshaling.

- **Fiields with MarshalJSON as Pointers:** To leverage custom JSON marshaling for private fields, they must be defined as pointers in the struct. This ensures that the library can access and marshal the respective private fields correctly. Example:

```go
// MagicJSON can only marshal dateB field (time.Time implements MarshalJSON internally)
type yourStruct struct {
    dateA time.Time
    dateB *time.Time
}
```

## Installation

```bash
go get -u github.com/parinpan/magicjson
```

## Usage

### Importing the Package

```go
import "github.com/parinpan/magicjson"
```

### Example Usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/parinpan/magicjson"
)

type ExampleStruct struct {
	// your private fields here
	a int
	b string
	c *time.Time
}

func main() {
	t := time.Now()
	
	obj := ExampleStruct{
		a: 10,
		b: "hello, world",
		c: &t,
	}

	payload, err := magicjson.Marshal(obj)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// output: {"a": 10, "b": "hello, world", "c": "2024-01-06 14:28:06.094184"}
	fmt.Println("JSON Payload:", string(payload))
}
```

Please see `json_test.go` for detailed example.

## How It Works

MagicJSON employs a combination of reflection and custom parsing strategies to marshal private fields. The `magicjson.MarshalJSON` function dynamically explores the structure of the provided object and uses custom parsers to handle different types.

## Custom Parsing Strategies

MagicJSON allows you to define custom parsing strategies for specific types or structures. By implementing the `MarshalJSON` method on your struct or type, you can control how it is marshaled.

```go
func (e ExampleStruct) MarshalJSON() ([]byte, error) {
    // Implement your custom marshaling logic here
}
```
