// Copyright 2020 The Logsuck Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package events

import (
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/jackbister/logsuck/internal/config"

	_ "github.com/mattn/go-sqlite3"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz ABCDEFGHIJKLMNOPQRSTUVWXYZ .~!@$%^&*()_+=")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomStr(length int) string {
	result := make([]rune, length)
	for i := range result {
		result[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(result)
}

func TestAddBatchTrueBatch(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("TestRexPipelineStep got error when creating in-memory SQLite database: %v", err)
	}
	repo, err := SqliteRepository(db, &config.SqliteConfig{
		DatabaseFile: ":memory:",
		TrueBatch:    true,
	})
	if err != nil {
		t.Fatalf("TestRexPipelineStep got error when creating events repo: %v", err)
	}

	repo.AddBatch([]Event{
		{
			Raw:       "2021-02-01 00:00:00 log event",
			Timestamp: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
			Host:      "localhost",
			Source:    "log.txt",
			Offset:    0,
		},
	})

	evts, err := repo.GetByIds([]int64{0}, SortModeNone)
	if err != nil {
		t.Fatalf("got error when retrieving event: %v", err)
	}
	if len(evts) != 1 {
		t.Fatalf("got unexpected number of events, expected 1 event but got %v", len(evts))
	}
}

func TestAddBatchOneByOne(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("TestRexPipelineStep got error when creating in-memory SQLite database: %v", err)
	}
	repo, err := SqliteRepository(db, &config.SqliteConfig{
		DatabaseFile: ":memory:",
		TrueBatch:    false,
	})
	if err != nil {
		t.Fatalf("TestRexPipelineStep got error when creating events repo: %v", err)
	}

	repo.AddBatch([]Event{
		{
			Raw:       "2021-02-01 00:00:00 log event",
			Timestamp: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
			Host:      "localhost",
			Source:    "log.txt",
			Offset:    0,
		},
	})

	evts, err := repo.GetByIds([]int64{0}, SortModeNone)
	if err != nil {
		t.Fatalf("got error when retrieving event: %v", err)
	}
	if len(evts) != 1 {
		t.Fatalf("got unexpected number of events, expected 1 event but got %v", len(evts))
	}
}

// TestCompression tests the effect of EnableFTSCompression on db size
func TestCompression(t *testing.T) {
	_ = os.Remove("db1.db")
	_ = os.Remove("db2.db")
	cfg1 := config.SqliteConfig{
		DatabaseFile:         "db1.db",
		TrueBatch:            false,
		EnableFTSCompression: false,
	}
	cfg2 := config.SqliteConfig{
		DatabaseFile:         "db2.db",
		TrueBatch:            false,
		EnableFTSCompression: true,
	}
	var db1, db2 *sql.DB
	var err error
	sql.Register("sqlite3_with_compression", &sqlite3.SQLiteDriver{
		Extensions: []string{"./../../compress"},
	})
	db1, _ = sql.Open("sqlite3", cfg1.DatabaseFile)
	db2, err = sql.Open("sqlite3_with_compression", cfg2.DatabaseFile)
	if err != nil {
		t.Fatal(err)
	}
	repo1, err := SqliteRepository(db1, &cfg1)
	if err != nil {
		t.Fatal(err)
	}
	repo2, err := SqliteRepository(db2, &cfg2)
	if err != nil {
		t.Fatal(err)
	}
	// Generate events
	eLen := 10_000
	events := make([]Event, eLen)
	for i := 0; i < eLen; i++ {
		events[i] = Event{
			Raw:       randomStr(i),
			Timestamp: time.Date(2021, 1, rand.Intn(30), rand.Intn(12), rand.Intn(60), 0, 0, time.UTC),
			Host:      fmt.Sprintf("localhost%v", i),
			Source:    "log.txt",
			Offset:    0,
		}
	}
	// Add some events
	_ = repo1.AddBatch(events)
	_ = repo2.AddBatch(events)

	// Fetch events. Compare random records. Assert if compression/decompression works fine!
	e1, err := repo1.GetByIds([]int64{1, 10, 100, 1000}, 1)
	if err != nil {
		t.Fatal(err)
	}
	e2, err := repo1.GetByIds([]int64{1, 10, 100, 1000}, 1)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(e1); i++ {
		if e1[i].Raw != e2[i].Raw {
			t.Errorf("Events MUST be equal!")
		}
	}

	// Compare db size
	f1, _ := os.Open("db1.db")
	f2, _ := os.Open("db2.db")
	info1, _ := f1.Stat()
	info2, _ := f2.Stat()
	fmt.Printf("Without compression: %v\nWith compression: %v\n", info1.Size(), info2.Size())
}
