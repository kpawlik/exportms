package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kpawlik/exportms/utils"
	"path/filepath"
)

const (
	OffersQ = `SELECT ID_OFERTY FROM OFERTY_WWW WHERE AKTYWNA='1' AND ID_OFERTY NOT IN
	(SELECT DISTINCT ID_OFERTY FROM OFERTY_WWW_WYROZNIENIA WHERE ID_WYROZNIENIA = 'NIE_EKSPORTUJ') `
	SpecialsQ = `SELECT ID_OFERTY,
					ID_WYROZNIENIA FROM OFERTY_WWW_WYROZNIENIA`
	OfferQ = `SELECT ID_OFERTY,
						ID_FIRMY,
					    NUMER_OFERTY,
					    AKTYWNA,
					    TYP_OFERTY,
					    RODZAJ_OFERTY,
					    ID_GRUPY,
					    NAZWA_OFERTY,
					    INFO,
					    TRESC_OFERTY,
					    OFERTA_ZAWIERA,
					    UWAGI,
					    CENA,
					    DATA_ZGLOSZENIA,
					    GODZINA_ZGLOSZENIA,
					    ID_WPROWADZAJACEGO,
					    DATA_MODYFIKACJI,
					    GODZINA_MODYFIKACJI,
					    ID_MODYFIKUJACEGO,
					    ILOSC_WYSWIETLEN,
					    CENA_NETTO,
					    STAWKA_VAT,
					    SWW,
					    JEDNOSTKA_MIARY,
					    WAGA,
					    ID_WLASCICIELA,
					    KOD_OFERTY
					FROM OFERTY_WWW WHERE ID_OFERTY = ?`
	AdditionslQ = `SELECT POLE_01,
						POLE_02,
						POLE_03,
						POLE_04,
						POLE_05,
						POLE_06,
						POLE_07,
						POLE_08,
						POLE_09,
						POLE_10,
						POLE_11,
						POLE_60,
						POLE_61,
						POLE_62,
						POLE_63,
						POLE_64,
						POLE_65,
						POLE_21,
						POLE_66,
						POLE_67,
						POLE_68,
						POLE_69,
						POLE_22,
						POLE_23,
						POLE_24,
						POLE_25,
						POLE_26,
						POLE_27,
						POLE_41,
						POLE_42,
						POLE_43,
						POLE_44,
						POLE_28
					FROM REJESTRY_UZYTKOWNIKA_DANE WHERE ID_REJESTRY_UZYTKOWNIKA_DANE=?`
	ImagesQ = `SELECT ID_OBRAZKA FROM OFERTY_WWW_OBRAZKI WHERE ID_OFERTY=?
           		AND ( ( STATUS<>'WWW_LOCKED' ) OR ( STATUS IS NULL ) ) ORDER BY KOLEJNOSC`
	ImageQ                     = `SELECT SQL_BIG_RESULT OBRAZEK FROM OBRAZKI WHERE ID_OBRAZKA=?`
	PersonQ                    = `SELECT IMIE, NAZWISKO, TELEFON, TELEFON_2, TELEFON_3, EMAIL FROM OSOBY WHERE ID_OSOBY=?`
	ID_WYROZNIENIA utils.ColId = 1
)
const (
	ID_OFERTY utils.ColId = iota
	ID_FIRMY
	NUMER_OFERTY
	AKTYWNA
	TYP_OFERTY
	RODZAJ_OFERTY
	ID_GRUPY
	NAZWA_OFERTY
	INFO
	TRESC_OFERTY
	OFERTA_ZAWIERA
	UWAGI
	CENA
	DATA_ZGLOSZENIA
	GODZINA_ZGLOSZENIA
	ID_WPROWADZAJACEGO
	DATA_MODYFIKACJI
	GODZINA_MODYFIKACJI
	ID_MODYFIKUJACEGO
	ILOSC_WYSWIETLEN
	CENA_NETTO
	STAWKA_VAT
	SWW
	JEDNOSTKA_MIARY
	WAGA
	ID_WLASCICIELA
	KOD_OFERTY
)

