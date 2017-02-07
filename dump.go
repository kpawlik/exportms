package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/kpawlik/exportms/db"
	"github.com/kpawlik/exportms/utils"
	"github.com/kpawlik/exportms/xml"
)

// fillOffer fill offer record from database
func fillOffer(offer *db.Offer, additionals *db.Additional, person, contact *db.Person) *xml.Offer {

	creatorName := fmt.Sprintf("%s %s", person.StringAt(db.Imie), person.StringAt(db.Nazwisko))
	contactName := fmt.Sprintf("%s %s", contact.StringAt(db.Imie), contact.StringAt(db.Nazwisko))
	xmlOffer := xml.NewOffer()
	xmlOffer.Add(xml.NewElem("id_oferty", offer.IntAt(db.IDOferty)))
	xmlOffer.Add(xml.NewCdataElem("nazwa_oferty", offer.StringAt(db.NazwaOferty)))
	xmlOffer.Add(xml.NewCdataElem("grupa", offer.StringAt(db.IDGrupy)))
	xmlOffer.Add(xml.NewCdataElem("podgrupa", offer.StringAt(db.RodzajOferty)))
	xmlOffer.Add(xml.NewCdataElem("numer_oferty", offer.StringAt(db.NumerOferty)))
	xmlOffer.Add(xml.NewCdataElem("opis", offer.StringAt(db.TrescOferty)))
	xmlOffer.Add(xml.NewCdataElem("uwagi", offer.StringAt(db.Uwagi)))
	xmlOffer.Add(xml.NewCdataElem("oferta_zawiera", offer.StringAt(db.OfertaZamowienia)))

	xmlOffer.Add(xml.NewCdataElem("wprowadzajacy", creatorName))
	xmlOffer.Add(xml.NewCdataElem("kontakt_osoba", contactName))
	xmlOffer.Add(xml.NewCdataElem("kontakt_telefon_1", contact.DecodeStringAt(db.Telefon)))
	xmlOffer.Add(xml.NewCdataElem("kontakt_telefon_1", contact.DecodeStringAt(db.Telefon2)))
	xmlOffer.Add(xml.NewCdataElem("kontakt_telefon_1", contact.DecodeStringAt(db.Telefon3)))
	xmlOffer.Add(xml.NewCdataElem("kontakt_email", contact.DecodeStringAt(db.Email)))
	xmlOffer.Add(xml.NewElem("cena_groszy", offer.IntAt(db.Cena)))
	xmlOffer.Add(xml.NewCdataElem("waluta", "PLN"))
	xmlOffer.Add(xml.NewElem("data_ost_modyfikacji", offer.StringAt(db.DataModyfikacji)))
	xmlOffer.Add(xml.NewElem("data_wprowadzenia", offer.StringAt(db.DataZgloszenia)))

	xmlOffer.Add(xml.NewCdataElem("miejscowosc", additionals.StringAt(db.Pole01)))
	xmlOffer.Add(xml.NewCdataElem("region", additionals.StringAt(db.Pole02)))
	xmlOffer.Add(xml.NewCdataElem("wojewodztwo", additionals.StringAt(db.Pole03)))
	xmlOffer.Add(xml.NewCdataElem("dzielnica", additionals.StringAt(db.Pole04)))
	xmlOffer.Add(xml.NewCdataElem("ulica", additionals.StringAt(db.Pole05)))
	xmlOffer.Add(xml.NewCdataElem("kod_pocztowy", additionals.StringAt(db.Pole06)))
	xmlOffer.Add(xml.NewCdataElem("powiat", additionals.StringAt(db.Pole10)))
	xmlOffer.Add(xml.NewCdataElem("gmina", additionals.StringAt(db.Pole11)))

	xmlOffer.Add(xml.NewCdataElem("komunikacja", additionals.StringAt(db.Pole07)))
	xmlOffer.Add(xml.NewCdataElem("dojazd", additionals.StringAt(db.Pole08)))
	xmlOffer.Add(xml.NewCdataElem("polozenie", additionals.StringAt(db.Pole09)))
	xmlOffer.Add(xml.NewElemNull("powierzchnia_calkowita", additionals.StringAt(db.Pole60), "0"))
	xmlOffer.Add(xml.NewElemNull("powierzchnia_uzytkowa", additionals.StringAt(db.Pole61), "0"))
	xmlOffer.Add(xml.NewElemNull("powierzchnia_dzialki", additionals.StringAt(db.Pole62), "0"))
	xmlOffer.Add(xml.NewElemNull("szerokosc_dzialki", additionals.StringAt(db.Pole63), "0"))
	xmlOffer.Add(xml.NewElemNull("dlugosc_dzialki", additionals.StringAt(db.Pole64), "0"))
	xmlOffer.Add(xml.NewElemNull("kubatura", additionals.StringAt(db.Pole65), "0"))
	xmlOffer.Add(xml.NewElemNull("pochylenie_dzialki", additionals.StringAt(db.Pole21), "0"))
	xmlOffer.Add(xml.NewElemNull("ilosc_pokoi", additionals.StringAt(db.Pole66), "0"))
	xmlOffer.Add(xml.NewElemNull("ilosc_pomieszczen", additionals.StringAt(db.Pole67), "0"))
	xmlOffer.Add(xml.NewElemNull("ilosc_pieter", additionals.StringAt(db.Pole68), "0"))
	xmlOffer.Add(xml.NewElemNull("pietro", additionals.StringAt(db.Pole69), "0"))
	xmlOffer.Add(xml.NewElem("rodzaj_pomieszczen", additionals.StringAt(db.Pole22)))
	xmlOffer.Add(xml.NewElem("rozklad_pomieszczen", additionals.StringAt(db.Pole23)))
	xmlOffer.Add(xml.NewElem("rodzaj_zabudowy", additionals.StringAt(db.Pole24)))
	xmlOffer.Add(xml.NewElemNull("rok_budowy", additionals.StringAt(db.Pole25), "0"))
	xmlOffer.Add(xml.NewElemNull("typ_konstrukcji", additionals.StringAt(db.Pole26), "0"))
	xmlOffer.Add(xml.NewElem("stan_techniczny", additionals.StringAt(db.Pole27)))
	xmlOffer.Add(xml.NewElem("cena_dodatkowa_sprzedazy", additionals.StringAt(db.Pole41)))
	xmlOffer.Add(xml.NewElem("cena_dodatkowa_najmu", additionals.StringAt(db.Pole42)))
	xmlOffer.Add(xml.NewElem("cena_dodatkowa_czynsz", additionals.StringAt(db.Pole43)))
	xmlOffer.Add(xml.NewElem("forma_wlasnosci", additionals.StringAt(db.Pole44)))
	xmlOffer.Add(xml.NewElem("rodzaj_dzialki", additionals.StringAt(db.Pole28)))
	return xmlOffer
}

