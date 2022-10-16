package gogo

type Comparable interface {
	Less(other interface{}) bool
}

type LinkedList[C Comparable] struct {
	Element C
	Next    *LinkedList[C]
}

func (ll LinkedList[C]) Add(element C) LinkedList[C] {
	// fix head
	if ll.Element.Less(element) {
		elementLinkedList := LinkedList[C]{
			Element: element,
			Next:    &ll,
		}
		return elementLinkedList
	}

	// fix current
	if ll.Next == nil || ll.Next.Element.Less(element) {
		elementLinkedList := LinkedList[C]{
			Element: element,
			Next:    ll.Next,
		}
		ll.Next = &elementLinkedList
		return ll
	}

	// fix recursively
	n := ll.Next.Add(element)
	ll.Next = &n
	return ll
}
