// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/moov-io/ofac"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
)

func TestSearch__Address(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/search?address=ibex+house+minories&limit=1", nil)

	router := mux.NewRouter()
	addSearchRoutes(log.NewNopLogger(), router, addressSearcher)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("bogus status code: %d", w.Code)
	}

	if v := w.Body.String(); !strings.Contains(v, `"match":1`) {
		t.Errorf("%#v", v)
	}

	var wrapper struct {
		Addresses []*ofac.Address `json:"addresses"`
	}
	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}
	if wrapper.Addresses[0].EntityID != "173" {
		t.Errorf("%#v", wrapper.Addresses[0])
	}

	// send an empty body and get an error
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/search?limit=1", nil)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusBadRequest {
		t.Errorf("bogus status code: %d", w.Code)
	}
}

func TestSearch__AddressCountry(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/search?country=united+kingdom&limit=1", nil)

	router := mux.NewRouter()
	addSearchRoutes(log.NewNopLogger(), router, addressSearcher)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("bogus status code: %d", w.Code)
	}

	if v := w.Body.String(); !strings.Contains(v, `"match":1`) {
		t.Errorf("%#v", v)
	}
}

func TestSearch__AddressMulti(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/search?address=ibex+house&country=united+kingdom&limit=1", nil)

	router := mux.NewRouter()
	addSearchRoutes(log.NewNopLogger(), router, addressSearcher)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("bogus status code: %d", w.Code)
	}

	if v := w.Body.String(); !strings.Contains(v, `"match":0.921`) {
		t.Errorf("%#v", v)
	}
}

func TestSearch__AddressProvidence(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/search?address=ibex+house&country=united+kingdom&providence=london+ec3n+1DY&limit=1", nil)

	router := mux.NewRouter()
	addSearchRoutes(log.NewNopLogger(), router, addressSearcher)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("bogus status code: %d", w.Code)
	}

	if v := w.Body.String(); !strings.Contains(v, `"match":0.947`) {
		t.Errorf("%#v", v)
	}
}

func TestSearch__AddressCity(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/search?address=ibex+house&country=united+kingdom&city=london+ec3n+1DY&limit=1", nil)

	router := mux.NewRouter()
	addSearchRoutes(log.NewNopLogger(), router, addressSearcher)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("bogus status code: %d", w.Code)
	}

	if v := w.Body.String(); !strings.Contains(v, `"match":0.947`) {
		t.Errorf("%#v", v)
	}
}

func TestSearch__AddressState(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/search?address=ibex+house&country=united+kingdom&state=london+ec3n+1DY&limit=1", nil)

	router := mux.NewRouter()
	addSearchRoutes(log.NewNopLogger(), router, addressSearcher)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("bogus status code: %d", w.Code)
	}

	if v := w.Body.String(); !strings.Contains(v, `"match":0.947`) {
		t.Errorf("%#v", v)
	}
}

func TestSearch__NameAndAddress(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/search?name=midco&address=rue+de+rhone&limit=1", nil)

	s := &searcher{
		Addresses: precomputeAddresses([]*ofac.Address{
			{
				EntityID:                    "2831",
				AddressID:                   "1965",
				Address:                     "57 Rue du Rhone",
				CityStateProvincePostalCode: "Geneva CH-1204",
				Country:                     "Switzerland",
			},
			{
				EntityID:                    "173",
				AddressID:                   "129",
				Address:                     "Ibex House, The Minories",
				CityStateProvincePostalCode: "London EC3N 1DY",
				Country:                     "United Kingdom",
			},
		}),
		SDNs: precomputeSDNs([]*ofac.SDN{
			{
				EntityID: "2831",
				SDNName:  "MIDCO FINANCE S.A.",
				SDNType:  "individual",
				Program:  "IRAQ2",
				Remarks:  "US FEIN CH-660-0-469-982-0 (United States); Switzerland.",
			},
		}),
	}

	router := mux.NewRouter()
	addSearchRoutes(log.NewNopLogger(), router, s)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("bogus status code: %d", w.Code)
	}
	var wrapper struct {
		SDNs      []*ofac.SDN     `json:"SDNs"`
		Addresses []*ofac.Address `json:"addresses"`
	}
	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}

	if len(wrapper.SDNs) != 1 || len(wrapper.Addresses) != 1 {
		t.Fatalf("sdns=%#v addresses=%#v", wrapper.SDNs[0], wrapper.Addresses[0])
	}

	if wrapper.SDNs[0].EntityID != "2831" || wrapper.Addresses[0].EntityID != "2831" {
		t.Errorf("SDNs[0].EntityID=%s Addresses[0].EntityID=%s", wrapper.SDNs[0].EntityID, wrapper.Addresses[0].EntityID)
	}

	// request with no results
	req = httptest.NewRequest("GET", "/search?name=midco&country=United+Kingdom&limit=1", nil)

	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("bogus status code: %d", w.Code)
	}
	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}
	if len(wrapper.SDNs) != 0 || len(wrapper.Addresses) != 0 {
		t.Fatalf("sdns=%#v addresses=%#v", wrapper.SDNs[0], wrapper.Addresses[0])
	}
}

