package xml

import (
	"testing"
)

func TestGetSectionID(t *testing.T) {
	offer := &XmlOffer{Uwagi: XmlCdata{Data: "&#xA;Ala ma kota<>&"}}
	if res := offer.Uwagi.Clean(); res != "\nAla ma kota" {
		t.Errorf("Replacer work wrong '%s'", res)
	}
}
