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

type Observer[T any] interface {
	Update(event BinEvent[T])
}

type Subject[T any] interface {
	Attach(observer Observer[T])
	Detach(observer Observer[T])
	Notify(event BinEvent[T])
}

type ConcreteSubject[T any] struct {
	observers []Observer[T]
}

func (s *ConcreteSubject[T]) Register(o Observer[T]) {
	s.observers = append(s.observers, o)
}

func (s *ConcreteSubject[T]) Deregister(o Observer[T]) {
	for i, observer := range s.observers {
		if o == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

func (s *ConcreteSubject[T]) Attach(observer Observer[T]) {
	s.observers = append(s.observers, observer)
}

func (s *ConcreteSubject[T]) Detach(observer Observer[T]) {
	for i, obs := range s.observers {
		if obs == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

func (s *ConcreteSubject[T]) Notify(event BinEvent[T]) {
	for _, observer := range s.observers {
		observer.Update(event)
	}
}
