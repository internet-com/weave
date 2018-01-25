package store

import "fmt"

////////////////////////////////////////////////
// Slice -> Iterator

// Model groups together key and value to help build Iterators
type Model struct {
	Key   []byte
	Value []byte
}

// pair is test helper.. make public??
func pair(key, value []byte) Model {
	return Model{
		Key:   key,
		Value: value,
	}
}

// SliceIterator wraps an Iterator over a slice of models
//
// TODO: make this private and only expose Iterator interface????
type SliceIterator struct {
	data []Model
	idx  int
}

var _ Iterator = (*SliceIterator)(nil)

// NewSliceIterator creates a new Iterator over this slice
func NewSliceIterator(data []Model) *SliceIterator {
	return &SliceIterator{
		data: data,
	}
}

// Valid implements Iterator and returns true iff it can be read
func (s *SliceIterator) Valid() bool {
	return s.idx < len(s.data)
}

// Next moves the iterator to the next sequential key in the database, as
// defined by order of iteration.
//
// If Valid returns false, this method will panic.
func (s *SliceIterator) Next() {
	s.assertValid()
	s.idx++
}

func (s *SliceIterator) assertValid() {
	if s.idx >= len(s.data) {
		panic("Passed end of slice")
	}
}

// Key returns the key of the cursor.
func (s *SliceIterator) Key() (key []byte) {
	s.assertValid()
	return s.data[s.idx].Key
}

// Value returns the value of the cursor.
func (s *SliceIterator) Value() (value []byte) {
	s.assertValid()
	return s.data[s.idx].Value
}

// Close releases the Iterator.
func (s *SliceIterator) Close() {
	s.data = nil
}

/////////////////////////////////////////////////////
// Empty KVStore

// EmptyKVStore never holds any data, used as a base layer to test caching
type EmptyKVStore struct{}

var _ KVStore = EmptyKVStore{}

// Get always returns nil
func (e EmptyKVStore) Get(key []byte) []byte { return nil }

// Has always returns false
func (e EmptyKVStore) Has(key []byte) bool { return false }

// Set is a noop
func (e EmptyKVStore) Set(key, value []byte) {}

// Delete is a noop
func (e EmptyKVStore) Delete(key []byte) {}

// Iterator is always empty
func (e EmptyKVStore) Iterator(start, end []byte) Iterator {
	return NewSliceIterator(nil)
}

// ReverseIterator is always empty
func (e EmptyKVStore) ReverseIterator(start, end []byte) Iterator {
	return NewSliceIterator(nil)
}

// NewBatch returns a batch that can write to this tree later
func (e EmptyKVStore) NewBatch() Batch {
	return NewNonAtomicBatch(e)
}

////////////////////////////////////////////////////
// Non-atomic batch (dummy implementation)

type opKind int32

const (
	setKind opKind = iota + 1
	delKind
)

// op is either set or delete
type op struct {
	kind  opKind
	key   []byte
	value []byte // only for set
}

func (o op) apply(out SetDeleter) {
	switch o.kind {
	case setKind:
		out.Set(o.key, o.value)
	case delKind:
		out.Delete(o.key)
	default:
		panic(fmt.Sprintf("Unknown kind: %d", o.kind))
	}
}

// setOp is a helper to create a set operation
func setOp(key, value []byte) op {
	return op{
		kind:  setKind,
		key:   key,
		value: value,
	}
}

// delOp is a helper to create a del operation
func delOp(key []byte) op {
	return op{
		kind: delKind,
		key:  key,
	}
}

//---------------- helpers -------------

func makeSetOps(ms []Model) []op {
	res := make([]op, len(ms))
	for i, m := range ms {
		res[i] = setOp(m.Key, m.Value)
	}
	return res
}

func makeDelOps(ms []Model) []op {
	res := make([]op, len(ms))
	for i, m := range ms {
		res[i] = delOp(m.Key)
	}
	return res
}

// NonAtomicBatch just piles up ops and executes them later
// on the underlying store. Can be used when there is no better
// option (for in-memory stores).
//
// NOTE: Never use this for KVStores that are persistent
type NonAtomicBatch struct {
	out SetDeleter
	ops []op
}

var _ Batch = (*NonAtomicBatch)(nil)

// NewNonAtomicBatch creates an empty batch to be later writen
// to the KVStore
func NewNonAtomicBatch(out SetDeleter) *NonAtomicBatch {
	return &NonAtomicBatch{
		out: out,
	}
}

// Set adds a set operation to the batch
func (b *NonAtomicBatch) Set(key, value []byte) {
	set := op{
		kind:  setKind,
		key:   key,
		value: value,
	}
	b.ops = append(b.ops, set)
}

// Delete adds a delete operation to the batch
func (b *NonAtomicBatch) Delete(key []byte) {
	del := op{
		kind: delKind,
		key:  key,
	}
	b.ops = append(b.ops, del)
}

// Write writes all the ops to the underlying store and resets
func (b *NonAtomicBatch) Write() {
	for _, op := range b.ops {
		op.apply(b.out)
	}
	b.ops = nil
}
