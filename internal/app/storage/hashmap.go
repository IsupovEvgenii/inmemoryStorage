package storage

//default number of buckets
const BucketSize = 2 ^ 4

type HashMap struct {
	M          map[int]*List
	BucketSize int
}

// Create new Hashmap
func NewHashMap() *HashMap {
	return &HashMap{
		make(map[int]*List),
		BucketSize,
	}
}

// Put an item into the map
func (h *HashMap) Put(Key []byte, Value Item) {
	bucket := len(Key) % h.BucketSize
	MapEntry := &MapEntry{Key, nil, Value}
	if _, ok := h.M[bucket]; !ok {
		h.M[bucket] = &List{}
	}
	h.M[bucket].Insert(MapEntry)
}

// Retrieve an item from the map using Key
func (h *HashMap) Get(Key []byte) (*Item, bool) {
	bucket := len(Key) % h.BucketSize
	if List, ok := h.M[bucket]; ok {
		for l := List.Head; l != nil; l = l.Next {
			if Equal(l.Key, Key) {
				return &l.Value, true
			}
		}
	}
	return nil, false
}

// Delete from the map using Key
func (h *HashMap) Delete(Key []byte) {
	bucket := len(Key) % h.BucketSize
	if _, ok := h.M[bucket]; ok {
		h.M[bucket].Delete(Key)
	}
}

type MapEntry struct {
	Key   []byte
	Next  *MapEntry
	Value Item
}

// LinkedList
type List struct {
	Head *MapEntry
	Size int
}

func (l *List) Insert(MapEntry *MapEntry) {
	if l.Head == nil {
		l.Head = MapEntry
	} else {
		for n := &l.Head; n != nil; {
			if Equal((*n).Key, MapEntry.Key) {
				(*n).Value = MapEntry.Value
				return
			}
			if (*n).Next == nil {
				(*n).Next = MapEntry
				break
			}
			n = &(*n).Next
		}
	}
	l.Size++
}

func (l *List) Delete(Key []byte) {
	if Equal(l.Head.Key, Key) {
		l.Head = l.Head.Next
		l.Size--
		return
	}
	n := l.Head
	var prev = n
	for n != nil {
		if Equal(n.Key, Key) {
			prev.Next = n.Next
			break
		}
		prev = n
		n = n.Next
	}
	l.Size--
}

func Equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
