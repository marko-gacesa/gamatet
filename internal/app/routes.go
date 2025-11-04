// Copyright (c) 2024, 2025 by Marko Gaćeša

package app

type route string

type routes struct {
	id   route
	prev *routes
}

func (r *routes) clear() {
	r.id = ""
	r.prev = nil
}

func (r *routes) push(id route) *routes {
	if id == "" {
		panic("can't push empty")
	}

	r.prev = &routes{
		id:   r.id,
		prev: r.prev,
	}
	r.id = id

	return r
}

func (r *routes) pop() route {
	if r.prev == nil {
		return ""
	}

	id := r.id
	*r = *(r.prev)

	return id
}

func (r *routes) curr() route {
	return r.id
}

func (r *routes) path() []route {
	var list []route
	for curr := r; curr != nil; curr = curr.prev {
		list = append(list, curr.id)
	}
	return list
}
