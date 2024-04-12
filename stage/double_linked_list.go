package stage

type LinkedList[T any] struct {
	Item   T
	Prev   *LinkedList[T]
	Next   *LinkedList[T]
	Root   *LinkedList[T]
	tail   *LinkedList[T]
	Length uint32
	dead   bool
}

func NewLinkedList[T any]() *LinkedList[T] {
	lst := LinkedList[T]{
		Prev:   nil,
		Next:   nil,
		Root:   nil,
		tail:   nil,
		Length: 0,
	}

	lst.Root = &lst
	lst.tail = &lst
	return &lst
}

func (h *LinkedList[T]) Drop() {
	if *&h.dead {
		return
	}

	if h == (*h).Root {
		panic("cant drop root of List")
	}

	if (*h.Root).tail == h {
		(*h.Root).tail = (*h).Prev
	}

	if (*h).Next != nil {
		(*(*h).Next).Prev = (*h).Prev
	}
	if h.Prev != nil {
		(*(*h).Prev).Next = (*h).Next
	}

	h.dead = true
	(*h.Root).Length--
}

// appending to the Root will always append to the end of the list
// appending to any none root entry will append after that entry
func (h *LinkedList[T]) Append(item T) *LinkedList[T] {
	// if the element knows the lists tail, it is the Root
	if (*h).tail != nil {
		h = (*h).tail
	}

	handler := LinkedList[T]{
		Item: item,
		Prev: h,
		Next: (*h).Next,
		Root: (*h).Root,
	}
	(*h.Root).Length = 1
	// Set a new tail if the current entry is the tail
	if (*h.Root).tail == h {
		(*h.Root).tail = &handler
	}
	(*h).Next = &handler

	return &handler
}

// pretends to the root, will always be attached after the root
// prepending to any none root entry adds the element before that entry
func (h *LinkedList[T]) Prepend(item T) *LinkedList[T] {
	// if the element knows the lists tail, it is the Root
	if (*h).tail != nil {
		// cant prepend to root,
		if (*h).Next != nil {
			// so either prepend to the next elment
			h = (*h).Next
		} else {
			// or treat this as an append
			return h.Append(item)
		}
	}

	handler := LinkedList[T]{
		Item: item,
		Prev: (*h).Prev,
		Next: h,
		Root: (*h).Root,
	}
	(*h.Root).Length++

	(*h).Prev = &handler

	return &handler
}
