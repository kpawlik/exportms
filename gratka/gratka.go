package gratka

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/kpawlik/exportms/utils"
	mxml "github.com/kpawlik/exportms/xml"
)

const (
	zipNameDateFormat = "20060102"
	// ExportType type of export
	ExportType = "full"
)

type fillOfferFunc func(*record, mxml.XmlOffer, *Dicts) error

var (
	fillFuncMap = map[string]fillOfferFunc{
		"mieszkanie": fillApartment,
		"dom":        fillHouse,
		"dzialka":    fillParcele,
		"lokal":      fillLocal,
		"garaz":      fillOffice, // same fields as office
		"biuro":      fillOffice}
	pictChan = make(chan error)
	// wg       sync.WaitGroup
	// mutex    sync.Mutex
)

type actions struct {
	*mxml.BaseElem
}

func newActions(exportType string) *actions {
	elem := newXMLElem("export", "full")
	return &actions{mxml.NewBaseElem("actions", []mxml.ElemWriter{elem})}
}

type company struct {
	*mxml.BaseElem
}

func newCompany(code string) *company {
	elem := newXMLElem("kod_offline", code)
	return &company{mxml.NewBaseElem("firma", []mxml.ElemWriter{elem})}
}

type pictures struct {
	*mxml.BaseElem
}

// Gratka type and methods
type Gratka struct {
	*mxml.BaseElem
}

func newGratka() *Gratka {
	return &Gratka{mxml.NewBaseElem("dane", []mxml.ElemWriter{})}
}

type record struct {
	*mxml.BaseElem
}

func newRecord() *record {
	return &record{mxml.NewBaseElem("record", []mxml.ElemWriter{})}
}

func newXMLElem(tag, value string) *mxml.Elem {
	return mxml.NewElem(tag, value)
}

//Convert converts source dir contetnt into GRATKA format
func Convert(sourceFolder, offlineCode, domain string) (gratkaDir string, err error) {
	var (
		header *mxml.XmlHeader
		inPath string
	)
	log.Printf("Start convert to 'GRATKA' format\n")
	if inPath = utils.GetLastIDFileName(sourceFolder); len(inPath) == 0 {
		err = utils.Errorf("No file for last export")
		return
	}
	if header, err = mxml.Unmarshall(inPath); err != nil {
		return
	}
	if gratkaDir, err = makeDirs(sourceFolder); err != nil {
		return
	}
	imagesDirPath := filepath.Join(sourceFolder, "images")
	outFilePath := filepath.Join(gratkaDir, "dane.xml")

	if err = convertXML(header, offlineCode, imagesDirPath, outFilePath, gratkaDir); err != nil {
		return
	}
	log.Printf("Data converted\n")
	return
}

// ZipPath returns apth to archive file
func ZipPath(sourceDir, domain string) string {
	zipFileName := fmt.Sprintf("%s_%s.zip", domain, time.Now().Format(zipNameDateFormat))
	return filepath.Join(sourceDir, zipFileName)
}

//DirPath path to folder with gratka
func DirPath(baseDir string) string {
	return filepath.Join(baseDir, "gratka")
}

//Zip compress folder
func Zip(sourceFolder, gratkaFolder, domain string) (zipPath string, err error) {
	zipPath = ZipPath(sourceFolder, domain)
	log.Printf("Zip start. Directory '%s' in to '%s'\n", gratkaFolder, zipPath)
	if err = utils.ZipFolder(gratkaFolder, zipPath); err != nil {
		return
	}
	log.Printf("Zip end")
	return
}

func makeDirs(baseFolder string) (path string, err error) {
	path = DirPath(baseFolder)
	if err = os.RemoveAll(path); err != nil {
		return
	}
	err = os.MkdirAll(path, os.ModePerm)
	return
}

func addPictures(rec *record, offer mxml.XmlOffer, imgFolder, gratkaFolder string) (err error) {
	picts := &pictures{mxml.NewBaseElem("zdjecia", []mxml.ElemWriter{})}
	if imgLen := offer.Pictures.Len(); imgLen == 0 {
		rec.Add(picts)
		pictChan <- err
		return
	}
	var (
		pictureOrder     int
		srcPath, dstPath string
	)
	for i, pictureName := range offer.Pictures.Pictures {
		pictureOrder = i + 1
		srcPath = filepath.Join(imgFolder, pictureName)
		dstPath = filepath.Join(gratkaFolder, fmt.Sprintf("%d_%d.jpg", offer.Id, pictureOrder))
		if err = copyFile(srcPath, dstPath); err != nil {
			break
		}
		picts.Add(newXMLElem(fmt.Sprintf("z%d", pictureOrder), "true"))
	}
	// mutex.Lock()
	rec.Add(picts)
	// mutex.Unlock()
	pictChan <- err
	return
}

