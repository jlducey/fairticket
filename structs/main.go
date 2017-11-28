package main

import "fmt"

type contactInfo struct { //struct can have mixed data types like string and int
	email string
	zip   int
}
type person struct {
	firstName string
	lastName  string
	contact   contactInfo // struct can be contained in another struct
}

func main() {
	//alex := person{"Alex", "Anderson"}                      // not recommended as it relies on position only for fields
	//kelly := person{firstName: "Kelly", lastName: "Walker"} // better to update later as order isn't relied on
	//fmt.Println(alex)
	//fmt.Println(kelly)
	var john person // third way to declare a struct
	fmt.Println(john)
	fmt.Printf("%+v", john) // note the + sign changes output over just a %v
	john.firstName = "John"
	john.lastName = "Cougar"
	fmt.Println(john)
	fmt.Printf("%+v", john)

	jim := person{
		firstName: "Jim",   // for an item of a given type struct.. must have comma at end of each line except last
		lastName:  "Ducey", // for embedded struct in a struct.. comma before
		contact: contactInfo{
			email: "jlducey@hotmail.com", // for each item in struct a comma at end
			zip:   67147,                 //last item even needs comma
		}, // this needs comma!
	}
	jimpointer := &jim // & means give me the memory address of the variable jim... its a reference to the struct in memory
	jimpointer.updateName("Jimmy")
	jim.print()
}

func (p person) print() {
	fmt.Printf("%+v", p)
}

func (pointerToPerson *person) updateName(newFirstName string) { // * says give me value this memory address is pointing at
	(pointerToPerson).firstName = newFirstName
}
