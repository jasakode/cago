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

// Struktur `database` merepresentasikan koneksi database dengan fitur penguncian (mutex)
// untuk memastikan akses thread-safe ke database.
//
// Field-field:
//   - mu: Mutex yang digunakan untuk mengamankan akses ke database agar thread-safe.
//   - sqldb: Pointer ke objek sql.DB yang merepresentasikan koneksi database SQLite.
//   - tableName: Nama tabel yang digunakan dalam operasi database.
type database struct {
	mu        sync.Mutex // Mutex untuk menghindari race condition.
	sqldb     *sql.DB    // Koneksi ke database SQLite.
	tableName string     // Nama tabel yang digunakan dalam query.
}

// Struktur `model` merepresentasikan entitas data yang disimpan dalam tabel database.
// Struktur ini menyimpan id unik, kunci, dan nilai dalam bentuk byte array.
//
// Field-field:
//   - ID: ID unik dari setiap entri dalam tabel, yang di-auto-increment oleh database.
//   - Key: Kunci (key) untuk setiap entri yang bertipe string.
//   - Value: Nilai (value) yang disimpan dalam bentuk byte array.
type model struct {
	ID    uint64 `json:"id"`    // ID unik dari setiap entri, di-generate otomatis.
	Key   string `json:"key"`   // Kunci untuk mengidentifikasi entri.
	Value []byte `json:"value"` // Nilai data yang disimpan dalam format byte.
}

// InitializeDB menginisialisasi koneksi database SQLite dan menyimpannya dalam aplikasi.
// Fungsi ini menetapkan nama tabel yang digunakan, membuka koneksi ke database,
// dan menyimpan objek database ke dalam field aplikasi.
//
// Langkah-langkah:
//  1. Membuat objek database baru dengan nama tabel yang ditentukan.
//  2. Membuka koneksi ke SQLite menggunakan jalur database dari konfigurasi aplikasi.
//  3. Menyimpan koneksi database ke dalam aplikasi dengan penguncian untuk memastikan thread safety.
//
// Mengembalikan:
//   - error: Kesalahan jika koneksi database gagal dibuka.
func (app *App) InitializeDB() error {
	// Membuat instance baru dari struct database dan menetapkan nama tabel.
	db := database{}
	db.tableName = "cagos"

	// Membuka koneksi ke SQLite menggunakan path yang disimpan dalam konfigurasi aplikasi.
	d, err := sql.Open("sqlite3", app.config.Path)
	if err != nil {
		return err // Mengembalikan kesalahan jika koneksi gagal.
	}

	// Mengunci akses ke aplikasi untuk mencegah race condition saat menginisialisasi database.
	app.mu.Lock()
	defer app.mu.Unlock()

	// Menetapkan koneksi database ke objek database.
	db.sqldb = d
	// Menyimpan objek database ke dalam aplikasi.
	app.db = &db

	return nil // Mengembalikan nil jika inisialisasi berhasil.
}

