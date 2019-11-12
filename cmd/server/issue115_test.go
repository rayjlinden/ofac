// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"github.com/moov-io/ofac"

	"github.com/go-kit/kit/log"
)

func TestIssue115__TopSDNs(t *testing.T) {
	score := jaroWinkler("georgehabbash", "georgebush")
	eql(t, "george bush jaroWinkler", score, 0.896)

	score = jaroWinkler("g", "geoergebush")
	eql(t, "g vs geoergebush", score, 0.697)

	s := &searcher{logger: log.NewNopLogger()}

	// Issue 115 (https://github.com/moov-io/ofac/issues/115) talks about how "george bush" is a false positive (90%) match against
	// several other "George ..." records. This is too sensitive and so we need to tone that down.

	// was 89.6% match
	s.SDNs = precomputeSDNs([]*ofac.SDN{{EntityID: "2680", SDNName: "HABBASH, George", SDNType: "INDIVIDUAL"}})
	out := s.TopSDNs(1, "george bush")
	eql(t, "issue115: top SDN 2680", out[0].match, 0.500)

	// was 88.3% match
	s.SDNs = precomputeSDNs([]*ofac.SDN{{EntityID: "9432", SDNName: "CHIWESHE, George", SDNType: "INDIVIDUAL"}})
	out = s.TopSDNs(1, "george bush")
	eql(t, "issue115: top SDN 18996", out[0].match, 0.729)

	// another example
	s.SDNs = precomputeSDNs([]*ofac.SDN{{EntityID: "0", SDNName: "Bush, George W", SDNType: "INDIVIDUAL"}})
	if s.SDNs[0].name != "george w bush" {
		t.Errorf("s.SDNs[0].name=%s", s.SDNs[0].name)
	}

	out = s.TopSDNs(1, "george bush")
	eql(t, "issue115: top SDN 0", out[0].match, 1.0)

	out = s.TopSDNs(1, "george w bush")
	eql(t, "issue115: top SDN 0", out[0].match, 1.0)
}
