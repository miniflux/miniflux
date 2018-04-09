// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package storage

import (
	"fmt"
	"strconv"

	"github.com/miniflux/miniflux/logger"
	"github.com/miniflux/miniflux/sql"
)

const schemaVersion = 16

// Migrate run database migrations.
func (s *Storage) Migrate() {
	var currentVersion int
	s.db.QueryRow(`select version from schema_version`).Scan(&currentVersion)

	fmt.Println("Current schema version:", currentVersion)
	fmt.Println("Latest schema version:", schemaVersion)

	for version := currentVersion + 1; version <= schemaVersion; version++ {
		fmt.Println("Migrating to version:", version)

		tx, err := s.db.Begin()
		if err != nil {
			logger.Fatal("[Storage:Migrate] %v", err)
		}

		rawSQL := sql.SqlMap["schema_version_"+strconv.Itoa(version)]
		// fmt.Println(rawSQL)
		_, err = tx.Exec(rawSQL)
		if err != nil {
			tx.Rollback()
			logger.Fatal("[Storage:Migrate] %v", err)
		}

		if _, err := tx.Exec(`delete from schema_version`); err != nil {
			tx.Rollback()
			logger.Fatal("[Storage:Migrate] %v", err)
		}

		if _, err := tx.Exec(`insert into schema_version (version) values($1)`, version); err != nil {
			tx.Rollback()
			logger.Fatal("[Storage:Migrate] %v", err)
		}

		if err := tx.Commit(); err != nil {
			logger.Fatal("[Storage:Migrate] %v", err)
		}
	}
}