// CreateTableIfNotExist membuat tabel baru jika tabel dengan nama yang sama belum ada di database.
// Fungsi ini digunakan untuk memastikan tabel tersedia sebelum melakukan operasi lain.
//
// Tabel yang dibuat memiliki kolom:
//   - id: Kunci utama (autoincrement).
//   - key: Teks unik yang tidak boleh NULL.
//   - value: Data dalam bentuk BLOB.
//
// Mengembalikan:
//   - error: Kesalahan jika terjadi kegagalan dalam eksekusi query.
func (db *database) CreateTableIfNotExist() error {
	// Query untuk membuat tabel jika belum ada, menggunakan SQL CREATE TABLE IF NOT EXISTS.
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS %s (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        key TEXT NOT NULL UNIQUE,
        value BLOB
    );`

	// Mengunci akses database untuk mencegah race condition saat membuat tabel.
	db.mu.Lock()
	defer db.mu.Unlock()

	// Menjalankan query untuk membuat tabel.
	_, err := db.sqldb.Exec(fmt.Sprintf(createTableQuery, db.tableName))
	if err != nil {
		return err // Mengembalikan kesalahan jika query gagal.
	}

	return nil // Mengembalikan nil jika tabel berhasil dibuat atau sudah ada.
}

// Update memperbarui nilai (value) yang terkait dengan key tertentu dalam tabel.
// Jika key tidak ditemukan, tidak ada perubahan yang akan dilakukan.
//
// Parameter:
//   - key: Kunci (key) yang ingin diperbarui.
//   - data: Data baru yang akan di-update ke dalam kolom value.
//
// Mengembalikan:
//   - error: Kesalahan jika terjadi kegagalan dalam eksekusi query.
func (db *database) Update(key string, data []byte) error {
	// Query untuk memperbarui nilai berdasarkan key yang diberikan.
	updateQuery := `
		UPDATE %s 
		SET value = ? 
		WHERE key = ?;
	`

	// Mengunci akses database untuk mencegah race condition saat memperbarui data.
	db.mu.Lock()
	defer db.mu.Unlock()

	// Menjalankan query untuk memperbarui data.
	_, err := db.sqldb.Exec(fmt.Sprintf(updateQuery, db.tableName), data, key)
	if err != nil {
		return err // Mengembalikan kesalahan jika query gagal.
	}

	return nil // Mengembalikan nil jika data berhasil diperbarui.
}

// InsertOrUpdate menambahkan data baru atau memperbarui data yang sudah ada berdasarkan key.
// Fungsi ini menggunakan ON CONFLICT untuk menangani situasi di mana key yang sama sudah ada dalam tabel.
//
// Parameter:
//   - key: Kunci unik yang digunakan untuk mengidentifikasi data.
//   - data: Data yang akan disimpan atau diperbarui.
//
// Mengembalikan:
//   - error: Kesalahan yang terjadi selama proses insert atau update.
func (db *database) InsertOrUpdate(key string, data []byte) error {
	// Mengunci akses ke database untuk menghindari kondisi balapan (race condition).
	db.mu.Lock()
	defer db.mu.Unlock()

	// Query untuk melakukan insert jika key belum ada, atau update jika key sudah ada.
	insertOrUpdateQuery := `
		INSERT INTO %s (key, value) 
		VALUES (?, ?)
		ON CONFLICT(key) 
		DO UPDATE SET value = excluded.value;
	`

	// Menjalankan query insert atau update dengan parameter key dan data.
	_, err := db.sqldb.Exec(fmt.Sprintf(insertOrUpdateQuery, db.tableName), key, data)
	if err != nil {
		return err // Mengembalikan kesalahan jika eksekusi query gagal.
	}

	return nil // Mengembalikan nil jika proses insert atau update berhasil.
}

// FindALL mengambil semua data dari tabel yang disimpan di database.
// Fungsi ini menggunakan mutex untuk memastikan akses ke database
// dilakukan secara aman dalam lingkungan multi-threaded.
//
// Mengembalikan:
//   - *[]model: Slice dari objek model yang berisi data dari tabel.
//   - error: Kesalahan jika ada masalah saat mengeksekusi query atau mengakses data.
func (db *database) FindALL() (*[]model, error) {
	// Mengunci database untuk mencegah kondisi balapan (race condition) selama pengaksesan.
	db.mu.Lock()
	defer db.mu.Unlock()

	// Menyiapkan query untuk mengambil semua data dari tabel.
	selectQuery := `SELECT id, key, value FROM %s;`

	// Menjalankan query SELECT untuk mendapatkan semua baris dari tabel yang ditentukan.
	rows, err := db.sqldb.Query(fmt.Sprintf(selectQuery, db.tableName))
	if err != nil {
		return nil, err // Mengembalikan kesalahan jika query gagal dieksekusi.
	}
	defer rows.Close() // Menutup hasil setelah selesai digunakan.

	// Inisialisasi slice untuk menampung hasil query.
	result := []model{}

	// Iterasi melalui setiap baris hasil query.
	for rows.Next() {
		r := model{} // Inisialisasi objek model untuk setiap baris.
		// Memindai kolom hasil ke dalam objek model.
		err := rows.Scan(&r.ID, &r.Key, &r.Value)
		if err != nil {
			return nil, err // Mengembalikan kesalahan jika proses pemindaian gagal.
		}
		// Menambahkan hasil pemindaian ke slice result.
		result = append(result, r)
	}

	// Mengembalikan slice dari objek model dan nil (tanpa kesalahan).
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
