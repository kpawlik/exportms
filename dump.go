package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"github.com/kpawlik/exportms/db"
	"github.com/kpawlik/exportms/utils"
	"github.com/kpawlik/exportms/xml"
	"os"
	"path/filepath"
	"time"
)

// fillOffer fill offer record from database
func fillOffer(offer *db.Offer, additionals *db.Additional, person, contact *db.Person) *xml.Offer {

	creatorName := fmt.Sprintf("%s %s", person.StringAt(db.IMIE), person.StringAt(db.NAZWISKO))
	contactName := fmt.Sprintf("%s %s", contact.StringAt(db.IMIE), contact.StringAt(db.NAZWISKO))
	xmlOffer := xml.NewOffer()
	xmlOffer.Add(xml.NewElem("id_oferty", offer.IntAt(db.ID_OFERTY)))
	xmlOffer.Add(xml.NewCdataElem("nazwa_oferty", offer.StringAt(db.NAZWA_OFERTY)))
	xmlOffer.Add(xml.NewCdataElem("grupa", offer.StringAt(db.ID_GRUPY)))
	xmlOffer.Add(xml.NewCdataElem("podgrupa", offer.StringAt(db.RODZAJ_OFERTY)))
	xmlOffer.Add(xml.NewCdataElem("numer_oferty", offer.StringAt(db.NUMER_OFERTY)))
	xmlOffer.Add(xml.NewCdataElem("opis", offer.StringAt(db.TRESC_OFERTY)))
	xmlOffer.Add(xml.NewCdataElem("uwagi", offer.StringAt(db.UWAGI)))
	xmlOffer.Add(xml.NewCdataElem("oferta_zawiera", offer.StringAt(db.OFERTA_ZAWIERA)))

	xmlOffer.Add(xml.NewCdataElem("wprowadzajacy", creatorName))
	xmlOffer.Add(xml.NewCdataElem("kontakt_osoba", contactName))
	xmlOffer.Add(xml.NewCdataElem("kontakt_telefon_1", contact.DecodeStringAt(db.TELEFON)))
	xmlOffer.Add(xml.NewCdataElem("kontakt_telefon_1", contact.DecodeStringAt(db.TELEFON_2)))
	xmlOffer.Add(xml.NewCdataElem("kontakt_telefon_1", contact.DecodeStringAt(db.TELEFON_3)))
	xmlOffer.Add(xml.NewCdataElem("kontakt_email", contact.DecodeStringAt(db.EMAIL)))
	xmlOffer.Add(xml.NewElem("cena_groszy", offer.IntAt(db.CENA)))
	xmlOffer.Add(xml.NewCdataElem("waluta", "PLN"))
	xmlOffer.Add(xml.NewElem("data_ost_modyfikacji", offer.StringAt(db.DATA_MODYFIKACJI)))
	xmlOffer.Add(xml.NewElem("data_wprowadzenia", offer.StringAt(db.DATA_ZGLOSZENIA)))

	xmlOffer.Add(xml.NewCdataElem("miejscowosc", additionals.StringAt(db.POLE_01)))
	xmlOffer.Add(xml.NewCdataElem("region", additionals.StringAt(db.POLE_02)))
	xmlOffer.Add(xml.NewCdataElem("wojewodztwo", additionals.StringAt(db.POLE_03)))
	xmlOffer.Add(xml.NewCdataElem("dzielnica", additionals.StringAt(db.POLE_04)))
	xmlOffer.Add(xml.NewCdataElem("ulica", additionals.StringAt(db.POLE_05)))
	xmlOffer.Add(xml.NewCdataElem("kod_pocztowy", additionals.StringAt(db.POLE_06)))
	xmlOffer.Add(xml.NewCdataElem("powiat", additionals.StringAt(db.POLE_10)))
	xmlOffer.Add(xml.NewCdataElem("gmina", additionals.StringAt(db.POLE_11)))

	xmlOffer.Add(xml.NewCdataElem("komunikacja", additionals.StringAt(db.POLE_07)))
	xmlOffer.Add(xml.NewCdataElem("dojazd", additionals.StringAt(db.POLE_08)))
	xmlOffer.Add(xml.NewCdataElem("polozenie", additionals.StringAt(db.POLE_09)))
	xmlOffer.Add(xml.NewElemNull("powierzchnia_calkowita", additionals.StringAt(db.POLE_60), "0"))
	xmlOffer.Add(xml.NewElemNull("powierzchnia_uzytkowa", additionals.StringAt(db.POLE_61), "0"))
	xmlOffer.Add(xml.NewElemNull("powierzchnia_dzialki", additionals.StringAt(db.POLE_62), "0"))
	xmlOffer.Add(xml.NewElemNull("szerokosc_dzialki", additionals.StringAt(db.POLE_63), "0"))
	xmlOffer.Add(xml.NewElemNull("dlugosc_dzialki", additionals.StringAt(db.POLE_64), "0"))
	xmlOffer.Add(xml.NewElemNull("kubatura", additionals.StringAt(db.POLE_65), "0"))
	xmlOffer.Add(xml.NewElemNull("pochylenie_dzialki", additionals.StringAt(db.POLE_21), "0"))
	xmlOffer.Add(xml.NewElemNull("ilosc_pokoi", additionals.StringAt(db.POLE_66), "0"))
	xmlOffer.Add(xml.NewElemNull("ilosc_pomieszczen", additionals.StringAt(db.POLE_67), "0"))
	xmlOffer.Add(xml.NewElemNull("ilosc_pieter", additionals.StringAt(db.POLE_68), "0"))
	xmlOffer.Add(xml.NewElemNull("pietro", additionals.StringAt(db.POLE_69), "0"))
	xmlOffer.Add(xml.NewElem("rodzaj_pomieszczen", additionals.StringAt(db.POLE_22)))
	xmlOffer.Add(xml.NewElem("rozklad_pomieszczen", additionals.StringAt(db.POLE_23)))
	xmlOffer.Add(xml.NewElem("rodzaj_zabudowy", additionals.StringAt(db.POLE_24)))
	xmlOffer.Add(xml.NewElemNull("rok_budowy", additionals.StringAt(db.POLE_25), "0"))
	xmlOffer.Add(xml.NewElemNull("typ_konstrukcji", additionals.StringAt(db.POLE_26), "0"))
	xmlOffer.Add(xml.NewElem("stan_techniczny", additionals.StringAt(db.POLE_27)))
	xmlOffer.Add(xml.NewElem("cena_dodatkowa_sprzedazy", additionals.StringAt(db.POLE_41)))
	xmlOffer.Add(xml.NewElem("cena_dodatkowa_najmu", additionals.StringAt(db.POLE_42)))
	xmlOffer.Add(xml.NewElem("cena_dodatkowa_czynsz", additionals.StringAt(db.POLE_43)))
	xmlOffer.Add(xml.NewElem("forma_wlasnosci", additionals.StringAt(db.POLE_44)))
	xmlOffer.Add(xml.NewElem("rodzaj_dzialki", additionals.StringAt(db.POLE_28)))
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
		offerId := offer.Id()
		go getOfferData(dbConn, offerId, specials, c)
		out <- <-c
	}
}

