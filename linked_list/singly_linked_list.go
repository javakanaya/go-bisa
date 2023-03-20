package main

import "fmt"

type sListNode struct {
	data int
	next *sListNode
}

type sList struct {
	head *sListNode
}

func main() {
	fmt.Println("This is a singly linked list")

	myList := sList{}
	for i := 0; i < 5; i++ {
		sList_pushBack(&myList, i)
		sList_printList(&myList)
	}

	for i := 5; i < 10; i++ {
		sList_pushFront(&myList, i)
		sList_printList(&myList)
	}

	sList_popBack(&myList)
	sList_printList(&myList)
	sList_popBack(&myList)
	sList_printList(&myList)
	sList_popBack(&myList)
	sList_printList(&myList)

	sList_popFront(&myList)
	sList_printList(&myList)

}

func sList_pushFront(list *sList, value int) {
	newNode := &sListNode{
		data: value,
		next: nil,
	}

	if list.head == nil {
		list.head = newNode
	} else {
		newNode.next = list.head
		list.head = newNode
	}
}

func sList_pushBack(list *sList, value int) {
	newNode := &sListNode{
		data: value,
		next: nil,
	}

	if list.head == nil {
		list.head = newNode
	} else {
		temp := list.head
		for temp.next != nil {
			temp = temp.next
		}
		temp.next = newNode
	}
}

func sList_popFront(list *sList) {
	fmt.Println("popFront")
	if list.head != nil {
		list.head = list.head.next
	}
}

func sList_popBack(list *sList) {
	fmt.Println("popBack")
	if list.head != nil {
		if list.head.next == nil {
			list.head = nil
			return
		}

		temp := list.head
		for temp.next.next != nil {
			temp = temp.next
		}
		temp.next = nil
	}
}

func sList_printList(list *sList) {
	if list.head != nil {
		fmt.Printf("%d", list.head.data)

		for temp := list.head.next; temp != nil; temp = temp.next {
			fmt.Printf(" -> %d", temp.data)
		}

		fmt.Printf("\n")
	} else {
		fmt.Println("List is Empty")
	}
}
