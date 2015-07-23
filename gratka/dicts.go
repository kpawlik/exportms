package gratka

import (
	//	"text/template"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type dict map[int]string

func (d dict) GetVal(val string) (result string, ok bool) {
	val = strings.ToLower(val)
	for k, v := range d {
		if ok = (v == val); ok {
			result = strconv.Itoa(k)
			return
		}
	}
	return
}

func regions() dict { //region
	return dict{
		0:   "cały kraj",
		1:   "dolnośląskie",
		2:   "kujawsko-pomorskie",
		3:   "opolskie",
		4:   "lubelskie",
		5:   "lubuskie",
		6:   "łódzkie",
		7:   "małopolskie",
		8:   "mazowieckie",
		9:   "podkarpackie",
		10:  "podlaskie",
		11:  "pomorskie",
		12:  "śląskie",
		13:  "świętokrzyskie",
		14:  "warmińsko-mazurskie",
		15:  "wielkopolskie",
		16:  "zachodnio-pomorskie",
		100: "zagranica"}
}

func buildingType() dict { //typ_budynku
	return dict{
		1:  "wolnostojący",
		2:  "segment środkowy",
		3:  "segment skrajny",
		4:  "bliźniak",
		5:  "pół bliźniaka",
		6:  "kamienica",
		7:  "willa",
		8:  "rezydencja",
		9:  "dworek",
		10: "szeregowy",
		11: "piętro domu",
		12: "rekreacyjny"}
}

func typeOfParcele() dict { //rodzaj_dzialki
	return dict{
		1:  "budowlana",
		2:  "rolna",
		3:  "inwestycyjna",
		4:  "usługowa",
		5:  "przemysłowa",
		6:  "rekreacyjna",
		7:  "leśna",
		8:  "siedliskowa",
		9:  "gospodarstwo",
		10: "rzemieślnicza",
		11: "rolno-budowlana"}
}

func typeOfLocal() dict { //typ_lokalu
	return dict{
		1:  "biurowiec",
		2:  "pawilon",
		3:  "sklep",
		4:  "centrum handlowe",
		5:  "kiosk",
		6:  "blok",
		7:  "kamienica",
		8:  "magazyn",
		9:  "hala produkcyjna",
		10: "hotel",
		11: "kawiarnia",
		12: "motel",
		13: "pensjonat",
		14: "pub",
		15: "restauracja",
		16: "salon",
		17: "warsztat",
		18: "zajazd",
		19: "stacja benzynowa",
		20: "biuro",
		21: "gospodarstwo rolne",
		22: "dwór",
		23: "pałac",
		24: "zamek",
		25: "młyn",
		26: "spichlerz",
		27: "folwark"}
}

func characterOfLocal() dict { //charakter_lokalu
	return dict{
		1: "usługowy",
		2: "handlowy",
		3: "produkcyjny",
		4: "magazynowy",
		5: "mieszkalny"}
}

func kindOfParcele() dict { //rodzaj_dzialki
	return dict{
		1:  "budowlana",
		2:  "rolna",
		3:  "inwestycyjna",
		4:  "usługowa",
		5:  "przemysłowa",
		6:  "rekreacyjna",
		7:  "leśna",
		8:  "siedliskowa",
		9:  "gospodarstwo",
		10: "rzemieślnicza",
		11: "rolno-budowlana"}
}

func caracterOfHouse() dict { //domy_charakter
	return dict{
		1: "jednorodzinny",
		2: "dwurodzinny",
		3: "wielorodzinny"}
}

func access() dict { //dojazd
	return dict{
		1: "asfaltowy",
		2: "utwardzony",
		3: "polny"}
}

func building() dict { // typ_zabudowy
	return dict{
		1: "blok",
		2: "kamienica",
		3: "dom",
		4: "apartamentowiec",
		5: "wieżowiec"}
}

func currency() dict {
	return dict{
		1: "pln",
		2: "usd",
		4: "eur"}
}

type Dicts struct {
	cache map[string]dict
}

func NewDicts() *Dicts {
	return &Dicts{make(map[string]dict)}
}

func (d *Dicts) Get(name string) dict {
	var (
		di dict
		ok bool
	)
	//print(name, "   ", "typeOfBuilding" == name, "\n")
	if di, ok = d.cache[name]; ok {
		return di
	}
	if di = getDict(name); di == nil {
		log.Panic(fmt.Sprintf("Unrecognized dictionary '%s'", name))
	}
	d.cache[name] = di
	return di

}

func getDict(name string) dict {
	name = strings.ToLower(name)
	switch name {
	case "region":
		return regions()
	case "currency":
		return currency()
	case "building":
		return building()
	case "buildingtype":
		return buildingType()
	case "local":
		return typeOfLocal()
	case "parcele":
		return typeOfParcele()
	}

	return nil
}
