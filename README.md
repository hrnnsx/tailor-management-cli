# Tailor Management CLI

Aplikasi Command Line Interface (CLI) berbasis Go untuk mengelola operasional bisnis penjahit. Sistem ini menangani pemesanan pelanggan, inventaris kain, profil ukuran badan, penugasan pekerja, serta verifikasi pembayaran.

## Prasyarat

Sebelum menjalankan aplikasi, pastikan perangkat Anda telah terpasang:
* Go (versi 1.20 atau yang lebih baru)
* MySQL atau MariaDB
* Terminal atau Command Prompt

## Instalasi dan Cara Penggunaan

1. Clone repositori dan masuk ke direktori proyek:
   ```bash
   git clone <url-repository>
   cd tailor-management-cli
   ```

2. Unduh seluruh dependensi Go:
   ```bash
   go mod tidy
   ```

3. Buat file `.env` pada direktori utama proyek dengan isi sebagai berikut:
   ```env
   DB_USER=root
   DB_PASSWORD=password_anda
   DB_HOST=127.0.0.1
   DB_PORT=3306
   DB_NAME=tailor_db
   ```

4. Eksekusi skrip SQL untuk membuat struktur database dan tabel:
   ```bash
   mysql -u root -p < database-schema/query.sql
   ```

5. Jalankan aplikasi:
   ```bash
   go run .
   ```

## Alur Sistem

Sistem ini membagi akses pengguna ke dalam tiga peran: **Customer**, **Admin**, dan **Worker**.

### Transisi Status Pembayaran
1. **Unpaid**: Status default saat pelanggan selesai membuat pesanan.
2. **Pending**: Status setelah pelanggan memasukkan data/bukti pembayaran.
3. **Paid**: Status setelah admin memverifikasi dan menyetujui pembayaran.

### Transisi Status Pengerjaan (Worker)
Setelah admin menugaskan pesanan ke pekerja yang tersedia, pekerja memperbarui progres secara berurutan:
1. Order diterima
2. Persiapan bahan
3. Pemotongan kain
4. Proses jahit
5. Finishing (Perubahan status akhir memerlukan status pembayaran "Paid")

## Struktur Database

### Daftar Tabel
* **users**: Menyimpan data akun pengguna dengan peran `customer`, `admin`, atau `worker`.
* **workers**: Tabel ekstensi untuk pengguna dengan peran `worker` untuk mencatat status ketersediaan (*availability*).
* **user_measurements**: Menyimpan profil ukuran badan pelanggan (tinggi, lingkar dada, lingkar pinggang, lingkar pinggul, lebar bahu, dan panjang lengan).
* **fabrics**: Inventaris bahan kain dasar beserta harga per meter.
* **patterns**: Daftar motif atau corak pakaian.
* **fabric_patterns**: Tabel penghubung antara kain dan motif yang mencatat sisa stok dalam satuan sentimeter.
* **size_requirements**: Tabel acuan ukuran standar (XS-XXL) dan kebutuhan panjang kain (cm).
* **orders**: Mencatat data transaksi pesanan, kode pesanan, ukuran hasil kalkulasi, penggunaan kain, harga total, status pembayaran, dan status pengerjaan.
* **payments**: Mencatat transaksi pembayaran, jumlah transfer, bukti bayar, status verifikasi, dan admin yang memverifikasi.

### Trigger Database
* **trg_after_user_insert**: Otomatis menambahkan data ke tabel `workers` dengan status ketersediaan aktif jika akun baru terdaftar dengan peran 'worker'.
* **trg_after_user_update**: Otomatis menyinkronkan tabel `workers` jika status peran pengguna diperbarui (menambahkan data jika berubah menjadi 'worker', atau menghapus data jika berubah dari 'worker').

Selengkapnya di :  [database-schema](./database-schema/database-schema.md)