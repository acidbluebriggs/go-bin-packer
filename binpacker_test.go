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

import (
	"testing"
)

type testObserver struct {
	counter *struct {
		created int
		added   int
		ignored int
	}
}

func (o *testObserver) Update(event BinEvent[string]) {
	switch event.Type() {
	case BinCreated:
		o.counter.created++
	case BlockAdded:
		o.counter.added++
	case BlockIgnored:
		o.counter.ignored++
	}
}

func TestEventSystem(t *testing.T) {
	counter := struct {
		created int
		added   int
		ignored int
	}{}

	observer := &testObserver{counter: &counter}

	packer := &BinPacker[string]{
		max:   100,
		count: 0,
	}

	packer.observers.observers = make([]Observer[string], 0)

	packer.observers.Attach(observer)

	bin, err := packer.createBin()
	if err != nil {
		t.Fatalf("Failed to create bin: %v", err)
	}

	if counter.created != 1 {
		t.Errorf("Expected 1 BinCreated event, got %d", counter.created)
	}

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

	added1 := bin.Add(block1)
	added2 := bin.Add(block2)
	added3 := bin.Add(block3)

	if !added1 || !added2 || added3 {
		t.Errorf("Expected blocks 1 and 2 to be added, and block 3 to be rejected")
	}
}