func convertXML(header *mxml.XmlHeader, offlineCode, imagesDirPath, outFilePath, gratkaFolder string) (err error) {
	var (
		rec    *record
		file   *os.File
		gratka *Gratka
	)
	dicts := NewDicts()
	actions := newActions(ExportType)
	company := newCompany(offlineCode)
	gratka = newGratka()
	gratka.Add(actions)
	gratka.Add(company)
	for _, offer := range header.Offers {
		rec = newRecord()
		if err = fillCommon(rec, offer, dicts); err != nil {
			return
		}
		offerType := unifyGroup(offer.Podgrupa.Clean())
		if fillFunc, ok := fillFuncMap[offerType]; ok {
			fillFunc(rec, offer, dicts)
		} else {
			log.Panicf("Unrecognized offer type '%s' offer no '%s'\n", offerType, offer.NumerOferty.Clean())
		}
		go addPictures(rec, offer, imagesDirPath, gratkaFolder)
		gratka.Add(rec)

	}
	for range header.Offers {
		if err = <-pictChan; err != nil {
			return
		}
	}
	if file, err = os.Create(outFilePath); err != nil {
		return
	}
	defer file.Close()
	file.WriteString(xml.Header)
	gratka.Write(file)
	return
}

func fillCommon(rec *record, offer mxml.XmlOffer, dicts *Dicts) (err error) {
	insertTime := time.Now()
	expireTime := insertTime.Add(time.Duration(2 * 24 * time.Hour))
	rec.Add(newXMLElem("action", "replace"))
	rec.Add(newXMLElem("id_inspert", strconv.Itoa(int(offer.Id))))
	rec.Add(newXMLElem("numer_oferty", fmt.Sprintf("%sSW", offer.NumerOferty.Clean())))
	//rec.Add(NewXmlElem("data_zalozenia", insertTime.Format(DateTimeFormat)))
	rec.Add(newXMLElem("data_usuniecia", expireTime.Format(DateFormat)))
	rec.Add(newXMLElem("id_rubryka", getSectionId(offer.Grupa.Clean(), "rubryka")))
	rec.Add(newXMLElem("id_podrubryka", getSectionId(offer.Podgrupa.Clean(), "podrubryka")))
	if reg, ok := dicts.Get("region").GetVal(offer.Wojewodztwo.Clean()); ok {
		rec.Add(newXMLElem("id_region", reg))
	} else {
		log.Printf("Warning: code for region '%s' not found\n", offer.Wojewodztwo.Clean())
	}
	rec.Add(newXMLElem("miejscowosc", offer.Miejscowosc.Clean()))
	rec.Add(newXMLElem("dzielnica", offer.Dzielnica.Clean()))
	rec.Add(newXMLElem("ulica", offer.Ulica.Clean()))
	rec.Add(newXMLElem("id_jednostka_pow", strconv.Itoa(1)))
	if cur, ok := dicts.Get("currency").GetVal(offer.Waluta); ok {
		rec.Add(newXMLElem("id_waluta", cur))
	}
	rec.Add(newXMLElem("opis", description(offer, dicts)))
	rec.Add(newXMLElem("kontakt_email", offer.KontaktEmail.Clean()))
	rec.Add(newXMLElem("kontakt_telefon", phones(offer)))
	rec.Add(newXMLElem("kontakt_osoba", offer.KontaktOsoba.Clean()))
	rec.Add(newXMLElem("kod_pocztowy", offer.KodPocztowy.Clean()))
	if exclusive := offer.Wylacznosc.Clean(); len(exclusive) > 0 {
		rec.Add(newXMLElem("na_wylacznosc", "T"))
	}
	return
}

func description(o mxml.XmlOffer, dicts *Dicts) string {
	var (
		out  = ""
		strs = [][]string{
			[]string{"Nazwa oferty: %s ", o.NazwaOferty.Clean()},
			[]string{"Opis oferty: %s ", o.Opis.NoHtml()},
			[]string{"Oferta zawiera %s ", o.OfertaZawiera.Clean()},
			[]string{"Uwagi: %s ", o.Uwagi.Clean()},
			[]string{"Komunikacja: %s ", o.Komunikacja.Clean()},
			[]string{"Dojazd: %s ", o.Dojazd.Clean()},
			[]string{"Położenie: %s ", o.Polozenie.Clean()},
			[]string{"Rodzaj pomieszczeń: %s ", o.RodzajPomieszczen.Clean()},
			[]string{"Rozkład pomieszczeń: %s ", o.RozkladPomieszczen.Clean()},
			[]string{"Typ zabudowy: %s ", o.RodzajZabudowy.Clean()},
			[]string{"Stan techniczny: %s ", o.StanTechniczny.Clean()},
			[]string{"Rodzaj zabudowy: %s ", o.RodzajZabudowy.Clean()},
			[]string{"Wylaczność %s ", o.Wylacznosc.Clean()}}
	)

	for _, str := range strs {
		if len(str[1]) > 0 && str[1] != "nieistotne" {
			out = fmt.Sprintf("%s %s", out, fmt.Sprintf(str[0], str[1]))
		}
	}
	out = strings.Replace(out, "&", "", -1)
	return out
}

