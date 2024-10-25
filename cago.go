// Copyright (c) 2024, Jasakode Authors.
// All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

/*
	Package cago provides documentation for Go's predeclared identifiers.
	The items documented here are not actually in package cago
	but their descriptions here allow godoc to present documentation
	for the language's special identifiers.
*/

package cago

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/jasakode/cago/lib"
	"github.com/jasakode/cago/store"
)

// Config menyimpan konfigurasi utama aplikasi yang berhubungan dengan database dan penggunaan memori.
//
// Field-field:
//   - Path: Lokasi file database di sistem. Jika path tidak ditentukan, aplikasi akan menggunakan database sementara yang datanya hilang setelah proses berakhir.
//   - MAX_MEM: Batas memori maksimum yang diperbolehkan untuk aplikasi, dinyatakan dalam bit. Default adalah 8.589.934.592 bit (1 GB).
//   - MIN_MEM_ALLOCATION: Jumlah memori minimum yang dialokasikan, dinyatakan dalam bit. Default adalah 8.388.608 bit (1 MB).
type Config struct {
	// Path ke file database. Jika kosong, data akan hilang setelah proses selesai.
	// File seperti "database.db" akan menyimpan data untuk mengantisipasi jika
	// program terhenti, sehingga data yang telah dicache dapat dimuat ulang.
	Path string
	// Memori maksimal yang diperbolehkan digunakan (dalam bit).
	// 8.388.608 bit = 1 MB.
	// default: 8589934592 bit (1 GB).
	MAX_MEM uint
	// Memori minimal yang akan dialokasikan (dalam bit).
	// 8.388.608 bit = 1 MB.
	// default: 8388608 bit (1 MB).
	MIN_MEM_ALLOCATION uint64
	// Jika true, data yang ditambahkan paling awal akan dihapus
	// ketika batas memori maksimal tercapai.
	// default : false
	EvictOldestOnMaxMem bool
	// Timeout untuk pemeriksaan entri yang kedaluwarsa (dalam milidetik).
	// Ini menentukan interval waktu antara setiap pemeriksaan data dalam cache.
	// Default: 10000 (10 detik).
	TimeoutCheck uint64
}

// Struktur `App` digunakan untuk mengelola seluruh aplikasi, termasuk konfigurasi, database, dan data cache.
//
// Field-field:
//   - mu: Mutex untuk memastikan operasi thread-safe pada aplikasi, mencegah race condition.
//   - start: Waktu start aplikasi dalam format Unix timestamp (uint64).
//   - config: Objek konfigurasi aplikasi (Config) yang menyimpan pengaturan aplikasi.
//   - db: Pointer ke objek database yang mengelola koneksi dan operasi database.
//   - data: Cache data dalam bentuk map, yang menggunakan string sebagai key dan store.Store sebagai value.
type App struct {
	mu        sync.Mutex             // Mutex untuk memastikan thread-safe akses ke field dalam struct App.
	db        *database              // Pointer ke objek database yang digunakan aplikasi.
	data      map[string]store.Store // Cache data aplikasi dalam map, dengan string sebagai key dan store.Store sebagai value.
	data_size uint64                 // ukuran total data berserta key
	start     uint64                 // Timestamp yang merepresentasikan waktu mulai aplikasi.
	config    Config                 // Konfigurasi aplikasi, berisi pengaturan penting.
}

// Variabel global `app` adalah instance dari struct `App` yang digunakan di seluruh aplikasi.
var app App = App{}

// New menginisialisasi aplikasi dengan konfigurasi yang diberikan.
// Jika konfigurasi tidak disediakan, aplikasi akan menggunakan nilai default.
// Mengatur data cache dan memulai waktu aplikasi.
// Jika Path untuk database diberikan, aplikasi akan menginisialisasi
// database dan memuat data dari database ke dalam cache.
func New(config ...Config) error {
	app = App{}
	// Mengatur konfigurasi default
	app.config = Config{}
	// Jika ada konfigurasi yang diberikan, gunakan konfigurasi tersebut
	if len(config) > 0 {
		app.config = config[0]
	}
	// Menginisialisasi aplikasi
	app.init()
	// Jika Path database tidak kosong, inisialisasi database
	if app.config.Path != "" {
		if err := app.InitializeDB(); err != nil {
			return err
		}
		// Membuat tabel jika belum ada
		if err := app.db.CreateTableIfNotExist(); err != nil {
			return err
		}
		// Mengambil semua data dari database
		rows, err := app.db.FindALL()
		if err != nil {
			return err
		}
		// Memasukkan data yang diambil dari database ke dalam cache
		for i := range *rows {
			val := (*rows)[i]
			// Menambahkan data ke cache berdasarkan key tertentu
			app.data[val.Key] = store.ParseStore(val.Value)
		}
		return nil
	}
	return nil
}

