package db

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/kpawlik/exportms/utils"
)

const (
	// OffersQ query to get offers id list
	OffersQ = `SELECT ID_OFERTY FROM OFERTY_WWW WHERE AKTYWNA='1' AND ID_OFERTY NOT IN
	(SELECT DISTINCT ID_OFERTY FROM OFERTY_WWW_WYROZNIENIA WHERE ID_WYROZNIENIA = 'NIE_EKSPORTUJ') `
	// SpecialsQ gets all specials offers ids
	SpecialsQ = `SELECT ID_OFERTY,
					ID_WYROZNIENIA FROM OFERTY_WWW_WYROZNIENIA`
	// OfferQ gets offers details
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
	// AdditionslQ additionals data
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
	// ImagesQ query to gets pictures for order
	ImagesQ = `SELECT ID_OBRAZKA FROM OFERTY_WWW_OBRAZKI WHERE ID_OFERTY=?
           		AND ( ( STATUS<>'WWW_LOCKED' ) OR ( STATUS IS NULL ) ) ORDER BY KOLEJNOSC`
	// ImageQ gets picture data
	ImageQ = `SELECT SQL_BIG_RESULT OBRAZEK FROM OBRAZKI WHERE ID_OBRAZKA=?`
	// PersonQ gets person data
	PersonQ = `SELECT IMIE, NAZWISKO, TELEFON, TELEFON_2, TELEFON_3, EMAIL FROM OSOBY WHERE ID_OSOBY=?`
	// IDWyroznienia speciall id
	IDWyroznienia utils.ColID = 1
)

// database ids
const (
	IDOferty utils.ColID = iota
	IDFirmy
	NumerOferty
	Aktywna
	TypOferty
	RodzajOferty
	IDGrupy
	NazwaOferty
	Info
	TrescOferty
	OfertaZamowienia
	Uwagi
	Cena
	DataZgloszenia
	GodzinaZgloszenia
	IDWprowadzajacego
	DataModyfikacji
	GodzinaModyfikacji
	IDModyfikujacego
	IloscWyswietlen
	CenaNetto
	StawkaVat
	Sww
	JednostaMiary
	Waga
	IDWlasciciela
	KodOferty
)

// Additional fields
const (
	Pole01 utils.ColID = iota
	Pole02
	Pole03
	Pole04
	Pole05
	Pole06
	Pole07
	Pole08
	Pole09
	Pole10
	Pole11
	Pole60
	Pole61
	Pole62
	Pole63
	Pole64
	Pole65
	Pole21
	Pole66
	Pole67
	Pole68
	Pole69
	Pole22
	Pole23
	Pole24
	Pole25
	Pole26
	Pole27
	Pole41
	Pole42
	Pole43
	Pole44
	Pole28
)

// User details ids
const (
	Imie utils.ColID = iota
	Nazwisko
	Telefon
	Telefon2
	Telefon3
	Email
)

// SpecialsMap type to store map of specials
type SpecialsMap map[string][]string

// DB struct and methods
type DB struct {
	Db *sql.DB
}

// Connection to databases
func (d *DB) Connection() *sql.DB {
	return d.Db
}

// Connect to database
func (d *DB) Connect(user, pass, host, dbname string) (err error) {
	if d.Db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pass, host, dbname)); err != nil {
		return
	}
	return d.Db.Ping()
}

// Close close db connections
func (d *DB) Close() (err error) {
	err = d.Db.Close()
	utils.LogErr(err, "Close db")
	return
}

// Query returns query results
func (d *DB) Query(query string) (rows *sql.Rows, err error) {
	rows, err = d.Db.Query(query)
	utils.LogErr(err, "Query")
	return
}

// CloseRows close all db rows
func (d *DB) CloseRows(rows *sql.Rows) {
	utils.LogErr(rows.Close(), "Close rows")
}

//
// Images struct with methods
//
type Images struct {
	*DB
	ids     []string
	offerID string
}

// NewImages new struct with db context
func NewImages(db *DB, offerID string) *Images {
	return &Images{db, []string{}, offerID}
}

// AddID add image id
func (i *Images) AddID(id string) {
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
	if rows, err = stm.Query(i.offerID); err != nil {
		return
	}
	for rows.Next() {
		var imageID string
		if err = rows.Scan(&imageID); err != nil {
			return
		}
		i.AddID(imageID)
	}
	ids = i.ids
	return
}