func TestSearch__NameAndAltName(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/search?limit=1&q=nayif", nil)

	s := &searcher{
		Alts:      altSearcher.Alts,
		SDNs:      sdnSearcher.SDNs,
		Addresses: addressSearcher.Addresses,
		DPs:       dplSearcher.DPs,
	}

	router := mux.NewRouter()
	addSearchRoutes(log.NewNopLogger(), router, s)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("bogus status code: %d", w.Code)
	}

	// read response body
	var wrapper struct {
		SDNs          []*ofac.SDN               `json:"SDNs"`
		AltNames      []*ofac.AlternateIdentity `json:"altNames"`
		Addresses     []*ofac.Address           `json:"addresses"`
		DeniedPersons []*ofac.DPL               `json:"deniedPersons"`
	}
	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}
	if wrapper.SDNs[0].EntityID != "2681" {
		t.Errorf("%#v", wrapper.SDNs[0])
	}
	if wrapper.AltNames[0].EntityID != "4691" {
		t.Errorf("%#v", wrapper.AltNames[0].EntityID)
	}
	if wrapper.Addresses[0].EntityID != "735" {
		t.Errorf("%#v", wrapper.Addresses[0].EntityID)
	}
	if wrapper.DeniedPersons[0].StreetAddress != "P.O. BOX 28360" {
		t.Errorf("%#v", wrapper.DeniedPersons[0].StreetAddress)
	}
}

func TestSearch__Name(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/search?name=Dr+AL+ZAWAHIRI&limit=1", nil)

	router := mux.NewRouter()
	addSearchRoutes(log.NewNopLogger(), router, sdnSearcher)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("bogus status code: %d", w.Code)
	}

	if v := w.Body.String(); !strings.Contains(v, `"match":1`) {
		t.Error(v)
	}

	var wrapper struct {
		SDNs []*ofac.SDN `json:"SDNs"`
	}
	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}
	if wrapper.SDNs[0].EntityID != "2676" {
		t.Errorf("%#v", wrapper.SDNs[0])
	}

	// directly check searcher.TopSDNs
	sdns := sdnSearcher.TopSDNs(1, "Dr AL ZAWAHIRI")
	if len(sdns) != 1 {
		t.Errorf("got SDNs: %#v", sdns)
	}
	if sdns[0].EntityID != "2676" {
		t.Errorf("%#v", sdns[0])
	}
	eql(t, "name match", sdns[0].match, 1.0)
}

func TestSearch__AltName(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/search?altName=sogo+KENKYUSHO&limit=1", nil)

	router := mux.NewRouter()
	addSearchRoutes(log.NewNopLogger(), router, altSearcher)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("bogus status code: %d", w.Code)
	}

	if v := w.Body.String(); !strings.Contains(v, `"match":1`) {
		t.Error(v)
	}

	var wrapper struct {
		Alts []*ofac.AlternateIdentity `json:"altNames"`
	}
	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}
	if wrapper.Alts[0].EntityID != "4691" {
		t.Errorf("%#v", wrapper.Alts[0])
	}
}

func TestSearch__ID(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/search?id=5892464&limit=2", nil)

	router := mux.NewRouter()
	addSearchRoutes(log.NewNopLogger(), router, idSearcher)
	router.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusOK {
		t.Errorf("bogus status code: %d", w.Code)
	}

	if v := w.Body.String(); !strings.Contains(v, `"match":1`) {
		t.Error(v)
	}

	var wrapper struct {
		SDNs []*ofac.SDN `json:"SDNs"`
	}
	if err := json.NewDecoder(w.Body).Decode(&wrapper); err != nil {
		t.Fatal(err)
	}
	if wrapper.SDNs[0].EntityID != "22790" {
		t.Errorf("%#v", wrapper.SDNs[0])
	}
}