// runNode menjalankan proses yang terus-menerus untuk memeriksa data dalam cache.
// Fungsi ini berfungsi untuk menghapus entri yang sudah kedaluwarsa berdasarkan MaxAge yang ditentukan.
func (app *App) runNode() {
	// Loop tanpa henti untuk terus memeriksa data dalam cache
	for {
		// Tidur selama waktu yang ditentukan oleh TimeoutCheck dalam milidetik
		// untuk mengatur interval pemeriksaan entri yang kedaluwarsa.
		time.Sleep(time.Duration(app.config.TimeoutCheck) * time.Millisecond)

		// Iterasi melalui setiap entri dalam cache
		for k, v := range app.data {
			// Memeriksa apakah MaxAge untuk entri ini tidak sama dengan 0
			if v.MaxAge() != 0 {
				// Jika waktu sekarang dikurangi waktu pembuatan entri masih dalam batas waktu
				if uint64(time.Now().UnixMilli())-v.CreateAt() >= v.MaxAge() {
					// Menghapus entri dari cache berdasarkan kunci
					Remove(k)
				}
			}
		}
	}
}

// init menginisialisasi nilai maksimum dan minimum memori untuk aplikasi.
// Jika MAX_MEM dan MIN_MEM_ALLOCATION tidak ditentukan, akan diatur
// ke nilai default yang sesuai.
func (app *App) init() {
	// Menentukan nilai MAX_MEM default jika belum ditentukan
	if app.config.MAX_MEM == 0 {
		app.config.MAX_MEM = 8388608 * 1204 // 10 MB
	}
	// Menentukan nilai MIN_MEM_ALLOCATION default jika belum ditentukan
	if app.config.MIN_MEM_ALLOCATION == 0 {
		app.config.MIN_MEM_ALLOCATION = 8388608 // 1 MB
	}
	if app.config.TimeoutCheck == 0 {
		app.config.TimeoutCheck = 10000 // 1 MB
	}

	// Menginisialisasi data cache untuk menyimpan store
	app.data = make(map[string]store.Store)
	// Menyimpan waktu mulai aplikasi dalam milidetik
	app.start = uint64(time.Now().UnixMilli())
	app.data_size = uint64(0)

	go app.runNode()
}

// TotalSize menghitung ukuran total dari semua key dan nilai yang disimpan dalam map app.data.
// Ukuran dihitung sebagai jumlah byte dari panjang string key dan panjang nilai (store)
// yang disimpan. Fungsi ini efisien dan tidak menggunakan banyak memori tambahan.
//
// Mengembalikan:
// - uint64: Total ukuran data (key dan value) dalam byte.
func Size() uint64 {
	var totalSize uint64
	// Iterasi melalui setiap pasangan key-value di dalam map data
	for key, store := range app.data {
		// Hitung ukuran key (string) dalam byte
		totalSize += uint64(len(key))
		// Hitung ukuran nilai (store) dengan fungsi Length(true)
		// Length(true) menghitung ukuran store secara keseluruhan
		totalSize += store.Length(true)
	}
	return totalSize
}