// Ids return list of all images ids
func (i *Images) Ids() (ids []string, err error) {
	if len(i.ids) > 0 {
		ids = i.ids
		return
	}
	return i.getIdsForOffer()
}

// FileNames return list of file names
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

//
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

// SaveImages save images
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
// Offers struct and methods
//
type Offers struct {
	*DB
	ids []int
}

// NewOffers new structure with db context
func NewOffers(db *DB) *Offers {
	return &Offers{db, []int{}}
}

//Ids returns ids
func (o *Offers) Ids() []int {
	return o.ids
}

// AddID adds id to offer
func (o *Offers) AddID(id int) {
	o.ids = append(o.ids, id)
}

// GetAllIds gets all id for offer
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
		o.AddID(id)
	}
	err = rows.Err()
	return
}

//
// Offer struct and methods
//
type Offer struct {
	*DB
	id  int
	row utils.Row
}

// NewOffer new struct with db context
func NewOffer(db *DB, id int) *Offer {
	return &Offer{db, id, utils.Row{}}
}

// Get record from db
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

// StrID returns id as string
func (o *Offer) StrID() string {
	return fmt.Sprintf("%d", o.id)
}

// ID returns id value
func (o *Offer) ID() int {
	return o.id
}

// StringAt column value with id index
func (o *Offer) StringAt(id utils.ColID) string {
	return utils.GetString(o.row, id)
}

// IntAt column value with id index
func (o *Offer) IntAt(id utils.ColID) string {
	if val, ok := utils.GetInt64(o.row, id); ok {
		return fmt.Sprintf("%d", val)
	}
	return ""
}

//
// Specials struct and methods
//
type Specials struct {
	*DB
	rows []utils.Row
}

// NewSpecials new struct with db context
func NewSpecials(db *DB) *Specials {
	return &Specials{db, []utils.Row{}}
}

// GetAll gets all specials from db
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

// Add row to specials s
func (s *Specials) Add(row utils.Row) {
	s.rows = append(s.rows, row)
}

// AsMap return specials as map [id] strings{}
func (s *Specials) AsMap() SpecialsMap {
	m := make(SpecialsMap)
	for _, row := range s.rows {
		id := utils.GetString(row, IDOferty)
		val := utils.GetString(row, IDWyroznienia)
		if arr, ok := m[id]; ok {
			arr = append(arr, val)
			m[id] = arr
		} else {
			m[id] = []string{val}
		}
	}
	return m
}

// Additional struct and methods
type Additional struct {
	*DB
	Row utils.Row
}

//NewAdditional create new Additional struct with db context
func NewAdditional(db *DB) *Additional {
	return &Additional{db, utils.Row{}}
}

// Get gets record from db
func (a *Additional) Get(offerID int) (err error) {
	db := a.Connection()
	row := db.QueryRow(AdditionslQ, offerID)

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

// StringAt column value with id index
func (a *Additional) StringAt(id utils.ColID) string {
	val := utils.GetString(a.Row, id)
	if val == "0" {
		return ""
	}
	return val
}

// IntAt column value with id index
func (a *Additional) IntAt(id utils.ColID) string {
	if val, ok := utils.GetInt64(a.Row, id); ok {
		return fmt.Sprintf("%d", val)
	}
	return ""
}

//
// Person struct and methods
//
type Person struct {
	*DB
	row utils.Row
}

// NewPerson creates new structure with database context
func NewPerson(db *DB) *Person {
	return &Person{db, utils.Row{}}
}

// Get return person record with id
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

// StringAt return row string from p.row with id index
func (p *Person) StringAt(id utils.ColID) string {
	return utils.GetString(p.row, id)
}

// DecodeStringAt return decoded string from column with id index
func (p *Person) DecodeStringAt(id utils.ColID) string {
	return utils.DecodeStr(utils.GetString(p.row, id))
}

// IntAt return column value from p.row with id index
func (p *Person) IntAt(id utils.ColID) string {
	if val, ok := utils.GetInt64(p.row, id); ok {
		return fmt.Sprintf("%d", val)
	}
	return ""
}
