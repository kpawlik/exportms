package gratka

import (
	"encoding/xml"
	"fmt"
	"log"
	"github.com/kpawlik/exportms/utils"
	mxml "github.com/kpawlik/exportms/xml"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	zipNameDateFormat = "20060102"
	ExportType        = "full"
)

type fillOfferFunc func(*Record, mxml.XmlOffer, *Dicts) error

var (
	fillFuncMap = map[string]fillOfferFunc{
		"mieszkanie": fillApartment,
		"dom":        fillHouse,
		"dzialka":    fillParcele,
		"lokal":      fillLocal,
		"garaz":      fillOffice, // same fields as office
		"biuro":      fillOffice}
	pictChan = make(chan error)
)

type Actions struct {
	*mxml.BaseElem
}

func NewActions(exportType string) *Actions {
	elem := NewXmlElem("export", "full")
	return &Actions{mxml.NewBaseElem("actions", []mxml.ElemWriter{elem})}
}

type Company struct {
	*mxml.BaseElem
}

func NewCompany(code string) *Company {
	elem := NewXmlElem("kod_offline", code)
	return &Company{mxml.NewBaseElem("firma", []mxml.ElemWriter{elem})}
}

type Pictures struct {
	*mxml.BaseElem
}

type Gratka struct {
	*mxml.BaseElem
}

func NewGratka() *Gratka {
	return &Gratka{mxml.NewBaseElem("dane", []mxml.ElemWriter{})}
}

type Record struct {
	*mxml.BaseElem
}

func NewRecord() *Record {
	return &Record{mxml.NewBaseElem("record", []mxml.ElemWriter{})}
}

func NewXmlElem(tag, value string) *mxml.Elem {
	return mxml.NewElem(tag, value)
}

func Convert(sourceFolder, offlineCode, domain string) (gratkaDir string, err error) {
	var (
		header *mxml.XmlHeader
	)
	log.Printf("Start convert to 'GRATKA' format\n")
	if inPath := utils.GetLastIdFileName(sourceFolder); len(inPath) == 0 {
		err = utils.Errorf("No file for last export")
		return
	} else {
		if header, err = mxml.Unmarshall(inPath); err != nil {
			return
		}
	}
	if gratkaDir, err = makeDirs(sourceFolder); err != nil {
		return
	}
	imagesDirPath := filepath.Join(sourceFolder, "images")
	outFilePath := filepath.Join(gratkaDir, "dane.xml")

	if err = convertXml(header, offlineCode, imagesDirPath, outFilePath, gratkaDir); err != nil {
		return
	}
	log.Printf("Data converted\n")

	//	if err = SendToFtp(domain, zipPath, ftp); err != nil {
	//		return
	//	}

	return
}

func ZipPath(sourceDir, domain string) string {
	zipFileName := fmt.Sprintf("%s_%s.zip", domain, time.Now().Format(zipNameDateFormat))
	return filepath.Join(sourceDir, zipFileName)
}

func DirPath(baseDir string) string {
	return filepath.Join(baseDir, "gratka")
}

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

func addPictures(rec *Record, offer mxml.XmlOffer, imgFolder, gratkaFolder string) (err error) {
	pictures := &Pictures{mxml.NewBaseElem("zdjecia", []mxml.ElemWriter{})}
	if imgLen := offer.Pictures.Len(); imgLen == 0 {
		rec.Add(pictures)
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
		pictures.Add(NewXmlElem(fmt.Sprintf("z%d", pictureOrder), "true"))
	}
	rec.Add(pictures)
	pictChan <- err
	return
}

