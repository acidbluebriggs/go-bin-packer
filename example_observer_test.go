/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package packer

import "fmt"

type BinEventLogger[T any] struct {
	Name string
}

func (l *BinEventLogger[T]) Update(event BinEvent[T]) {
	switch event.Type() {
	case BinCreated:
		fmt.Printf("[%s] Bin created: ID=%d\n", l.Name, event.Bin().id)

	case BlockAdded:
		if addEvent, ok := event.(BlockAddedEvent[T]); ok {
			fmt.Printf("[%s] Block added to bin %d, size=%d\n",
				l.Name, event.Bin().id, addEvent.Block().Size())
		}

	case BlockIgnored:
		if ignoreEvent, ok := event.(BlockIgnoredEvent[T]); ok {
			fmt.Printf("[%s] Block ignored for bin %d, size=%d\n",
				l.Name, event.Bin().id, ignoreEvent.Block().Size())
		}
	}
}

func ExampleBinPacker_Add() {
	// Create a bin packer
	packer := &BinPacker[string]{
		max:   100,
		count: 0,
	}

	packer.observers.observers = make([]Observer[string], 0)

	logger1 := &BinEventLogger[string]{Name: "Logger1"}
	logger2 := &BinEventLogger[string]{Name: "Logger2"}

	packer.observers.Attach(logger1)
	packer.observers.Attach(logger2)

	bin, _ := packer.createBin()
	fmt.Printf("[BinCreated] Bin ID: %d\n", bin.id)

	block1 := &Block[string]{
		Size: func() int { return 30 },
		Data: "Block 1",
	}

	block2 := &Block[string]{
		Size: func() int { return 50 },
		Data: "Block 2",
	}

	block3 := &Block[string]{
		Size: func() int { return 40 },
		Data: "Block 3",
	}

	result1 := bin.Add(block1)
	fmt.Printf("Add block1 (size 30): %t\n", result1)

	result2 := bin.Add(block2)
	fmt.Printf("Add block2 (size 50): %t\n", result2)

	result3 := bin.Add(block3)
	fmt.Printf("Add block3 (size 40): %t\n", result3)

	// Output:
	// [Logger1] Bin created: ID=0
	// [Logger2] Bin created: ID=0
	// [BinCreated] Bin ID: 0
	// Add block1 (size 30): true
	// Add block2 (size 50): true
	// Add block3 (size 40): false
}
