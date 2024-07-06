package importer

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/arednch/phonebook/data"
)

const (
	headerFirstName   = "first_name"
	headerLastName    = "name"
	headerCallsign    = "callsign"
	headerPhoneNumber = "telephone"
	headerPrivate     = "privat"
)

func ReadFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func ReadFromFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func ReadPhonebook(path string) ([]*data.Entry, error) {
	var blob []byte
	var err error
	switch {
	case strings.HasPrefix(path, "http://"):
		fallthrough
	case strings.HasPrefix(path, "https://"):
		blob, err = ReadFromURL(path)
	case strings.HasPrefix(path, "/"):
		blob, err = ReadFromFile(path)
	default:
		err = errors.New("unknown or unsupported path scheme (needs to be a valid, absolute file path or http/https URL)")
	}
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(bytes.NewBuffer(blob))
	// read and index headers
	hdrs, err := reader.Read()
	if err != nil {
		return nil, err
	}
	headers := make(map[string]int)
	for i, v := range hdrs {
		headers[strings.ToLower(v)] = i
	}
	firstIdx, ok := headers[headerFirstName]
	if !ok {
		return nil, errors.New("unable to locate first name column in CSV")
	}
	lastIdx, ok := headers[headerLastName]
	if !ok {
		return nil, errors.New("unable to locate last name column in CSV")
	}
	callIdx, ok := headers[headerCallsign]
	if !ok {
		return nil, errors.New("unable to locate callsign column in CSV")
	}
	phoneIdx, ok := headers[headerPhoneNumber]
	if !ok {
		return nil, errors.New("unable to locate phone number column in CSV")
	}
	privateIdx, privateIdxAvailable := headers[headerPrivate]

	var records []*data.Entry
	for {
		r, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// skip if we encounter the first empty line
		if strings.TrimSpace(r[firstIdx]) == "" && strings.TrimSpace(r[lastIdx]) == "" &&
			strings.TrimSpace(r[callIdx]) == "" && strings.TrimSpace(r[phoneIdx]) == "" {
			break
		}
		// check if entry is marked as private and if so, skip it
		if privateIdxAvailable && strings.ToLower(strings.TrimSpace(r[privateIdx])) == "y" {
			continue
		}

		entry := &data.Entry{
			FirstName:   strings.TrimSpace(r[firstIdx]),
			LastName:    strings.TrimSpace(r[lastIdx]),
			Callsign:    strings.TrimSpace(r[callIdx]),
			PhoneNumber: strings.TrimSpace(r[phoneIdx]),
		}
		records = append(records, entry)
	}

	return records, nil
}