const (
	POLE_01 utils.ColId = iota
	POLE_02
	POLE_03
	POLE_04
	POLE_05
	POLE_06
	POLE_07
	POLE_08
	POLE_09
	POLE_10
	POLE_11
	POLE_60
	POLE_61
	POLE_62
	POLE_63
	POLE_64
	POLE_65
	POLE_21
	POLE_66
	POLE_67
	POLE_68
	POLE_69
	POLE_22
	POLE_23
	POLE_24
	POLE_25
	POLE_26
	POLE_27
	POLE_41
	POLE_42
	POLE_43
	POLE_44
	POLE_28
)

const (
	IMIE utils.ColId = iota
	NAZWISKO
	TELEFON
	TELEFON_2
	TELEFON_3
	EMAIL
)

type SpecialsMap map[string][]string

//
//	DB
//
type DB struct {
	Db *sql.DB
}

func (d *DB) Connection() *sql.DB {
	return d.Db
}
func (d *DB) Connect(user, pass, host, dbname string) (err error) {
	if d.Db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pass, host, dbname)); err != nil {
		return
	}
	return d.Db.Ping()
}

func (d *DB) Close() (err error) {
	err = d.Db.Close()
	utils.LogErr(err, "Close db")
	return
}

func (d *DB) Query(query string) (rows *sql.Rows, err error) {
	rows, err = d.Db.Query(query)
	utils.LogErr(err, "Query")
	return
}

func (d *DB) CloseRows(rows *sql.Rows) {
	utils.LogErr(rows.Close(), "Close rows")
}

type Images struct {
	*DB
	ids     []string
	offerId string
}

func NewImages(db *DB, offerId string) *Images {
	return &Images{db, []string{}, offerId}
}

func (i *Images) AddId(id string) {
	i.ids = append(i.ids, id)
}

func (i *Images) getIdsForOffer() (ids []string, err error) {
	var (
		stm  *sql.Stmt
		rows *sql.Rows
	)
	db := i.Connection()
	if stm, err = db.Prepare(ImagesQ); err != nil {
		return
	}
	if rows, err = stm.Query(i.offerId); err != nil {
		return
	}
	for rows.Next() {
		var imageId string
		if err = rows.Scan(&imageId); err != nil {
			return
		}
		i.AddId(imageId)
	}
	ids = i.ids
	return
}

func (i *Images) Ids() (ids []string, err error) {
	if len(i.ids) > 0 {
		ids = i.ids
		return
	}
	return i.getIdsForOffer()
}

func (i *Images) FileNames() (names []string, err error) {
	var (
		ids []string
	)
	ids, err = i.Ids()
	for _, id := range ids {
		names = append(names, fmt.Sprintf("%s.jpg", id))
	}
	return
}

func (i *Images) getRawImage(id string) (buff []byte, err error) {
	var (
		smt *sql.Stmt
	)
	db := i.Connection()
	if smt, err = db.Prepare(ImageQ); err != nil {
		return
	}
	row := smt.QueryRow(id)
	if err = row.Scan(&buff); err != nil {
		return
	}
	return
}

func (i *Images) saveImage(id, dest string) (err error) {
	var (
		buff []byte
	)
	if buff, err = i.getRawImage(id); err != nil {
		return
	}
	path := filepath.Join(dest, fmt.Sprintf("%s.jpg", id))
	return utils.SaveImageFromBytes(buff, path)
}

func (i *Images) SaveImages(dest string) (err error) {
	if _, err = i.Ids(); err != nil {
		return
	}
	for _, id := range i.ids {
		if err = i.saveImage(id, dest); err != nil {
			return
		}
	}
	return
}

//
// Offers
//

type Offers struct {
	*DB
	ids []int
}

