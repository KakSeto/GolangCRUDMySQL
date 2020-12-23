package main

/* Import Kebutuhan Database,
download driver jika blm punya : go get -u github.com/go-sql-driver/mysql*/

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

/* Buat Struct dengan variabel sejumlah kolom tabel
PERHATIKAN PENGGUNAAN HURUF BESAR / KECIL, pastikan KONSISTEN */
type Pegawai struct {
	Id      int
	Nama    string
	Alamat  string
	Jabatan string
}

/* Koneksi dengan Database */
func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "golang"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

/* Inisialisasi Folder 'Form', karena semua view diletakkan di form/ */
var tmpl = template.Must(template.ParseGlob("form/*"))

/* Function Index - Mengambil data dari database */
func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM Pegawai ORDER BY id ASC")
	if err != nil {
		panic(err.Error())
	}

	/* Penggunaan Struct ,dipanggil*/
	emp := Pegawai{}
	res := []Pegawai{}

	/* Looping untuk mengambil data */
	for selDB.Next() {

		//buat variabel untuk menampung data
		//sesuaikan sama nama kolom database (huruf kecil)
		var id int
		var nama, alamat, jabatan string

		err = selDB.Scan(&id, &nama, &alamat, &jabatan)
		if err != nil {
			panic(err.Error())
		}

		//kanan nama var struct - kiri nama kolom database yang diinisialisasikan diatas
		emp.Id = id
		emp.Nama = nama
		emp.Alamat = alamat
		emp.Jabatan = jabatan
		res = append(res, emp)
	}

	//kirim data ke view Index
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

/* Function Show - Detail data yang dipilih */
func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Pegawai WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Pegawai{}
	for selDB.Next() {
		//buat variabel untuk menampung data
		//sesuaikan sama nama kolom database (huruf kecil)
		var id int
		var nama, alamat, jabatan string

		err = selDB.Scan(&id, &nama, &alamat, &jabatan)
		if err != nil {
			panic(err.Error())
		}

		//kanan nama var struct - kiri nama kolom database yang diinisialisasikan diatas
		emp.Id = id
		emp.Nama = nama
		emp.Alamat = alamat
		emp.Jabatan = jabatan
	}
	tmpl.ExecuteTemplate(w, "Show", emp)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Pegawai WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Pegawai{}
	for selDB.Next() {
		//buat variabel untuk menampung data
		//sesuaikan sama nama kolom database (huruf kecil)
		var id int
		var nama, alamat, jabatan string

		err = selDB.Scan(&id, &nama, &alamat, &jabatan)
		if err != nil {
			panic(err.Error())
		}

		//kanan nama var struct - kiri nama kolom database yang diinisialisasikan diatas
		emp.Id = id
		emp.Nama = nama
		emp.Alamat = alamat
		emp.Jabatan = jabatan
	}
	tmpl.ExecuteTemplate(w, "Edit", emp)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		nama := r.FormValue("nama")
		alamat := r.FormValue("alamat")
		jabatan := r.FormValue("jabatan")
		insForm, err := db.Prepare("INSERT INTO Pegawai(nama, alamat, jabatan) VALUES(?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(nama, alamat, jabatan)
		log.Println("INSERT: Nama: " + nama + " | Alamat: " + alamat + " | Jabatan: " + jabatan)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		nama := r.FormValue("nama")
		alamat := r.FormValue("alamat")
		jabatan := r.FormValue("jabatan")
		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE Pegawai SET nama=?, alamat=?, jabatan=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(nama, alamat, jabatan, id)
		log.Println("UPDATE: Nama: " + nama + " | Alamat: " + alamat + " | Jabatan: " + jabatan)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	emp := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Pegawai WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)
}
