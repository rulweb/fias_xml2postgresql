package current_status

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"

	"github.com/go-gorp/gorp"
	_ "github.com/lib/pq"
)

const dateformat = "2006-01-02"

// Статус актуальности КЛАДР 4.0
type XmlObject struct {
	XMLName    xml.Name `xml:"CurrentStatus"`
	CURENTSTID int      `xml:"CURENTSTID,attr"`
	NAME       string   `xml:"NAME,attr"`
}

type DBObject struct {
	CURENTSTID int    `db:"currentst_id, primarykey"`
	NAME       string `db:"name"`
}

func xml2db(xml XmlObject) *DBObject {
	obj := &DBObject{
		CURENTSTID: xml.CURENTSTID,
		NAME:       xml.NAME}

	return obj
}

func Export(dbmap *gorp.DbMap) {
	// Создаем таблицу
	dbmap.AddTableWithName(DBObject{}, "curentst")
	err := dbmap.DropTableIfExists(DBObject{})
	if err != nil {
		fmt.Println("Error on drop table:", err)
		return
	}
	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		fmt.Println("Error on creating table:", err)
		return
	}

	xmlFile, err := os.Open("xml/AS_CURENTST_20150705_a4c01a56-bea6-4cf3-84f8-10d201df820d.XML")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	total := 0
	var inElement string
	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		// Inspect the type of the token just read.
		switch se := t.(type) {
		case xml.StartElement:
			// If we just read a StartElement token
			inElement = se.Name.Local

			if inElement == "CurrentStatus" {
				total++
				var item XmlObject

				// decode a whole chunk of following XML into the
				// variable item which is a ActualStatus (se above)
				decoder.DecodeElement(&item, &se)
				obj := xml2db(item)
				err := dbmap.Insert(obj)
				if err != nil {
					fmt.Println("Error on creating table:", err)
					return
				}

				s := strconv.Itoa(total)
				fmt.Printf("\rCurrentStatus: %s rows", s)
			}
		default:
		}

	}

	fmt.Printf("Total processed items in CurrentStatus: %d \n", total)
}