func phones(o mxml.XmlOffer) string {
	out := o.KontaktTel1.Clean()
	for _, s := range []string{
		o.KontaktTel2.Clean(),
		o.KontaktTel3.Clean()} {
		if len(s) > 0 {
			out += s
		}
	}
	return out
}

func fillApartment(rec *record, offer mxml.XmlOffer, dicts *Dicts) (err error) {
	offerType := "mieszkanie"
	rec.Add(newXMLElem("id_pietro", getFloor(offer.Pietro)))
	rec.Add(newXMLElem("id_liczba_pieter", getNoOfFloors(offer.IloscPieter, offerType)))
	rec.Add(newXMLElem("powierzchnia", formatFloat(offer.Area("apartment"))))
	rec.Add(newXMLElem("id_liczba_pokoi", getNoOfRooms(offer.IloscPokoi, offerType)))
	rec.Add(newXMLElem("id_rok_budowy", offer.RokBudowy.Clean()))
	if bType, ok := dicts.Get("building").GetVal(offer.RodzajZabudowy.Clean()); ok {
		rec.Add(newXMLElem("id_typ_zabudowy_nowe", bType))
	}
	switch getSectionId(offer.Grupa.Clean(), "rubryka") {
	case "1": //sprzedaz
		rec.Add(newXMLElem("cena_calkowita", formatPrice(offer.CenaGroszy)))
		rec.Add(newXMLElem("wysokosc_czynszu", formatPrice(offer.CenaDodatkowaCzynsz)))
	case "5": //wynajem
		rec.Add(newXMLElem("cena", formatPrice(offer.CenaGroszy)))
		rec.Add(newXMLElem("cena_wynajmu", formatPrice(offer.CenaGroszy)))
		rec.Add(newXMLElem("cena_czynsz", formatPrice(offer.CenaGroszy)))
	}
	return
}

func fillHouse(rec *record, offer mxml.XmlOffer, dicts *Dicts) (err error) {
	offerType := "dom"
	rec.Add(newXMLElem("id_liczba_pieter", getNoOfFloors(offer.IloscPieter, offerType)))
	rec.Add(newXMLElem("id_liczba_pokoi", getNoOfRooms(offer.IloscPokoi, offerType)))
	rec.Add(newXMLElem("powierzchnia", formatFloat(offer.PowierzchniaCalkowita)))
	rec.Add(newXMLElem("id_rok_budowy", offer.RokBudowy.Clean()))
	rec.Add(newXMLElem("powierzchnia_dzialki", formatFloat(offer.PowierzchniaDzialki)))
	if bType, ok := dicts.Get("buildingtype").GetVal(offer.TypKonstrukcji.Clean()); ok {
		rec.Add(newXMLElem("id_typ_budynku", bType))
	}
	rec.Add(newXMLElem("dlugosc", formatFloat(offer.DlugoscDzialki)))
	rec.Add(newXMLElem("szerokosc", formatFloat(offer.DlugoscDzialki)))
	rec.Add(newXMLElem("cena", formatPrice(offer.CenaGroszy)))
	return
}

func fillParcele(rec *record, offer mxml.XmlOffer, dicts *Dicts) (err error) {

	if pType, ok := dicts.Get("parcele").GetVal(offer.RodzajDzialki.Clean()); ok {
		rec.Add(newXMLElem("id_rodzaj_dzialki", pType))
	}
	rec.Add(newXMLElem("powierzchnia", formatFloat(offer.Area("parcele"))))
	rec.Add(newXMLElem("dlugosc", formatFloat(offer.DlugoscDzialki)))
	rec.Add(newXMLElem("szerokosc", formatFloat(offer.DlugoscDzialki)))
	rec.Add(newXMLElem("id_rok_budowy", offer.RokBudowy.Clean()))
	rec.Add(newXMLElem("cena", formatPrice(offer.CenaGroszy)))
	return
}

func fillLocal(rec *record, offer mxml.XmlOffer, dicts *Dicts) (err error) {
	rec.Add(newXMLElem("id_rok_budowy", offer.RokBudowy.Clean()))
	rec.Add(newXMLElem("powierzchnia_calkowita", formatFloat(offer.Area("local"))))
	if loc, ok := dicts.Get("local").GetVal(offer.TypKonstrukcji.Clean()); ok {
		rec.Add(newXMLElem("id_typ_lokalu", loc))
	}
	rec.Add(newXMLElem("cena", formatPrice(offer.CenaGroszy)))
	return
}

func fillOffice(rec *record, offer mxml.XmlOffer, dicts *Dicts) (err error) {
	rec.Add(newXMLElem("id_rok_budowy", offer.RokBudowy.Clean()))
	rec.Add(newXMLElem("powierzchnia", formatFloat(offer.PowierzchniaCalkowita)))
	rec.Add(newXMLElem("cena", formatFloat(offer.CenaGroszy)))
	return
}
