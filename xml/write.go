package xml

import (
	"fmt"
	"io"
	"os"
)

//
// Structures to dump data as XML imtermediate.
//

//************************************************************

type ElemWriter interface {
	Write(io.Writer)
	ElemName() string
	Elems() []ElemWriter
}

//************************************************************
type BaseElem struct {
	XMLName string
	elems   []ElemWriter
}

func NewBaseElem(name string, elems []ElemWriter) *BaseElem {
	if elems == nil {
		elems = []ElemWriter{}
	}
	return &BaseElem{name, elems}
}

func (b *BaseElem) ElemName() string {
	return b.XMLName
}

func (b *BaseElem) Elems() []ElemWriter {
	return b.elems
}

func (b *BaseElem) Write(w io.Writer) {
	WriteElements(b, w)
}

func (b *BaseElem) writeBeginElemt(w io.Writer) {
	fmt.Fprintf(w, "<%s>\n", b.ElemName())
}

func (b *BaseElem) writeEndElement(w io.Writer) {
	fmt.Fprintf(w, "</%s>\n", b.ElemName())
}

func (b *BaseElem) writeElements(w io.Writer) {
	for _, e := range b.Elems() {
		e.Write(w)
	}
}
func (b *BaseElem) Add(elem ElemWriter) {
	b.elems = append(b.elems, elem)
}

//************************************************************

type Elem struct {
	*BaseElem
	value     string
	cdata     bool
	nullValue string
}

func NewElem(tag, value string) *Elem {
	return &Elem{NewBaseElem(tag, nil), value, false, ""}
}

func NewCdataElem(tag, value string) *Elem {
	return &Elem{NewBaseElem(tag, nil), value, true, ""}
}

func NewElemNull(tag, value string, null string) *Elem {
	return &Elem{NewBaseElem(tag, nil), value, false, null}
}

func (e *Elem) ElemName() string {
	return e.BaseElem.ElemName()
}

func (e *Elem) Write(w io.Writer) {
	if len(e.value) == 0 || e.value == e.nullValue {
		return
	}
	fmt.Fprintf(w, "<%s>", e.ElemName())
	if e.cdata {
		fmt.Fprintf(w, "<![CDATA[%s]]>", e.value)
	} else {
		fmt.Fprintf(w, "%s", e.value)
	}
	fmt.Fprintf(w, "</%s>\n", e.ElemName())
}

//************************************************************

type ListElem struct {
	*BaseElem
	TagName string
}

func NewListElem(xmlName, tagName string) *ListElem {
	return &ListElem{
		NewBaseElem(xmlName, nil),
		tagName}
}

func (l *ListElem) Add(elem string) {
	l.elems = append(l.elems, NewElem(l.TagName, elem))
}

func (l *ListElem) AddMany(elems ...string) {
	for _, elem := range elems {
		l.Add(elem)
	}
}

//************************************************************

type System struct {
	*BaseElem
}

func NewSystem(version, date, hour, www string) *System {
	return &System{
		NewBaseElem("system", []ElemWriter{
			NewElem("wersja_eksportu", version),
			NewElem("data_eksportu", date),
			NewElem("godzina_eksportu", hour),
			NewElem("strona_www", www)})}
}

//************************************************************

type Header struct {
	*BaseElem
}

func NewHeader(system, offers ElemWriter) *Header {
	return &Header{NewBaseElem("swdp", []ElemWriter{system, offers})}
}

//************************************************************

type Offers struct {
	*BaseElem
}

func NewOffers() *Offers {
	return &Offers{NewBaseElem("oferty_nieruchomosci", nil)}
}

func (o *Offers) Add(offer *Offer) {
	o.elems = append(o.elems, offer)
}

func (o *Offers) Write(w io.Writer) {
	fmt.Fprintf(w, "<%s>\n", o.XMLName)
	for _, elem := range o.elems {
		elem.Write(w)
	}
	fmt.Fprintf(w, "</%s>\n", o.XMLName)
}

type Offer struct {
	*BaseElem
	Specials *ListElem
	Pictures *ListElem
}

func NewOffer() *Offer {
	return &Offer{BaseElem: NewBaseElem("oferta", nil)}
}

func (o *Offer) Add(elem *Elem) {
	o.elems = append(o.elems, elem)
}

func (o *Offer) Write(w io.Writer) {
	o.writeBeginElemt(w)
	o.writeElements(w)
	o.Specials.Write(w)
	o.Pictures.Write(w)
	o.writeEndElement(w)
}

func WriteElements(elem ElemWriter, w io.Writer) {
	fmt.Fprintf(w, "<%s>\n", elem.ElemName())
	for _, e := range elem.Elems() {
		e.Write(w)
	}
	fmt.Fprintf(w, "</%s>\n", elem.ElemName())
}

func Write(header *Header, path string) (err error) {
	var (
		f *os.File
	)
	if f, err = os.Create(path); err != nil {
		return
	}
	defer f.Close()
	f.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	header.Write(f)
	return
}
