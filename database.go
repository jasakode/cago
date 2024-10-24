// Copyright (c) 2024, Jasakode Authors.
// All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package cago

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type database struct {
	mu        sync.Mutex
	sqldb     *sql.DB
	tableName string
}

type model struct {
	ID    uint64 `json:"id"`
	Key   string `json:"key"`
	Value []byte `json:"value"`
}

func (app *App) InitializeDB() error {
	db := database{}
	db.tableName = "cagos"

	d, err := sql.Open("sqlite3", app.config.Path)
	if err != nil {
		return err
	}

	app.mu.Lock()
	defer app.mu.Unlock()
	db.sqldb = d
	app.db = &db
	return nil
}

// create table if not exist
func (db *database) CreateTableIfNotExist() error {
	createTableQuery := `
  	CREATE TABLE IF NOT EXISTS %s (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        key TEXT NOT NULL,
        value BLOB
    );`
	db.mu.Lock()
	defer db.mu.Unlock()
	_, err := db.sqldb.Exec(fmt.Sprintf(createTableQuery, db.tableName))
	if err != nil {
		return err
	}
	return nil
}

func (db *database) InsertOrUpdate(key string, data []byte) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Menggunakan ON CONFLICT untuk update jika key sudah ada
	insertOrUpdateQuery := `
		INSERT INTO %s (key, value) 
		VALUES (?, ?)
		ON CONFLICT(key) 
		DO UPDATE SET value = excluded.value;
	`
	_, err := db.sqldb.Exec(fmt.Sprintf(insertOrUpdateQuery, db.tableName), key, data)
	if err != nil {
		return err
	}
	return nil
}

func (db *database) FindALL() (*[]model, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Menyiapkan query untuk mengambil semua data
	selectQuery := `SELECT id, key, value FROM %s;`
	// Menjalankan query dan mendapatkan hasil
	rows, err := db.sqldb.Query(fmt.Sprintf(selectQuery, db.tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []model{}
	for rows.Next() {
		r := model{}
		err := rows.Scan(&r.ID, &r.Key, &r.Value)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return &result, nil
}

// RemoveByKey menghapus entri dari database berdasarkan kunci yang diberikan.
// Fungsi ini mengunci database untuk memastikan tidak ada akses bersamaan
// saat melakukan penghapusan. Jika terjadi kesalahan saat mengeksekusi
// perintah SQL, kesalahan tersebut akan dikembalikan.
//
// Parameter:
//   - key: Kunci dari entri yang ingin dihapus.
//
// Mengembalikan:
//   - error: Kesalahan jika terjadi selama proses penghapusan.
func (db *database) RemoveByKey(key string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Menyiapkan query untuk menghapus entri berdasarkan kunci
	removeQuery := `
		DELETE FROM %s 
		WHERE key = ?;
	`
	_, err := db.sqldb.Exec(fmt.Sprintf(removeQuery, db.tableName), key)
	if err != nil {
		return err
	}
	return nil
}

// RemoveAll menghapus semua entri dari tabel dalam database.
// Fungsi ini mengunci database untuk memastikan tidak ada akses bersamaan
// saat melakukan penghapusan. Jika terjadi kesalahan saat mengeksekusi
// perintah SQL, kesalahan tersebut akan dikembalikan.
//
// Mengembalikan:
//   - error: Kesalahan jika terjadi selama proses penghapusan.
func (db *database) RemoveAll() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Menyiapkan query untuk menghapus semua entri dari tabel
	removeAllQuery := `
		DELETE FROM %s;
	`
	_, err := db.sqldb.Exec(fmt.Sprintf(removeAllQuery, db.tableName))
	if err != nil {
		return err
	}
	return nil
}