// getOfferData gather all data for offer and return filled object to channel
func getOfferData(dbConn *db.DB, offerId int, specialsMap db.SpecialsMap, offersChan chan *xml.Offer) {
	var (
		err error
	)
	offer := db.NewOffer(dbConn, offerId)
	utils.LogErrf(offer.Get(), "Get offer (%d)", offerId)
	// contact and owner data
	personId := offer.StringAt(db.ID_WPROWADZAJACEGO)
	contactId := offer.StringAt(db.ID_WLASCICIELA)
	person := db.NewPerson(dbConn)
	utils.LogErrf(person.Get(personId), "Get person (%s) for offer (%d)", personId, offerId)
	contact := db.NewPerson(dbConn)
	utils.LogErrf(contact.Get(contactId), "Get contact (%s) for offer (%d)", personId, offerId)
	// additionals data
	add := db.NewAdditional(dbConn)
	utils.LogErrf(add.Get(offerId), "Get additional data for offer (%d)", offerId)
	//
	xmlOffer := fillOffer(offer, add, person, contact)
	sid := offer.StrId()
	var imgIds []string
	// images and specials
	images := db.NewImages(dbConn, sid)
	if imgIds, err = images.FileNames(); err != nil {
		utils.LogErrf(err, "Images ids for offer (%d)", offerId)
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
	exportId := fmt.Sprintf("%s_00.xml", time.Now().Format(DateFormat))
	exportXMLName := filepath.Join(workDir, exportId)

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
		for _, offerId := range offers.Ids() {
			offer := db.NewOffer(dbConn, offerId)
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
		counter += 1
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
	if err = ioutil.WriteFile(filepath.Join(workDir, "last.id"), []byte(exportId), os.ModePerm); err != nil {
		return
	}
	count = len(offers.Ids())
	return
}