func NewOffers(db *DB) *Offers {
	return &Offers{db, []int{}}
}

func (o *Offers) Ids() []int {
	return o.ids
}

func (o *Offers) AddId(id int) {
	o.ids = append(o.ids, id)
}

func (o *Offers) GetAllIds() (err error) {
	db := o.Connection()
	rows, err := db.Query(OffersQ)
	if err != nil {
		return
	}
	defer o.CloseRows(rows)
	for rows.Next() {
		var id int
		if err = rows.Scan(&id); err != nil {
			return
		}
		o.AddId(id)
	}
	err = rows.Err()
	return
}

//
// Offer
//
type Offer struct {
	*DB
	id  int
	row utils.Row
}

func NewOffer(db *DB, id int) *Offer {
	return &Offer{db, id, utils.Row{}}
}

func (o *Offer) Get() (err error) {
	db := o.Connection()
	row := db.QueryRow(OfferQ, o.id)
	dest := utils.Row{
		new(sql.NullInt64),  //ID_OFERTY,
		new(sql.NullInt64),  //ID_FIRMY,
		new(sql.NullString), //NUMER_OFERTY,
		new(sql.NullInt64),  //AKTYWNA
		new(sql.NullString), //TYP_OFERTY,
		new(sql.NullString), //RODZAJ_OFERTY,
		new(sql.NullString), //ID_GRUPY,
		new(sql.NullString), //NAZWA_OFERTY,
		new(sql.NullString), //INFO,
		new(sql.NullString), //TRESC_OFERTY,
		new(sql.NullString), //OFERTA_ZAWIERA,
		new(sql.NullString), //UWAGI,
		new(sql.NullInt64),  //CENA,
		new(sql.NullString), //DATA_ZGLOSZENIA,
		new(sql.NullString), //GODZINA_ZGLOSZENIA,
		new(sql.NullString), //ID_WPROWADZAJACEGO,
		new(sql.NullString), //DATA_MODYFIKACJI,
		new(sql.NullString), //GODZINA_MODYFIKACJI,
		new(sql.NullInt64),  //ID_MODYFIKUJACEGO,
		new(sql.NullInt64),  //ILOSC_WYSWIETLEN,
		new(sql.NullInt64),  //CENA_NETTO,
		new(sql.NullString), //STAWKA_VAT,
		new(sql.NullString), //SWW,
		new(sql.NullString), //JEDNOSTKA_MIARY,
		new(sql.NullInt64),  //WAGA,
		new(sql.NullString), //ID_WLASCICIELA,
		new(sql.NullString), //KOD_OFERTY
	}
	err = row.Scan(dest...)
	o.row = dest
	return
}

func (o *Offer) StrId() string {
	return fmt.Sprintf("%d", o.id)
}

func (o *Offer) Id() int {
	return o.id
}

func (o *Offer) StringAt(id utils.ColId) string {
	return utils.GetString(o.row, id)
}

func (o *Offer) IntAt(id utils.ColId) string {
	if val, ok := utils.GetInt64(o.row, id); ok {
		return fmt.Sprintf("%d", val)
	} else {
		return ""
	}
}

//
// specials
//
type Specials struct {
	*DB
	rows []utils.Row
}

func NewSpecials(db *DB) *Specials {
	return &Specials{db, []utils.Row{}}
}

func (s *Specials) GetAll() (err error) {
	db := s.Connection()
	rows, err := db.Query(SpecialsQ)
	if err != nil {
		return
	}
	defer s.CloseRows(rows)
	for rows.Next() {
		dest := utils.Row{
			new(sql.NullString), //ID_OFERTY,
			new(sql.NullString)} //ID_WYROZNIENIA,
		if err = rows.Scan(dest...); err != nil {
			return
		}
		s.Add(dest)
	}
	err = rows.Err()
	return
}

func (s *Specials) Add(row utils.Row) {
	s.rows = append(s.rows, row)
}

