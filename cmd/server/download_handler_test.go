// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/moov-io/base"
	"github.com/moov-io/ofac/internal/database"

	"github.com/go-kit/kit/log"
)

func TestDownload__manualRefreshPath(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		return
	}

	check := func(t *testing.T, repo *sqliteDownloadRepository) {
		searcher := &searcher{}

		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/data/refresh", nil)
		req.Header.Set("x-request-id", base.ID())

		manualRefreshHandler(log.NewNopLogger(), searcher, repo)(w, req)
		w.Flush()

		if w.Code != http.StatusOK {
			t.Errorf("bogus status code: %d", w.Code)
		}
		var stats downloadStats
		if err := json.NewDecoder(w.Body).Decode(&stats); err != nil {
			t.Error(err)
		}
		if stats.SDNs == 0 {
			t.Errorf("stats.SDNs=%d but expected non-zero", stats.SDNs)
		}
	}

	// SQLite tests
	sqliteDB := database.CreateTestSqliteDB(t)
	defer sqliteDB.Close()
	check(t, &sqliteDownloadRepository{sqliteDB.DB, log.NewNopLogger()})

	// MySQL tests
	mysqlDB := database.CreateTestMySQLDB(t)
	defer mysqlDB.Close()
	check(t, &sqliteDownloadRepository{mysqlDB.DB, log.NewNopLogger()})
}
