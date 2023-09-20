# 📜 QuickSQL

[![GoDoc](https://godoc.org/github.com/NublyBR/go-quicksql?status.png)](http://godoc.org/github.com/NublyBR/go-quicksql)
[![Go Report Card](https://goreportcard.com/badge/github.com/NublyBR/go-quicksql)](https://goreportcard.com/report/github.com/NublyBR/go-quicksql)

A set of Go utility tools for writing SQL quickly.

# ⚡️ Basic Usage

Create insert statements with specified columns

```go
ins := NewInsert(writer, "table_name", "id", "name")
ins.Add(1, "QuickSQL")
ins.Add(2, "Hello, World!")
ins.Flush()

// INSERT INTO `table_name` (`id`, `name`) VALUES
//     (1, "QuickSQL"),
//     (2, "Hello, World!");
```

Create insert statements with columns from struct

```go
type demo struct {
    ID   int
    Name string
    Data []byte
}

ins := NewInsert(writer, "demo", demo{})
ins.Add(demo{
    ID:   1,
    Name: "QuickSQL",
    Data: "Hello, World!",
})
ins.Add(demo{
    ID:   2,
    Name: "Another Row",
    Data: "Hello, World!",
})
ins.Flush()

// INSERT INTO `demo` (`id`, `name`, `data`) VALUES
//     (1, "QuickSQL", 0x48656c6c6f2c20576f726c6421),
//     (2, "Another Row", 0x48656c6c6f2c20576f726c6421);
```

Split every `n` rows

```go
ins := NewInsert(buf, "split", "number").Every(2)

for i := 0; i < 10; i++ {
    ins.Add(i)
}

ins.Flush()

// INSERT INTO `split` (`number`) VALUES
// 	(0),
// 	(1);
//
// INSERT INTO `split` (`number`) VALUES
// 	(2),
// 	(3);
```