func (s *Specials) AsMap() SpecialsMap {
	m := make(SpecialsMap)
	for _, row := range s.rows {
		id := utils.GetString(row, ID_OFERTY)
		val := utils.GetString(row, ID_WYROZNIENIA)
		if arr, ok := m[id]; ok {
			arr = append(arr, val)
			m[id] = arr
		} else {
			m[id] = []string{val}
		}
	}
	return m
}

//
// Additions
//

type Additional struct {
	*DB
	Row utils.Row
}

func NewAdditional(db *DB) *Additional {
	return &Additional{db, utils.Row{}}
}

func (a *Additional) Get(offerId int) (err error) {
	db := a.Connection()
	row := db.QueryRow(AdditionslQ, offerId)

	dest := utils.Row{
		new(sql.NullString), //POLE_01
		new(sql.NullString), //POLE_02
		new(sql.NullString), //POLE_03
		new(sql.NullString), //POLE_04
		new(sql.NullString), //POLE_05
		new(sql.NullString), //POLE_06
		new(sql.NullString), //POLE_07
		new(sql.NullString), //POLE_08
		new(sql.NullString), //POLE_09
		new(sql.NullString), //POLE_10
		new(sql.NullString), //POLE_11
		new(sql.NullString), //POLE_60
		new(sql.NullString), //POLE_61
		new(sql.NullString), //POLE_62
		new(sql.NullString), //POLE_63
		new(sql.NullString), //POLE_64
		new(sql.NullString), //POLE_65
		new(sql.NullString), //POLE_21
		new(sql.NullString), //POLE_66
		new(sql.NullString), //POLE_67
		new(sql.NullString), //POLE_68
		new(sql.NullString), //POLE_69
		new(sql.NullString), //POLE_22
		new(sql.NullString), //POLE_23
		new(sql.NullString), //POLE_24
		new(sql.NullString), //POLE_25
		new(sql.NullString), //POLE_26
		new(sql.NullString), //POLE_27
		new(sql.NullString), //POLE_41
		new(sql.NullString), //POLE_42
		new(sql.NullString), //POLE_43
		new(sql.NullString), //POLE_44
		new(sql.NullString), //POLE_28
	}
	err = row.Scan(dest...)
	a.Row = dest
	return
}

func (a *Additional) StringAt(id utils.ColId) string {
	val := utils.GetString(a.Row, id)
	if val == "0" {
		return ""
	} else {
		return val
	}
}

func (a *Additional) IntAt(id utils.ColId) string {
	if val, ok := utils.GetInt64(a.Row, id); ok {
		return fmt.Sprintf("%d", val)
	} else {
		return ""
	}
}

//
// person
//

type Person struct {
	*DB
	row utils.Row
}

func NewPerson(db *DB) *Person {
	return &Person{db, utils.Row{}}
}

func (p *Person) Get(id string) (err error) {
	var stm *sql.Stmt
	db := p.Connection()
	if stm, err = db.Prepare(PersonQ); err != nil {
		return
	}
	row := stm.QueryRow(id)
	dest := utils.Row{
		new(sql.NullString), //IMIE
		new(sql.NullString), //NAZWISKO
		new(sql.NullString), //TELEFON
		new(sql.NullString), //TELEFON_2
		new(sql.NullString), //TELEFON_3
		new(sql.NullString), //EMAIL
	}
	err = row.Scan(dest...)
	p.row = dest
	return
}

func (p *Person) StringAt(id utils.ColId) string {
	return utils.GetString(p.row, id)
}

func (p *Person) DecodeStringAt(id utils.ColId) string {
	return utils.DecodeStr(utils.GetString(p.row, id))
}

func (p *Person) IntAt(id utils.ColId) string {
	if val, ok := utils.GetInt64(p.row, id); ok {
		return fmt.Sprintf("%d", val)
	} else {
		return ""
	}
}
