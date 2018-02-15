package main

import (
	"fmt"
)

type Testing struct {
	Name string
}

type Domain struct {
	Number  float64
	Cost    float64
	Message string
}

type Nested struct {
	Dom     Domain
	Message string
}

func (testing *Testing) SayHello(user string) string {
	return fmt.Sprintf("Hello from %s", user)
}

func (testing *Testing) SayHelloPointer(user *string) string {
	return fmt.Sprintf("Hello from pointer %s", *user)
}

func (testing *Testing) SayHelloFromPointer() *string {
	s := "Hello from pointer"
	return &s
}

func (testing *Testing) UseDomain(dom Domain) Domain {
	dom.Number = 100
	dom.Cost = 200
	return dom
}

func (testing *Testing) UseDomainPointer(dom *Domain) Domain {
	dom.Number = 10000
	dom.Cost = 20000
	return *dom
}

func (testing *Testing) UseDomainReturnPointer(dom *Domain) *Domain {
	dom.Number = 90000
	dom.Cost = 80000
	return dom
}

func (testing *Testing) UseSlice(someStrings []string) string {
	return fmt.Sprintf("Received %s", someStrings)
}

func (testing *Testing) UseMix(floats []float64, message string, dom Domain) Domain {
	dom.Number = floats[0]
	dom.Message = message
	dom.Cost = floats[0]
	return dom
}

func (testing *Testing) StructSlice(structs []Domain) Domain {
	d := structs[0]
	d.Number = 999
	d.Cost = 888
	d.Message = "I got updated"
	return d
}

func (testing *Testing) UseNested(nested Nested) Nested {
	return nested
}
