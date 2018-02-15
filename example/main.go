package main

import (
	"fmt"
	"github.com/gando999/jsonsock"
)

const Target = "localhost:63011"

func sendRequest(method string, args []interface{}) {
	client := jsonsock.CreateClient(Target)
	fmt.Println(client.Send(method, args))
}

func runTesting() {
	sendRequest("testing.SayHello", []interface{}{"Donald"})
	sendRequest("testing.SayHelloPointer", []interface{}{"Donald Pointer"})
	sendRequest("testing.SayHelloFromPointer", []interface{}{})
	d := Domain{256, 512, "Domain"}
	sendRequest("testing.UseDomain", []interface{}{d})
	sendRequest("testing.UseDomainPointer", []interface{}{d})
	sendRequest("testing.UseDomainReturnPointer", []interface{}{d})

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

func runPetShop() {
	dog := make(map[string]interface{})
	dog["Id"] = 1
	dog["Breed"] = "Collie"
	dog["Name"] = "Bonzo"
	sendRequest("petshop.SaveDog", []interface{}{dog})
	sendRequest("petshop.GetDogById", []interface{}{1})

	otherDog := Dog{2, "Alsatian", "Biffer"}
	sendRequest("petshop.SaveDog", []interface{}{otherDog})
	sendRequest("petshop.GetDogById", []interface{}{2})
}

func StartServer() {
	server := jsonsock.CreateServer(Target)
	repo := AnimalRepository{make(map[float64]Animal)}
	api := PetShopApi{"PetShopApi", &repo}
	server.Register("petshop", &api)
	server.Register("testing", &Testing{})
	server.Start()
}

func main() {
	go StartServer()
	runTesting()
	runPetShop()
}
