package exporter

import (
	"encoding/xml"

	"github.com/finack/phonebook/data"
)

type Yealink struct{}

func (y *Yealink) Export(entries []*data.Entry, pbx bool) ([]byte, error) {
	return xml.MarshalIndent(struct {
		*GenericPhoneBook
		XMLName struct{} `xml:"YealinkIPPhoneDirectory"`
	}{
		GenericPhoneBook: export(entries, pbx),
	}, "", "    ")
}