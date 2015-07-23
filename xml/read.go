package xml

import (
	"encoding/xml"
	"io/ioutil"
	"regexp"
	"strings"
)

var (
	strToRem = []string{
		"<", "",
		">", "",
		"&nbsp;", "",
		"&lt;", "",
		"&gt;", "",
		"&nbsp", "",
		"&lt", "",
		"&gt", "",
		"&oacute;", "ó",
		"&Oacute;", "Ó",
		"&bdquo;", "\"",
		"&rdquo;", "\"",
		"&ndash;", "-",
		"&quot;", "\"",
		"&frac34;", "3/4",
		"&frac12;", "1/2",
		"\n", " ",
		//"&", "",
	}
	repl   = strings.NewReplacer(strToRem...)
	htmlRe = regexp.MustCompile(`<[^>]*>`)
)

//
// Structs to read content from intermadiate XML files
//
type XmlHeader struct {
	XMLName       xml.Name   `xml:"swdp"`
	ExportVersion string     `xml:"system>wersja_eksportu"`
	ExportDate    string     `xml:"system>data_eksportu"`
	ExportHour    string     `xml:"system>godzina_eksportu"`
	WWW           string     `xml:"system>strona_www"`
	Offers        []XmlOffer `xml:"oferty_nieruchomosci>oferta"`
}

type XmlSpecials struct {
	XMLName  xml.Name `xml:"wyroznienia"`
	Specials []string `xml:"wyroznienie"`
}

type XmlPictures struct {
	XMLName  xml.Name `xml:"zdjecia"`
	Pictures []string `xml:"zdjecie"`
}

func (x XmlPictures) Len() int {
	return len(x.Pictures)
}

type XmlCdata struct {
	Data string `xml:",chardata"`
}

func (d XmlCdata) Clean() string {
	return strings.TrimSpace(repl.Replace(d.Data))
}

func (d XmlCdata) NoHtml() string {
	return strings.TrimSpace(repl.Replace(htmlRe.ReplaceAllString(d.Data, "")))
	//return htmlRe.ReplaceAllString(d.Data, "")
}

type XmlOffer struct {
	XMLName     xml.Name `xml:"oferta"`
	Id          int64    `xml:"id_oferty"`
	NazwaOferty XmlCdata `xml:"nazwa_oferty"`
	Grupa       XmlCdata `xml:"grupa"`
	Podgrupa    XmlCdata `xml:"podgrupa"`
	NumerOferty XmlCdata `xml:"numer_oferty"`
	Opis        XmlCdata `xml:"opis"`

	Uwagi         XmlCdata `xml:"uwagi"`
	OfertaZawiera XmlCdata `xml:"oferta_zawiera"`
	Wprowadzajacy XmlCdata `xml:"wprowadzajacy"`
	KontaktOsoba  XmlCdata `xml:"kontakt_osoba"`
	KontaktTel1   XmlCdata `xml:"kontakt_telefon_1"`
	KontaktTel2   XmlCdata `xml:"kontakt_telefon_2"`
	KontaktTel3   XmlCdata `xml:"kontakt_telefon_3"`
	KontaktEmail  XmlCdata `xml:"kontakt_email"`

	CenaGroszy string `xml:"cena_groszy"`
	Waluta     string `xml:"waluta"`

	DataModyfikacji  string `xml:"data_ost_modyfikacji"`
	DataWprowadzenia string `xml:"data_wprowadzenia"`

	Miejscowosc XmlCdata `xml:"miejscowosc"`
	Region      XmlCdata `xml:"region"`
	Powiat      XmlCdata `xml:"powiat"`
	Community   XmlCdata `xml:"gmina"`
	Wojewodztwo XmlCdata `xml:"wojewodztwo"`
	Dzielnica   XmlCdata `xml:"dzielnica"`
	Ulica       XmlCdata `xml:"ulica"`
	KodPocztowy XmlCdata `xml:"kod_pocztowy"`
	Komunikacja XmlCdata `xml:"komunikacja"`
	Dojazd      XmlCdata `xml:"dojazd"`
	Polozenie   XmlCdata `xml:"polozenie"`

	PowierzchniaCalkowita  string   `xml:"powierzchnia_calkowita"`
	PowierzchniaUzytkowa   string   `xml:"powierzchnia_uzytkowa"`
	PowierzchniaDzialki    string   `xml:"powierzchnia_dzialki"`
	SzerokoscDzialki       string   `xml:"szerokosc_dzialki"`
	DlugoscDzialki         string   `xml:"dlugosc_dzialki"`
	Kubatura               string   `xml:"kubatura"`
	Nachylenie             XmlCdata `xml:"pochylenie_dzialki"`
	IloscPokoi             string   `xml:"ilosc_pokoi"`
	IloscPomieszczen       string   `xml:"ilosc_pomieszczen"`
	IloscPieter            string   `xml:"ilosc_pieter"`
	Pietro                 string   `xml:"pietro"`
	RodzajPomieszczen      XmlCdata `xml:"rodzaj_pomieszczen"`
	RozkladPomieszczen     XmlCdata `xml:"rozklad_pomieszczen"`
	RodzajZabudowy         XmlCdata `xml:"rodzaj_zabudowy"`
	RokBudowy              XmlCdata `xml:"rok_budowy"`
	TypKonstrukcji         XmlCdata `xml:"typ_konstrukcji"`
	StanTechniczny         XmlCdata `xml:"stan_techniczny"`
	CenaDodatkowaSprzedazy string   `xml:"cena_dodatkowa_sprzedazy"`
	CenaDodatkowaNajmu     string   `xml:"cena_dodatkowa_najmu"`
	CenaDodatkowaCzynsz    string   `xml:"cena_dodatkowa_czynsz"`
	FormaWlasnosci         XmlCdata `xml:"forma_wlasnosci"`
	RodzajDzialki          XmlCdata `xml:"rodzaj_dzialki"`
	Wylacznosc             XmlCdata `xml:"wylacznosc"`

	Specials *XmlSpecials
	Pictures *XmlPictures
}

func Unmarshall(path string) (header *XmlHeader, err error) {
	var (
		data []byte
	)
	if data, err = ioutil.ReadFile(path); err != nil {
		return
	}
	header = &XmlHeader{}
	err = xml.Unmarshal([]byte(data), header)
	return
}
