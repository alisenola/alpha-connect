package utils

// A ring set is a fixed size set that store the N last elements
type RingSet struct {
	buffer []element
	head   int
	tail   int
	len    int
}

type element struct {
	key   interface{}
	value interface{}
}

func NewRingSet(n int) *RingSet {
	rs := RingSet{}
	rs.head = -1
	rs.tail = 0
	rs.len = n
	rs.buffer = make([]element, n, n)
	return &rs
}

// Check if the buffer contains an element.
// the data structure is optimised for a lot
// of hits near the head
func (rs *RingSet) Contains(key interface{}) bool {
	if rs.head < 0 {
		return false
	}
	i := rs.head
	for i != rs.tail {
		if rs.buffer[i].key == key {
			return true
		}
		// Golang neg modulo n is neg
		if i == 0 {
			i = rs.len - 1
		} else {
			i = i - 1
		}
	}
	// Check tail
	if rs.buffer[i].key == key {
		return true
	}

	return false
}

func (rs *RingSet) Get(key interface{}) interface{} {
	if rs.head < 0 {
		return nil
	}
	i := rs.head
	for i != rs.tail {
		if rs.buffer[i].key == key {
			return rs.buffer[i].value
		}
		// Golang neg modulo n is neg
		if i == 0 {
			i = rs.len - 1
		} else {
			i = i - 1
		}
	}
	// Check tail
	if rs.buffer[i].key == key {
		return rs.buffer[i].value
	}

	return nil
}

func (rs *RingSet) Insert(key interface{}, value interface{}) bool {
	if rs.Contains(key) {
		return false
	}

	rs.head = (rs.head + 1) % rs.len
	if rs.head == rs.tail {
		rs.tail = (rs.tail + 1) % rs.len
	}
	rs.buffer[rs.head] = element{key: key, value: value}

	return true
}

func (rs *RingSet) Upsert(key interface{}, value interface{}) bool {
	if rs.Contains(key) {
		i := rs.head
		for i != rs.tail {
			if rs.buffer[i] == key {
				rs.buffer[i] = element{key: key, value: value}
			}
			// Golang neg modulo n is neg
			if i == 0 {
				i = rs.len - 1
			} else {
				i = i - 1
			}
		}
		// Check tail
		if rs.buffer[i] == key {
			rs.buffer[i] = element{key: key, value: value}
		}
	}

	rs.head = (rs.head + 1) % rs.len
	if rs.head == rs.tail {
		rs.tail = (rs.tail + 1) % rs.len
	}
	rs.buffer[rs.head] = element{key: key, value: value}

	return true
}
