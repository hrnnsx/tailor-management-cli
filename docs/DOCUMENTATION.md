# Dokumentasi Sistem Tailor Management CLI

Sistem Tailor Management CLI adalah aplikasi antarmuka baris perintah (CLI) berbasis Go untuk mengelola operasional bisnis penjahit kustom, mulai dari pendaftaran pelanggan, pemesanan kustom, alokasi stok bahan, penugasan penjahit, hingga verifikasi pembayaran.

---

## Daftar Fitur Utama (Berdasarkan Role)

### 1. Fitur Umum & Otentikasi
* **Sign Up:** Pendaftaran akun baru dengan pengisian Nama, Email, Password, dan Nomor HP. Sistem memvalidasi format email (regex) serta format nomor HP (minimal 10 digit dan diawali "08").
* **Sign In:** Otentikasi pengguna menggunakan Email dan Password. Sistem mengarahkan pengguna ke antarmuka menu yang sesuai dengan peran (role) masing-masing.

### 2. Fitur Customer (Pelanggan)
* **Pemesanan Kustom (Bikin Baju):** Pelanggan memasukkan atau memilih profil ukuran badan (tinggi badan, lingkar dada, lingkar pinggang). Sistem secara otomatis mengkalkulasi kebutuhan kain dalam centimeter dan estimasi ukuran (XS hingga XXL). Pelanggan kemudian memilih bahan dan corak kain berdasarkan ketersediaan stok.
* **Cek Tagihan & Pembayaran:** Pelanggan dapat melihat rincian tagihan pesanan berstatus belum dibayar (unpaid) dan mengajukan konfirmasi pembayaran.
* **Cek Status Pesanan:** Pelanggan dapat memantau status pesanan dan progres pengerjaan terkini.

### 3. Fitur Admin
* **Assign Order:** Admin melihat daftar pesanan baru yang belum ditugaskan, kemudian menugaskannya ke penjahit (worker) yang tersedia.
* **Verifikasi Pembayaran:** Admin memeriksa daftar pembayaran pelanggan berstatus pending untuk diverifikasi (disetujui) menjadi paid.
* **Laporan Penjualan (Report):** Admin melihat ringkasan performa bisnis, mencakup total pesanan, pesanan selesai, status pembayaran, serta total pendapatan.
* **Manajemen Inventaris Kain:**
  * Menambahkan jenis kain dasar dan harga per centimeter.
  * Menambahkan kombinasi varian kain dengan corak (pattern) baru beserta stok awal.
  * Menambahkan stok (restock) pada kombinasi kain dan corak yang ada.

### 4. Fitur Worker (Penjahit)
* **Lihat Order Saya:** Pekerja melihat daftar pesanan yang secara spesifik ditugaskan kepada mereka oleh Admin.
* **Update Progress Order:** Pekerja memperbarui tahap penyelesaian pengerjaan baju secara bertahap.

---

## Alur Kerja Sistem (Workflow)

Sistem mengintegrasikan proses bisnis melalui alur kerja berikut:

### 1. Alur Pemesanan & Alokasi Stok (Customer)
1. Customer melakukan registrasi/login dan memilih menu "Bikin Baju".
2. Customer memasukkan data ukuran tubuh (tinggi, lingkar dada, lingkar pinggang).
3. Sistem secara otomatis menentukan estimasi ukuran (XS-XXL) dan menghitung jumlah kain yang dibutuhkan (dalam cm).
4. Customer memilih varian kain dan corak. Sistem memeriksa ketersediaan stok:
   * Jika stok cukup, transaksi diproses: stok kain dipotong secara otomatis di database, dan kode pesanan (ORD-xxx) dibuat dengan status pembayaran `unpaid`.
   * Jika stok tidak cukup, sistem menolak pilihan dengan status "SOLD OUT".

### 2. Alur Penugasan & Verifikasi (Admin)
1. Admin memeriksa daftar pesanan baru yang belum memiliki penjahit.
2. Admin menugaskan pesanan tersebut kepada penjahit yang sedang tersedia (availability = true).
3. Admin memeriksa pengajuan pembayaran berstatus `pending` dari Customer.
4. Admin memverifikasi bukti pembayaran. Setelah diverifikasi, status pembayaran pesanan berubah menjadi `paid`.

### 3. Alur Pengerjaan & Pembaruan Progres (Worker)
1. Worker membuka menu "Lihat Order Saya" untuk melihat daftar tugas yang diterima.
2. Worker memperbarui status progres pengerjaan secara berurutan:
   * Order diterima
   * Persiapan bahan
   * Pemotongan kain
   * Proses jahit
   * Finishing
3. Pada tahap penyelesaian akhir (Finishing), sistem memastikan pembayaran pelanggan telah berstatus `paid` sebelum pesanan diselesaikan sepenuhnya.

---

## Status State Machine

### Status Pembayaran (Payment Status)
* `unpaid` : Pesanan baru berhasil dibuat oleh Customer.
* `pending` : Customer telah mengajukan konfirmasi pembayaran.
* `paid` : Admin telah memverifikasi dan menyetujui pembayaran.

### Status Pengerjaan (Order Progress Status)
* `Order diterima` -> `Persiapan bahan` -> `Pemotongan kain` -> `Proses jahit` -> `Finishing`
