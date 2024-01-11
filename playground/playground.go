package main

/*
********************************************************************************
Golang - Asterisk and Ampersand Cheatsheet
********************************************************************************
Also available at: https://play.golang.org/p/lNpnS9j1ma
Allowed:
--------
p := Person{"Steve", 28} 	stores the value
p := &Person{"Steve", 28} 	stores the pointer address (reference)
PrintPerson(p) 			passes either the value or pointer address (reference)
PrintPerson(*p) 		passes the value
PrintPerson(&p) 		passes the pointer address (reference)
func PrintPerson(p Person)	ONLY receives the value
func PrintPerson(p *Person)	ONLY receives the pointer address (reference)
Not Allowed:
--------
p := *Person{"Steve", 28} 	illegal
func PrintPerson(p &Person)	illegal
*/

import "fmt"

type Person struct {
	Name string
	Age  int
}

func editPerson(person *Person) {
	person.Name = "Nils"
}

func wontEditPerson(person Person) {
	person.Name = "Oskar"
}

func main() {
	person := Person{Name: "Bob", Age: 89}
	personAsReference := &Person{Name: "Bob (as reference)", Age: 89}
	fmt.Println(person)

	person.Name = "Alice"
	fmt.Println(person)

	editPerson(&person)
	fmt.Println(person)

	fmt.Println(personAsReference)
	editPerson(personAsReference)
	fmt.Println(personAsReference)

	wontEditPerson(person)
	fmt.Println(person)
}
