# Go Bin Packer
A Go implementation of the bin packing algorithm with support for generic types.

## What is Bin Packing?
Bin packing is an optimization problem where the goal is to efficiently pack objects of different sizes into a finite number of bins or containers, each with a fixed capacity, while minimizing the number of bins used.

Common applications include:

* Shipping and logistics (fitting packages into containers)

* Memory allocation in computing

* Task scheduling on processors

* Cutting stock problems in manufacturing

## About This Project
This project provides a generic implementation of the First-Fit bin packing algorithm in Go. It leverages Go's generics to allow packing of any type of data.

## Key features:

* Generic implementation that works with any data type

* Observer pattern for monitoring bin packing events

* First-Fit algorithm implementation

```go

// Create a bin packer with bins of size 100
packer, _ := NewPacker[string](100)

// Create some blocks to pack
block1 := &Block[string]{
    Size: func() int { return 30 },
    Data: "Block 1",
}

block2 := &Block[string]{
    Size: func() int { return 50 },
    Data: "Block 2",
}

// add blocks
packer.Add(block1)
packer.Add(block2)

// packed bins
bins := packer.Bins()
```

## Event System
The library includes an observer pattern implementation that allows you to monitor bin packing events:
 
- BinCreated: When a new bin is created

- BlockAdded: When a block is successfully added to a bin

- BlockIgnored: When a block cannot be added to a bin

## Future Improvements

* Additional bin packing algorithms (Best-Fit, Next-Fit, etc.)

* Performance optimizations

* More comprehensive examples and documentation

## License

Apache 2.0
