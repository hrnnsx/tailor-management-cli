# Dokumentasi Pengujian (TESTING.md)

Dokumen ini berisi informasi strategi dan hasil pengujian aplikasi Tailor Management CLI. Pengujian dilakukan melalui dua pendekatan:
1. Pengujian Otomatis (Automated Unit Testing) menggunakan database mocking untuk handler autentikasi.
2. Pengujian Manual End-to-End (Manual E2E Testing) untuk seluruh alur antarmuka baris perintah (CLI) dan transaksi sistem.

---

## 1. Pengujian Otomatis (Automated Unit Testing)

Pengujian unit otomatis diimplementasikan pada modul `handler/auth_handler.go` menggunakan pustaka bawaan Go `testing`, `go-sqlmock` untuk simulasi query database, dan `testify/assert` untuk verifikasi hasil.

### Dependensi Pustaka
* `testing` (Go Standard Library)
* `github.com/DATA-DOG/go-sqlmock`
* `github.com/stretchr/testify/assert`

### Cara Menjalankan Pengujian Otomatis
Jalankan perintah berikut pada terminal di direktori utama proyek:

```bash
go test ./handler/... -v
```

### Ringkasan Test Case Auth Handler

#### A. Pengujian Fungsi `SignIn`

1. `TestAuthHandler_SignIn_Success`
   * Tujuan: Memastikan proses masuk akun berhasil ketika email dan password sesuai.
   * Ekspektasi: Query `SELECT` mengeksekusi parameter email dengan benar dan mengembalikan data pengguna lengkap tanpa error.

2. `TestAuthHandler_SignIn_UserNotFound`
   * Tujuan: Memastikan sistem mengembalikan pesan error saat email tidak terdaftar di database.
   * Ekspektasi: Database mengembalikan `sql.ErrNoRows`, fungsi mengembalikan error "email atau password salah".

3. `TestAuthHandler_SignIn_WrongPassword`
   * Tujuan: Memastikan sistem menolak autentikasi saat password yang dimasukkan salah.
   * Ekspektasi: Data pengguna ditemukan tetapi verifikasi password gagal, fungsi mengembalikan error "email atau password salah".

4. `TestAuthHandler_SignIn_DatabaseError`
   * Tujuan: Memastikan penanganan error saat koneksi database bermasalah.
   * Ekspektasi: Database mengembalikan error koneksi, fungsi mengembalikan error dengan pesan "gagal mengambil data user...".

#### B. Pengujian Fungsi `Register`

1. `TestRegister_Success`
   * Tujuan: Memastikan pendaftaran akun baru dengan peran customer berhasil.
   * Ekspektasi: Query `INSERT` dieksekusi dengan argumen nama, email, password, dan nomor telepon, mengembalikan hasil sukses tanpa error.

2. `TestRegister_Failed_DuplicateEmail`
   * Tujuan: Memastikan sistem menolak pendaftaran jika email sudah terdaftar.
   * Ekspektasi: Database mengembalikan error `Duplicate entry`, fungsi menangkap error tersebut dan mengembalikan pesan "Email sudah terdaftar, silakan gunakan email lain!".

---

## 2. Pengujian Manual End-to-End (Manual E2E Testing)

Pengujian manual dilakukan pada seluruh alur aplikasi secara end-to-end melalui terminal interaktif untuk memverifikasi fungsionalitas CLI, validasi input, serta aturan bisnis sistem.

### A. Modul Otentikasi & Validasi Input

| ID Test | Skenario Pengujian | Langkah Pengujian | Hasil yang Diharapkan | Status |
|---|---|---|---|---|
| TC-AUTH-01 | Pendaftaran akun baru valid | Input nama, email valid, password >= 6 karakter, nomor HP diawali "08" (>= 10 digit) | Akun berhasil dibuat dan tersimpan di database | Lulus |
| TC-AUTH-02 | Validasi format email | Input email tanpa karakter `@` atau domain tidak valid | Sistem menampilkan peringatan format email salah | Lulus |
| TC-AUTH-03 | Validasi nomor telepon | Input nomor HP tidak diawali "08" atau kurang dari 10 digit | Sistem menolak input dan meminta ulang nomor HP | Lulus |
| TC-AUTH-04 | Validasi konfirmasi password | Input password dan konfirmasi password berbeda | Sistem menolak pendaftaran dan menampilkan pesan ketidakcocokan | Lulus |
| TC-AUTH-05 | Sign In berhasil | Input email dan password terdaftar yang sesuai | Pengguna berhasil masuk ke menu sesuai role | Lulus |
| TC-AUTH-06 | Sign In gagal | Input password salah atau email tidak terdaftar | Sistem menampilkan pesan "email atau password salah" | Lulus |

