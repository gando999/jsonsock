package main

type PetRepository interface {
	Load(id float64) Animal
	Save(obj Animal) float64
}

type Animal interface {
	GetId() float64
	Talk() string
}

type AnimalRepository struct {
	m map[float64]Animal
}

type Dog struct {
	Id    float64
	Breed string
	Name  string
}

func (repo *AnimalRepository) Load(id float64) Animal {
	return repo.m[id]
}

func (repo *AnimalRepository) Save(dog Animal) float64 {
	repo.m[dog.GetId()] = dog
	return dog.GetId()
}

func (dog Dog) Talk() string {
	return "Woof!"
}

func (dog Dog) GetId() float64 {
	return dog.Id
}

type PetShopApi struct {
	Name string
	PetRepository
}

func (api *PetShopApi) GetDogById(id float64) Animal {
	return api.Load(id)
}

func (api *PetShopApi) SaveDog(dog Dog) float64 {
	return api.Save(dog)
}
