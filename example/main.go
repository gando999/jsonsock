package main

import (
	"fmt"
	"github.com/gando999/jsonsock"
)

const Target = "localhost:63011"

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
	dom.Number = 100
	dom.Cost = 200
	return *dom
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

func sendRequest(method string, args []interface{}) {
	client := jsonsock.CreateClient(Target)
	fmt.Println(client.Send(method, args))
}

func main() {
	go startServer()

	sendRequest("testing.SayHello", []interface{}{"Donald"})
	sendRequest("testing.SayHelloPointer", []interface{}{"Donald Pointer"})
	sendRequest("testing.SayHelloFromPointer", []interface{}{})
	d := Domain{256, 512, "Domain"}
	sendRequest("testing.UseDomain", []interface{}{d})
	sendRequest("testing.UseDomainPointer", []interface{}{d})

	m := make(map[string]interface{})
	m["Number"] = 55
	m["Cost"] = 65
	m["Message"] = "I was a map"
	sendRequest("testing.UseDomain", []interface{}{m})

	sl := []string{"Hello", "there", "sailor"}
	sendRequest("testing.UseSlice", []interface{}{sl})

	floats := []float64{10, 9, 8, 7, 6, 5}
	sendRequest("testing.UseMix", []interface{}{floats, "Man alive!", m})

	sendRequest("testing.StructSlice", []interface{}{[]Domain{d}})

	d.Message = "Inner nested"
	x := Nested{d, "Outer nested"}
	sendRequest("testing.UseNested", []interface{}{x})
}

func startServer() {
	server := jsonsock.CreateServer(Target)
	server.Register("testing", &Testing{})
	server.Start()
}