### B. Modul Customer (Pelanggan)

| ID Test | Skenario Pengujian | Langkah Pengujian | Hasil yang Diharapkan | Status |
|---|---|---|---|---|
| TC-CUST-01 | Pengisian ukuran badan | Memilih buat ukuran baru, mengisikan angka tinggi, dada, dan pinggang | Data ukuran tersimpan dan estimasi ukuran (XS-XXL) dikalkulasi otomatis | Lulus |
| TC-CUST-02 | Pemilihan kain & corak stok cukup | Memilih kain dan corak dengan stok cm yang mencukupi | Pesanan terbuat dengan status `unpaid`, stok kain otomatis terpotong | Lulus |
| TC-CUST-03 | Proteksi kain stok tidak cukup | Memilih kain dan corak dengan status "SOLD OUT" / stok < kebutuhan | Pilihan diblokir dan sistem menampilkan peringatan stok tidak cukup | Lulus |
| TC-CUST-04 | Cek tagihan pesanan | Memilih menu Cek Tagihan | Menampilkan daftar pesanan berstatus `unpaid` beserta total harga | Lulus |
| TC-CUST-05 | Konfirmasi pembayaran | Memilih pesanan `unpaid` dan mengunggah/memasukkan data bukti bayar | Status pesanan/pembayaran berubah menjadi `pending` | Lulus |
| TC-CUST-06 | Cek status pesanan | Memilih menu Cek Status Pesanan | Menampilkan status pengerjaan dan status pembayaran terkini | Lulus |

### C. Modul Admin

| ID Test | Skenario Pengujian | Langkah Pengujian | Hasil yang Diharapkan | Status |
|---|---|---|---|---|
| TC-ADM-01 | Assign Order ke Worker | Memilih pesanan baru dan menugaskan ke worker yang tersedia | Pesanan terhubung ke ID worker, worker terpilih siap mengerjakan | Lulus |
| TC-ADM-02 | Verifikasi Pembayaran | Memeriksa daftar pembayaran `pending` lalu menyetujui transaksi | Status pembayaran berubah dari `pending` menjadi `paid` | Lulus |
| TC-ADM-03 | Tambah Kain Dasar Baru | Input nama kain baru dan harga per cm | Data kain baru tersimpan di database | Lulus |
| TC-ADM-04 | Tambah Corak & Stok Awal | Menghubungkan kain dasar dengan corak baru serta mengisikan stok cm | Kombinasi `fabric_patterns` berhasil ditambahkan | Lulus |
| TC-ADM-05 | Restock Kain | Menambahkan kuantitas stok cm pada varian kain yang sudah ada | Jumlah stok kain bertambah sesuai input | Lulus |
| TC-ADM-06 | Laporan Penjualan | Memilih menu Report | Menampilkan total order, pesanan selesai, dan total pendapatan | Lulus |

### D. Modul Worker (Penjahit)

| ID Test | Skenario Pengujian | Langkah Pengujian | Hasil yang Diharapkan | Status |
|---|---|---|---|---|
| TC-WRK-01 | Lihat Order Saya | Worker masuk menu utama dan memilih Lihat Order Saya | Menampilkan daftar pesanan yang ditugaskan khusus ke worker tersebut | Lulus |
| TC-WRK-02 | Update Progres Pengerjaan | Memperbarui tahap dari "Order diterima" hingga "Proses jahit" | Status pengerjaan pesanan terbarui di database | Lulus |
| TC-WRK-03 | Syarat Pengerjaan Finishing | Memperbarui status ke "Finishing" pada pesanan yang belum `paid` | Sistem memberikan pemberitahuan/penahanan hingga status menjadi `paid` | Lulus |