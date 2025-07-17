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

import "errors"

type Block[T any] struct {
	Size func() int
	Data T
}

// EventType represents the type of bin packing event
type EventType int

const (
	BinCreated EventType = iota
	BlockAdded
	BlockIgnored
)

// BinEvent is the interface for all bin packing events
type BinEvent[T any] interface {
	Type() EventType
	Bin() *Bin[T]
}

// BinCreatedEvent is triggered when a new bin is created
type BinCreatedEvent[T any] struct {
	bin *Bin[T]
}

func (e BinCreatedEvent[T]) Type() EventType { return BinCreated }
func (e BinCreatedEvent[T]) Bin() *Bin[T]    { return e.bin }

// BlockAddedEvent is triggered when a block is added to a bin
type BlockAddedEvent[T any] struct {
	bin   *Bin[T]
	block Block[T]
}

func (e BlockAddedEvent[T]) Type() EventType { return BlockAdded }
func (e BlockAddedEvent[T]) Bin() *Bin[T]    { return e.bin }
func (e BlockAddedEvent[T]) Block() Block[T] { return e.block }

// BlockIgnoredEvent is triggered when a block cannot be added to a bin
type BlockIgnoredEvent[T any] struct {
	bin   *Bin[T]
	block Block[T]
}

func (e BlockIgnoredEvent[T]) Type() EventType { return BlockIgnored }
func (e BlockIgnoredEvent[T]) Bin() *Bin[T]    { return e.bin }
func (e BlockIgnoredEvent[T]) Block() Block[T] { return e.block }

type Sizer interface {
	Size() int
}

type Packer[T any] interface {
	Sizer
	Add(item *Block[T]) error
}

type Bin[T any] struct {
	limit int
	size  int
	id    int
	items []*Block[T]
}

func (bin *Bin[T]) Add(item *Block[T]) bool {
	if bin.hasRoom(item) {
		bin.items = append(bin.items, item)
		bin.size += item.Size()
		return true
	}
	return false
}

func (bin *Bin[T]) hasRoom(item *Block[T]) bool {
	return bin.size+item.Size() <= bin.limit
}

func (bin *Bin[T]) Size() int {
	return bin.size
}

func NewBin[T any](sizeLimit int, id int) (*Bin[T], error) {
	if sizeLimit < 1 {
		return nil, errors.New("size limit cannot be < 1")
	}

	return &Bin[T]{
		id:    id,
		limit: sizeLimit,
		items: make([]*Block[T], 0),
	}, nil
}

type BinPacker[T any] struct {
	max       int
	bins      []*Bin[T]
	count     int
	observers ConcreteSubject[T]
}

func (p *BinPacker[T]) Size() int {
	total := 0
	for _, bin := range p.bins {
		total += bin.size
	}
	return total
}

func (p *BinPacker[T]) createBin() (*Bin[T], error) {
	bin, err := NewBin[T](p.max, p.count)
	if err != nil {
		return nil, err
	}
	p.count++
	p.bins = append(p.bins, bin)

	event := BinCreatedEvent[T]{bin: bin}
	p.observers.Notify(event)

	return bin, nil
}

func (p *BinPacker[T]) Max() int {
	return p.max
}

func (p *BinPacker[T]) Bins() []*Bin[T] {
	return p.bins
}

func NewPacker[T any](maxSize int) (*BinPacker[T], error) {
	if maxSize < 1 {
		return nil, errors.New("maxSize must be > 0")
	}

	return &BinPacker[T]{}, nil
}

// First fit bin packer.. Maybe add a strategy for being able to veto an item
// to be added. Should it be the observers? More to think about this abstraction. or
// another strategy to intercept the attempt?

func (p *BinPacker[T]) Add(item *Block[T]) error {
	if item == nil {
		return errors.New("item cannot be nil")
	}

	var targetBin *Bin[T]
	added := false

	for _, bin := range p.bins {
		added = bin.Add(item)
		if added {
			targetBin = bin
			break
		}
	}

	if !added {
		var err error
		targetBin, err = NewBin[T](p.max, p.count)
		if err != nil {
			return err
		}
		p.bins = append(p.bins, targetBin)
		p.count++

		added = targetBin.Add(item)
		if !added {
			return errors.New("item too large for a new bin")
		}
		// TODO might not care about this event until we get to a vetoable
		// version, as created probably doesn't mean anything, but it can be ignore.
		// we will have to check if any observers reject this, or some other mechanism.
		// This may need a new tightly coupled observable.

		event := BinCreatedEvent[T]{bin: targetBin}
		p.observers.Notify(event)
	}

	event := BlockAddedEvent[T]{bin: targetBin, block: *item}
	p.observers.Notify(event)

	return nil
}