func convertXml(header *mxml.XmlHeader, offlineCode, imagesDirPath, outFilePath, gratkaFolder string) (err error) {
	var (
		rec    *Record
		file   *os.File
		gratka *Gratka
	)
	dicts := NewDicts()
	actions := NewActions(ExportType)
	company := NewCompany(offlineCode)
	gratka = NewGratka()
	gratka.Add(actions)
	gratka.Add(company)
	for _, offer := range header.Offers {
		rec = NewRecord()
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

func fillCommon(rec *Record, offer mxml.XmlOffer, dicts *Dicts) (err error) {
	insertTime := time.Now()
	expireTime := insertTime.Add(time.Duration(2 * 24 * time.Hour))
	rec.Add(NewXmlElem("action", "replace"))
	rec.Add(NewXmlElem("id_inspert", strconv.Itoa(int(offer.Id))))
	rec.Add(NewXmlElem("numer_oferty", fmt.Sprintf("%sSW", offer.NumerOferty.Clean())))
	//rec.Add(NewXmlElem("data_zalozenia", insertTime.Format(DateTimeFormat)))
	rec.Add(NewXmlElem("data_usuniecia", expireTime.Format(DateFormat)))
	rec.Add(NewXmlElem("id_rubryka", getSectionId(offer.Grupa.Clean(), "rubryka")))
	rec.Add(NewXmlElem("id_podrubryka", getSectionId(offer.Podgrupa.Clean(), "podrubryka")))
	if reg, ok := dicts.Get("region").GetVal(offer.Wojewodztwo.Clean()); ok {
		rec.Add(NewXmlElem("id_region", reg))
	} else {
		log.Printf("Warning: code for region '%s' not found\n", offer.Wojewodztwo.Clean())
	}
	rec.Add(NewXmlElem("miejscowosc", offer.Miejscowosc.Clean()))
	rec.Add(NewXmlElem("dzielnica", offer.Dzielnica.Clean()))
	rec.Add(NewXmlElem("ulica", offer.Ulica.Clean()))
	rec.Add(NewXmlElem("id_jednostka_pow", strconv.Itoa(1)))
	if cur, ok := dicts.Get("currency").GetVal(offer.Waluta); ok {
		rec.Add(NewXmlElem("id_waluta", cur))
	}
	rec.Add(NewXmlElem("opis", description(offer, dicts)))
	rec.Add(NewXmlElem("kontakt_email", offer.KontaktEmail.Clean()))
	rec.Add(NewXmlElem("kontakt_telefon", phones(offer)))
	rec.Add(NewXmlElem("kontakt_osoba", offer.KontaktOsoba.Clean()))
	rec.Add(NewXmlElem("kod_pocztowy", offer.KodPocztowy.Clean()))
	if exclusive := offer.Wylacznosc.Clean(); len(exclusive) > 0 {
		rec.Add(NewXmlElem("na_wylacznosc", "T"))
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

func fillApartment(rec *Record, offer mxml.XmlOffer, dicts *Dicts) (err error) {
	offerType := "mieszkanie"
	rec.Add(NewXmlElem("id_pietro", getFloor(offer.Pietro)))
	rec.Add(NewXmlElem("id_liczba_pieter", getNoOfFloors(offer.IloscPieter, offerType)))
	rec.Add(NewXmlElem("powierzchnia", formatFloat(offer.PowierzchniaCalkowita)))
	rec.Add(NewXmlElem("id_liczba_pokoi", getNoOfRooms(offer.IloscPokoi, offerType)))
	rec.Add(NewXmlElem("id_rok_budowy", offer.RokBudowy.Clean()))
	if bType, ok := dicts.Get("building").GetVal(offer.RodzajZabudowy.Clean()); ok {
		rec.Add(NewXmlElem("id_typ_zabudowy_nowe", bType))
	}
	switch getSectionId(offer.Grupa.Clean(), "rubryka") {
	case "1": //sprzedaz
		rec.Add(NewXmlElem("cena_calkowita", formatPrice(offer.CenaGroszy)))
		rec.Add(NewXmlElem("wysokosc_czynszu", formatPrice(offer.CenaDodatkowaCzynsz)))
	case "5": //wynajem
		rec.Add(NewXmlElem("cena", formatPrice(offer.CenaGroszy)))
		rec.Add(NewXmlElem("cena_wynajmu", formatPrice(offer.CenaGroszy)))
		rec.Add(NewXmlElem("cena_czynsz", formatPrice(offer.CenaGroszy)))
	}
	return
}

func fillHouse(rec *Record, offer mxml.XmlOffer, dicts *Dicts) (err error) {
	offerType := "dom"
	rec.Add(NewXmlElem("id_liczba_pieter", getNoOfFloors(offer.IloscPieter, offerType)))
	rec.Add(NewXmlElem("id_liczba_pokoi", getNoOfRooms(offer.IloscPokoi, offerType)))
	rec.Add(NewXmlElem("powierzchnia", formatFloat(offer.PowierzchniaCalkowita)))
	rec.Add(NewXmlElem("id_rok_budowy", offer.RokBudowy.Clean()))
	rec.Add(NewXmlElem("powierzchnia_dzialki", formatFloat(offer.PowierzchniaDzialki)))
	if bType, ok := dicts.Get("buildingtype").GetVal(offer.TypKonstrukcji.Clean()); ok {
		rec.Add(NewXmlElem("id_typ_budynku", bType))
	}
	rec.Add(NewXmlElem("dlugosc", formatFloat(offer.DlugoscDzialki)))
	rec.Add(NewXmlElem("szerokosc", formatFloat(offer.DlugoscDzialki)))
	rec.Add(NewXmlElem("cena", formatPrice(offer.CenaGroszy)))
	return
}

func fillParcele(rec *Record, offer mxml.XmlOffer, dicts *Dicts) (err error) {
	var (
		pow string
	)
	if pType, ok := dicts.Get("parcele").GetVal(offer.RodzajDzialki.Clean()); ok {
		rec.Add(NewXmlElem("id_rodzaj_dzialki", pType))
	}
	if len(offer.PowierzchniaDzialki) > 0 {
		pow = offer.PowierzchniaDzialki
	} else {
		pow = offer.PowierzchniaCalkowita
	}
	rec.Add(NewXmlElem("powierzchnia", formatFloat(pow)))
	rec.Add(NewXmlElem("dlugosc", formatFloat(offer.DlugoscDzialki)))
	rec.Add(NewXmlElem("szerokosc", formatFloat(offer.DlugoscDzialki)))
	rec.Add(NewXmlElem("id_rok_budowy", offer.RokBudowy.Clean()))
	rec.Add(NewXmlElem("cena", formatPrice(offer.CenaGroszy)))
	return
}

func fillLocal(rec *Record, offer mxml.XmlOffer, dicts *Dicts) (err error) {
	rec.Add(NewXmlElem("id_rok_budowy", offer.RokBudowy.Clean()))
	rec.Add(NewXmlElem("powierzchnia_calkowita", formatFloat(offer.PowierzchniaCalkowita)))
	if loc, ok := dicts.Get("local").GetVal(offer.TypKonstrukcji.Clean()); ok {
		rec.Add(NewXmlElem("id_typ_lokalu", loc))
	}
	rec.Add(NewXmlElem("cena", formatPrice(offer.CenaGroszy)))
	return
}

func fillOffice(rec *Record, offer mxml.XmlOffer, dicts *Dicts) (err error) {
	rec.Add(NewXmlElem("id_rok_budowy", offer.RokBudowy.Clean()))
	rec.Add(NewXmlElem("powierzchnia", formatFloat(offer.PowierzchniaCalkowita)))
	rec.Add(NewXmlElem("cena", formatFloat(offer.CenaGroszy)))
	return
}
