package spanmap

const minFill = 30 // shrink at 30%

type mapItem struct {
	index uint64
	item  interface{}
}

// Map is a map that is optimized for data that span contiguous indexes.
type Map struct {
	min, max uint64
	items    []mapItem
	mask     uint64
	shrinkAt int
	len      int
}

func (m *Map) grow() {
	items := make([]mapItem, len(m.items)*2)
	mask := uint64(len(items) - 1)
	for idx := m.min; idx <= m.max; idx++ {
		if m.items[idx&m.mask].index == idx &&
			m.items[idx&m.mask].item != nil {
			items[idx&mask] = m.items[idx&m.mask]
		}
	}
	m.items = items
	m.mask = uint64(len(m.items) - 1)
	m.shrinkAt = len(m.items) * minFill / 100
}

func (m *Map) shrink() {
	sz := 1
	for sz < m.shrinkAt {
		sz *= 2
	}
	items := make([]mapItem, sz)
	mask := uint64(len(items) - 1)
	for idx := m.min; idx <= m.max; idx++ {
		if m.items[idx&m.mask].index == idx &&
			m.items[idx&m.mask].item != nil {
			items[idx&mask] = m.items[idx&m.mask]
		}
	}
	m.items = items
	m.mask = uint64(len(m.items) - 1)
	m.shrinkAt = len(m.items) * minFill / 100
}

// Set an item at index. Item cannot be nil.
func (m *Map) Set(index uint64, item interface{}) interface{} {
	if item == nil {
		panic("nil item")
	}
	if len(m.items) == 0 {
		// create the initial slice
		m.items = make([]mapItem, 1)
		m.mask = uint64(len(m.items) - 1)
		m.shrinkAt = len(m.items) * minFill / 100
		m.min, m.max = index, index
		m.items[0].index = index
		m.items[0].item = item
		m.len = 1
		return nil
	}
	if m.items[index&m.mask].item != nil {
		if m.items[index&m.mask].index == index {
			prev := m.items[index&m.mask].item
			m.items[index&m.mask].item = item
			return prev
		}
		for m.items[index&m.mask].item != nil {
			m.grow()
		}
	}
	if index < m.min {
		m.min = index
	} else if index > m.max {
		m.max = index
	}
	m.items[index&m.mask].index = index
	m.items[index&m.mask].item = item
	m.len++
	return nil
}

// Get an item at index.
func (m *Map) Get(index uint64) interface{} {
	if len(m.items) == 0 {
		return nil
	}
	if m.items[index&m.mask].index == index {
		return m.items[index&m.mask].item
	}
	return nil
}

// Delete an index
func (m *Map) Delete(index uint64) interface{} {
	if len(m.items) == 0 || m.items[index&m.mask].index != index ||
		m.items[index&m.mask].item == nil {
		return nil
	}
	item := m.items[index&m.mask].item
	m.items[index&m.mask].index = 0
	m.items[index&m.mask].item = nil
	m.len--
	if m.min == index {
		if m.max == index {
			m.min, m.max = 0, 0
			m.items = nil
			return item
		}
		for {
			m.min++
			if m.items[(m.min)&m.mask].item == nil {
				continue
			}
			break
		}
	} else if m.max == index {
		for m.max > m.min {
			m.max--
			if m.items[(m.max)&m.mask].item == nil {
				continue
			}
			break
		}
	}
	if (m.max - m.min + 1) <= uint64(m.shrinkAt) {
		m.shrink()
	}
	return item
}

// Len returns the number of items in map
func (m *Map) Len() int {
	return m.len
}

// Min returns the minimum index
func (m *Map) Min() uint64 {
	return m.min
}

// Max returns the maximum index
func (m *Map) Max() uint64 {
	return m.max
}

// func printLog(l *Map) {
// 	fmt.Printf("%02d-%02d [ ", l.first, l.last)
// 	for i := 0; i < len(l.items); i++ {
// 		if l.items[i].item == nil {
// 			fmt.Printf("-- ")
// 		} else {
// 			fmt.Printf("%02d ", l.items[i].index)
// 		}
// 	}
// 	println("]")
// }