// getImages download images in concurrent way
func getImages(c chan error, imgs *db.Images, destFolder string) {
	c <- imgs.SaveImages(filepath.Join(destFolder, "images"))
}

// worker is concurrent function to handle len(in) concurrent offers
func worker(in chan *db.Offer, out chan *xml.Offer, specials db.SpecialsMap) {
	c := make(chan *xml.Offer)
	for offer := range in {
		dbConn := offer.DB
		offerID := offer.ID()
		go getOfferData(dbConn, offerID, specials, c)
		out <- <-c
	}
}

// getOfferData gather all data for offer and return filled object to channel
func getOfferData(dbConn *db.DB, offerID int, specialsMap db.SpecialsMap, offersChan chan *xml.Offer) {
	var (
		err error
	)
	offer := db.NewOffer(dbConn, offerID)
	utils.LogErrf(offer.Get(), "Get offer (%d)", offerID)
	// contact and owner data
	personID := offer.StringAt(db.IDWprowadzajacego)
	contactID := offer.StringAt(db.IDWlasciciela)

	person := db.NewPerson(dbConn)
	utils.LogErrf(person.Get(personID), "Get person (%s) for offer (%d)", personID, offerID)
	contact := db.NewPerson(dbConn)
	utils.LogErrf(contact.Get(contactID), "Get contact (%s) for offer (%d)", personID, offerID)
	// additionals data
	add := db.NewAdditional(dbConn)
	utils.LogErrf(add.Get(offerID), "Get additional data for offer (%d)", offerID)
	//
	xmlOffer := fillOffer(offer, add, person, contact)
	sid := offer.StrID()
	var imgIds []string
	// images and specials
	images := db.NewImages(dbConn, sid)
	if imgIds, err = images.FileNames(); err != nil {
		utils.LogErrf(err, "Images ids for offer (%d)", offerID)
	}
	go getImages(imagesChan, images, workDir)
	xmlOffer.Pictures = xml.NewListElem("zdjecia", "zdjecie")
	xmlOffer.Pictures.AddMany(imgIds...)
	xmlOffer.Specials = xml.NewListElem("wyroznienia", "wyroznienie")
	xmlOffer.Specials.AddMany(specialsMap[sid]...)
	offersChan <- xmlOffer
}

func dumpAsXML(conf *config) (count int, err error) {
	workDir := conf.workDir
	exportID := fmt.Sprintf("%s_00.xml", time.Now().Format(DateFormat))
	exportXMLName := filepath.Join(workDir, exportID)

	// sync channels
	imagesChan = make(chan error)
	inOffersChan := make(chan *db.Offer, 5)
	outOffersChan := make(chan *xml.Offer)
	// connect to db
	dbConn := &db.DB{}
	if err = dbConn.Connect(DbUser, DbPass, DbHost, DbName); err != nil {
		err = utils.Errorf("Database connection %v", err)
		return
	}
	defer dbConn.Close()
	// specials are common for all offers
	specials := db.NewSpecials(dbConn)
	if err = specials.GetAll(); err != nil {
		return
	}
	specialsMap := specials.AsMap()
	// start workers
	for i := 0; i < noOfWorkers; i++ {
		go worker(inOffersChan, outOffersChan, specialsMap)
	}
	// start getting offers
	offers := db.NewOffers(dbConn)
	if err = offers.GetAllIds(); err != nil {
		return
	}
	go func() {
		for _, offerID := range offers.Ids() {
			offer := db.NewOffer(dbConn, offerID)
			inOffersChan <- offer
		}
	}()
	step := 25
	total := len(offers.Ids())
	counter := 0
	// collect offers from workers
	xmlOffers := xml.NewOffers()
	for range offers.Ids() {
		xmlOffer := <-outOffersChan
		xmlOffers.Add(xmlOffer)
		counter++
		if (counter % step) == 0 {
			log.Printf("Completed %d/%d\n", counter, total)
		}
	}
	log.Printf("Completed %d/%d\n", counter, total)
	log.Printf("Wait for images\n")
	// wait for all images finish download
	for range offers.Ids() {
		<-imagesChan
	}
	system := xml.NewSystem("1", "date", "hour", "www")
	head := xml.NewHeader(system, xmlOffers)
	if err = xml.Write(head, exportXMLName); err != nil {
		return
	}
	if err = ioutil.WriteFile(filepath.Join(workDir, "last.id"), []byte(exportID), os.ModePerm); err != nil {
		return
	}
	count = len(offers.Ids())
	return
}