// Set menyimpan nilai ke dalam store dengan key yang diberikan.
// Fungsi ini juga dapat menerima parameter opsional untuk menentukan maxAge.
// Nilai yang disimpan harus sesuai dengan tipe yang didefinisikan oleh interface store.Compare.
//
// Parameter:
//   - key (string): Key unik yang digunakan untuk mengidentifikasi nilai dalam store.
//   - value (store.Compare): Nilai yang akan disimpan. Harus memiliki tipe data yang sesuai
//     dengan interface Compare, seperti integer, float, string, atau tipe apapun yang diizinkan.
//   - maxAge (opsional) (uint64): Waktu maksimal dalam milidetik selama nilai akan disimpan.
//     Jika tidak disertakan, nilai ini akan diabaikan.
//
// Mengembalikan:
// - error: Kesalahan jika terjadi selama penyimpanan data.
func Set(key string, value store.Compare, maxAge ...uint64) error {
	app.mu.Lock()
	defer app.mu.Unlock()
	_, ok := app.data[key]
	if ok {
		return fmt.Errorf("data already exists")
	}
	switch v := any(value).(type) {
	case string:
		data := store.NewStore([]byte(v), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
		if app.db != nil {
			if err := app.db.InsertOrUpdate(key, data); err != nil {
				return err
			}
		}
	case int:
		data := store.NewStore(lib.Int64ToByte(int64(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
		if app.db != nil {
			if err := app.db.InsertOrUpdate(key, data); err != nil {
				return err
			}
		}
	case int8:
		data := store.NewStore(lib.Int8ToByte(int8(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
		if app.db != nil {
			if err := app.db.InsertOrUpdate(key, data); err != nil {
				return err
			}
		}
	case int16:
		data := store.NewStore(lib.Int16ToByte(int16(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
		if app.db != nil {
			if err := app.db.InsertOrUpdate(key, data); err != nil {
				return err
			}
		}
	case int32:
		data := store.NewStore(lib.Int32ToByte(int32(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
		if app.db != nil {
			if err := app.db.InsertOrUpdate(key, data); err != nil {
				return err
			}
		}
	case int64:
		data := store.NewStore(lib.Int64ToByte(int64(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
		if app.db != nil {
			if err := app.db.InsertOrUpdate(key, data); err != nil {
				return err
			}
		}
	case uint:
		data := store.NewStore(lib.Uint64ToByte(uint64(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
		if app.db != nil {
			if err := app.db.InsertOrUpdate(key, data); err != nil {
				return err
			}
		}
	case uint8:
		data := store.NewStore(lib.Uint8ToByte(uint8(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
		if app.db != nil {
			if err := app.db.InsertOrUpdate(key, data); err != nil {
				return err
			}
		}
	case uint16:
		data := store.NewStore(lib.Uint16ToByte(uint16(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
		if app.db != nil {
			if err := app.db.InsertOrUpdate(key, data); err != nil {
				return err
			}
		}
	case uint32:
		data := store.NewStore(lib.Uint32ToByte(uint32(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
		if app.db != nil {
			if err := app.db.InsertOrUpdate(key, data); err != nil {
				return err
			}
		}
	case uint64:
		data := store.NewStore(lib.Uint64ToByte(uint64(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
		if app.db != nil {
			if err := app.db.InsertOrUpdate(key, data); err != nil {
				return err
			}
		}
	case float32, float64:
		by, err := json.Marshal(value)
		if err != nil {
			return err
		}
		data := store.NewStore(by, maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
		if app.db != nil {
			if err := app.db.InsertOrUpdate(key, data); err != nil {
				return err
			}
		}
	case any:
		by, err := json.Marshal(value)
		if err != nil {
			return err
		}
		data := store.NewStore(by, maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
		if app.db != nil {
			if err := app.db.InsertOrUpdate(key, data); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

// Get mengambil nilai dari store berdasarkan key yang diberikan.
// Fungsi ini mengembalikan pointer ke nilai yang ditemukan. Jika tidak ada nilai
// yang cocok dengan key, akan mengembalikan nil.
//
// Parameter:
//   - key (string): Key unik yang digunakan untuk mencari nilai dalam store.
//
// Tipe Parameter:
//   - K (store.Compare): Tipe data yang diharapkan sesuai dengan interface Compare,
//     seperti integer, float, string, atau tipe apapun yang diizinkan.
//
// Mengembalikan:
//   - *K: Pointer ke nilai yang diambil dari store. Jika nilai tidak ditemukan,
//     akan mengembalikan nil.
func Get[K store.Compare](key string) *K {
	app.mu.Lock()
	defer app.mu.Unlock()

	value, ok := app.data[key]
	if !ok {
		return nil // Mengembalikan nil jika key tidak ada
	}

	var result K

	// Menangani setiap tipe dalam switch
	switch any(result).(type) {
	case string:
		result = any(value.Text()).(K)
	case int:
		intValue, err := value.Int()
		if err != nil {
			fmt.Println("Error retrieving int:", err)
			return nil // Tangani kesalahan dengan baik
		}
		result = any(intValue).(K)
	case int8:
		intValue, err := value.Int()
		if err != nil {
			fmt.Println("Error retrieving int8:", err)
			return nil // Tangani kesalahan dengan baik
		}
		result = any(int8(intValue)).(K) // Konversi jika perlu
	case int16:
		intValue, err := value.Int()
		if err != nil {
			fmt.Println("Error retrieving int16:", err)
			return nil // Tangani kesalahan dengan baik
		}
		result = any(int16(intValue)).(K) // Konversi jika perlu
	case int32:
		intValue, err := value.Int()
		if err != nil {
			fmt.Println("Error retrieving int32:", err)
			return nil // Tangani kesalahan dengan baik
		}
		result = any(int32(intValue)).(K) // Konversi jika perlu
	case int64:
		intValue, err := value.Int()
		if err != nil {
			fmt.Println("Error retrieving int64:", err)
			return nil // Tangani kesalahan dengan baik
		}
		result = any(int64(intValue)).(K) // Konversi jika perlu
	case uint:
		intValue, err := value.Int()
		if err != nil {
			fmt.Println("Error retrieving uint:", err)
			return nil // Tangani kesalahan dengan baik
		}
		result = any(uint(intValue)).(K) // Konversi jika perlu
	case uint8:
		intValue, err := value.Int()
		if err != nil {
			fmt.Println("Error retrieving uint8:", err)
			return nil // Tangani kesalahan dengan baik
		}
		result = any(uint8(intValue)).(K) // Konversi jika perlu
	case uint16:
		intValue, err := value.Int()
		if err != nil {
			fmt.Println("Error retrieving uint16:", err)
			return nil // Tangani kesalahan dengan baik
		}
		result = any(uint16(intValue)).(K) // Konversi jika perlu
	case uint32:
		intValue, err := value.Int()
		if err != nil {
			fmt.Println("Error retrieving uint32:", err)
			return nil // Tangani kesalahan dengan baik
		}
		result = any(uint32(intValue)).(K) // Konversi jika perlu
	case uint64:
		intValue, err := value.Int()
		if err != nil {
			fmt.Println("Error retrieving uint64:", err)
			return nil // Tangani kesalahan dengan baik
		}
		result = any(uint64(intValue)).(K) // Konversi jika perlu
	case float32:
		intValue, err := value.Int()
		if err != nil {
			fmt.Println("Error retrieving float32:", err)
			return nil // Tangani kesalahan dengan baik
		}
		result = any(float32(intValue)).(K) // Konversi jika perlu
	case float64:
		intValue, err := value.Int()
		if err != nil {
			fmt.Println("Error retrieving float64:", err)
			return nil // Tangani kesalahan dengan baik
		}
		result = any(float64(intValue)).(K) // Konversi jika perlu
	default:
		err := value.JSON(&result)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return nil // Tangani kesalahan dengan baik
		}
	}

	return &result
}

// Exist memeriksa apakah nilai dengan key yang diberikan ada dalam store.
// Fungsi ini mengembalikan true jika key ditemukan, dan false jika tidak.
//
// Parameter:
//   - key (string): Key unik yang digunakan untuk memeriksa keberadaan nilai
//     dalam store.
//
// Mengembalikan:
// - bool: True jika nilai dengan key ditemukan; False jika tidak ditemukan.
func Exist(key string) bool {
	app.mu.Lock()
	defer app.mu.Unlock()
	_, ok := app.data[key]
	return ok
}

// Put menggantikan atau membuat nilai baru ke dalam store dengan key yang diberikan.
// Jika key sudah ada, nilai yang lama akan digantikan dengan nilai baru.
// Fungsi ini juga dapat menerima parameter opsional untuk menentukan maxAge.
//
// Parameter:
//   - key (string): Key unik yang digunakan untuk mengidentifikasi nilai dalam store.
//   - value (store.Compare): Nilai yang akan disimpan. Harus memiliki tipe data yang sesuai
//     dengan interface Compare, seperti integer, float, string, atau tipe apapun yang diizinkan.
//   - maxAge (opsional) (uint64): Waktu maksimal dalam milidetik selama nilai akan disimpan.
//     Jika tidak disertakan, nilai ini akan disimpan tanpa batasan waktu.
//
// Mengembalikan:
// - error: Kesalahan jika terjadi selama proses penggantian atau penyimpanan data.
func Put(key string, value store.Compare, maxAge ...uint64) error {
	app.mu.Lock()
	defer app.mu.Unlock()
	if len(maxAge) == 0 {
		old, ok := app.data[key]
		if ok {
			maxAge = append(maxAge, old.MaxAge())
		}
	}
	switch v := any(value).(type) {
	case string:
		data := store.NewStore([]byte(v), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
		if app.db != nil {
			if err := app.db.InsertOrUpdate(key, data); err != nil {
				return err
			}
		}
	case int:
		data := store.NewStore(lib.Int64ToByte(int64(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
		if app.db != nil {
			if err := app.db.InsertOrUpdate(key, data); err != nil {
				return err
			}
		}
	case int8:
		data := store.NewStore(lib.Int8ToByte(int8(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
		if app.db != nil {
			if err := app.db.InsertOrUpdate(key, data); err != nil {
				return err
			}
		}
	case int16:
		data := store.NewStore(lib.Int16ToByte(int16(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
		if app.db != nil {
			if err := app.db.InsertOrUpdate(key, data); err != nil {
				return err
			}
		}
	case int32:
		data := store.NewStore(lib.Int32ToByte(int32(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
		if app.db != nil {
			if err := app.db.InsertOrUpdate(key, data); err != nil {
				return err
			}
		}
	case int64:
		data := store.NewStore(lib.Int64ToByte(int64(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
		if app.db != nil {
			if err := app.db.InsertOrUpdate(key, data); err != nil {
				return err
			}
		}
	case uint:
		data := store.NewStore(lib.Uint64ToByte(uint64(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
		if app.db != nil {
			if err := app.db.InsertOrUpdate(key, data); err != nil {
				return err
			}
		}
	case uint8:
		data := store.NewStore(lib.Uint8ToByte(uint8(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
		if app.db != nil {
			if err := app.db.InsertOrUpdate(key, data); err != nil {
				return err
			}
		}
	case uint16:
		data := store.NewStore(lib.Uint16ToByte(uint16(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
		if app.db != nil {
			if err := app.db.InsertOrUpdate(key, data); err != nil {
				return err
			}
		}
	case uint32:
		data := store.NewStore(lib.Uint32ToByte(uint32(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
		if app.db != nil {
			if err := app.db.InsertOrUpdate(key, data); err != nil {
				return err
			}
		}
	case uint64:
		data := store.NewStore(lib.Uint64ToByte(uint64(v)), maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
		if app.db != nil {
			if err := app.db.InsertOrUpdate(key, data); err != nil {
				return err
			}
		}
	case float32, float64:
		by, err := json.Marshal(value)
		if err != nil {
			return err
		}
		data := store.NewStore(by, maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
		if app.db != nil {
			if err := app.db.InsertOrUpdate(key, data); err != nil {
				return err
			}
		}
	case any:
		by, err := json.Marshal(value)
		if err != nil {
			return err
		}
		data := store.NewStore(by, maxAge...)
		if app.db != nil {
			app.db.InsertOrUpdate(key, data)
		}
		app.data[key] = data
		if app.db != nil {
			if err := app.db.InsertOrUpdate(key, data); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

// Remove menghapus nilai yang terkait dengan key yang diberikan dari store.
// Fungsi ini juga menghapus data dari database jika ada.
//
// Parameter:
//   - key (string): Key unik yang digunakan untuk menghapus nilai dalam store.
//
// Mengembalikan:
// - bool: True jika key berhasil dihapus; False jika key tidak ditemukan.
func Remove(key string) bool {
	app.mu.Lock()
	defer app.mu.Unlock()
	_, ok := app.data[key]
	delete(app.data, key)
	if app.db != nil {
		if err := app.db.RemoveByKey(key); err != nil {
			fmt.Println(err.Error())
		}
	}
	return ok
}

// Clear menghapus semua nilai yang tersimpan dalam store dan database.
// Fungsi ini mengosongkan map data dan, jika ada, memanggil fungsi untuk
// menghapus semua data dari database.
//
// Mengembalikan:
// - error: Kesalahan jika terjadi selama proses penghapusan data dari database.
func Clear() error {
	app.mu.Lock()
	defer app.mu.Unlock()
	app.data = make(map[string]store.Store)
	if app.db != nil {
		return app.db.RemoveAll()
	}
	return nil
}
