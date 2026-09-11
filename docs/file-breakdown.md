# Penjelasan Struktur Berkas (File Breakdown)

Dokumen ini menjelaskan fungsi dan tanggung jawab dari setiap berkas kode program dalam proyek Tailor Management CLI.

---

## 1. Modul Utama & Konfigurasi

### `database.go`
Menangani koneksi awal ke database MySQL.
* Mengambil konfigurasi lingkungan (`DB_USER`, `DB_PASSWORD`, `DB_HOST`, `DB_PORT`, `DB_NAME`) menggunakan pustaka `godotenv`.
* Membuka koneksi database dengan parameter `parseTime=true` untuk menangani tipe data waktu.
* Memastikan koneksi berjalan dengan melakukan pengecekan `Ping()`.

### `color.go`
Menyediakan utilitas visual untuk antarmuka terminal.
* Menyimpan konstanta warna teks berbasis ANSI escape codes (seperti merah, hijau, kuning, cyan, dan lavender).
* Digunakan oleh seluruh modul CLI untuk membedakan status pesan (sukses, error, peringatan, atau petunjuk input).

### `menu.go`
Berperan sebagai pengatur navigasi (router) utama aplikasi.
* Mengelola perulangan utama (*infinite loop*) agar antarmuka tidak langsung keluar setelah menyelesaikan suatu tindakan.
* Menyediakan layar `MainMenu` untuk akses awal (Sign In, Sign Up, Exit).
* Mengarahkan pengguna ke menu spesifik sesuai peran (`CustomerMenu`, `AdminMenu`, atau `WorkerMenu`) setelah proses otentikasi berhasil.
* Menginisialisasi *handler* untuk menyuntikkan koneksi database ke modul fitur.

---

## 2. Modul Autentikasi

### `signin.go`
Menangani alur masuk akun pengguna.
* Meminta input email dan password.
* Menyembunyikan tampilan input password di terminal menggunakan pustaka `promptui`.
* Memvalidasi format email menggunakan ekspresi reguler (regex) sebelum mencocokkan kredensial ke database.

### `signup.go`
Menangani pendaftaran akun baru.
* Meminta input nama, email, password, konfirmasi password, dan nomor telepon.
* Menerapkan validasi input:
  * Nama tidak boleh kosong.
  * Email harus sesuai format standar.
  * Password minimal 6 karakter dan harus cocok dengan konfirmasi password.
  * Nomor telepon minimal 10 digit dan wajib diawali dengan "08".

---

## 3. Modul Logika Peran (Role Features)

### `customer.go`
Mengelola antarmuka dan interaksi untuk pengguna dengan peran Customer.
* **BikinBajuCLI:** Meminta atau menggunakan data ukuran tubuh pelanggan, memicu kalkulasi kain, dan menampilkan daftar pilihan bahan serta corak.
* **selectPatternWithSoldOutGuard:** Memeriksa ketersediaan stok kain. Jika kebutuhan kain melebihi stok yang ada, opsi akan ditandai "SOLD OUT" dan pilihan diblokir.
* **CheckBill & PayOrderCLI:** Menampilkan daftar pesanan yang belum dibayar (*unpaid*) dan mengisikan konfirmasi pembayaran untuk diverifikasi admin.

### `order.go`
Menangani logika bisnis dan transaksi terkait pemesanan.
* **determineSize:** Menentukan klasifikasi ukuran baju (XS hingga XXL) berdasarkan data pengukuran fisik pelanggan.
* **submitOrderTransaction:** Menjalankan transaksi database (`tx.Begin()`) untuk memotong stok kain pada tabel `fabric_patterns`, membuat kode pesanan unik berformat `ORD-xxx`, dan menyimpan detail pesanan. Jika terjadi kesalahan pada salah satu proses, transaksi di-*rollback*.

### `admin.go`
Mengelola antarmuka dan fungsi administratif untuk peran Admin.
* **CheckAllOrder & AssignOrder:** Menampilkan pesanan baru dan menugaskannya kepada penjahit yang sedang tersedia.
* **CheckPayment:** Menampilkan daftar pembayaran berstatus *pending* dan mengubah statusnya menjadi *paid* setelah diverifikasi.
* **Laporan Penjualan:** Menampilkan ringkasan total pesanan, pesanan selesai, status pembayaran, dan total pendapatan.
* **Manajemen Kain:** Menambah jenis kain dasar (`CreateFabricCLI`), menambah kombinasi kain dan corak beserta stok awal (`CreateFabricPatternCLI`), serta menambah stok kain (`RestockFabricCLI`).

### `worker.go`
Mengelola antarmuka dan alur kerja untuk peran Worker (Penjahit).
* **ShowMyOrders:** Menampilkan daftar pesanan yang khusus ditugaskan kepada penjahit yang sedang login.
* **UpdateOrderProgressCLI:** Memperbarui tahap progres pengerjaan baju secara berurutan (*Order diterima*, *Persiapan bahan*, *Pemotongan kain*, *Proses jahit*, *Finishing*).