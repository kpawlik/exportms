package gratka

import (
	"testing"
	"time"
)

func TestGetSectionID(t *testing.T) {
	testData := map[string]string{
		"Sprzedaż": "1",
		"sprzedaż": "1",
		"sprzedaz": "1",
		"dowynaj":  "4",
		"DoWynaj":  "4",
		"Kupno":    "2",
		"kupno":    "2",
		"wynaj":    "5",
		"Wynaj":    "5",
		"Wynajem":  "5",
	}
	for k, v := range testData {

		if sID := getSectionId(k, ExportType); sID != v {
			t.Errorf("Wrong Section Id. Is %d, should be %d\n", sID, v)
		}
	}
}

func TestTimeDiff(t *testing.T) {
	t1 := time.Now().Format(DateTimeFormat)
	t2 := time.Now().Add(2 * 24 * time.Hour).Format(DateFormat)
	t.Logf("D1: %s, D2: %s\n", t1, t2)
}
