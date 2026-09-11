# MODELS.md

Dokumentasi ini mencakup seluruh struct, enum, dan DTO (Data Transfer Object) yang digunakan dalam aplikasi[cite: 1].

---

## 1. Enums & Custom Types

### PaymentStatusType
Status pembayaran transaksi[cite: 1].
* `Paid` ("paid"): Pembayaran telah berhasil dilunasi[cite: 1].
* `Unpaid` ("unpaid"): Pembayaran belum dilunasi[cite: 1].

### OrderStatusType
Status alur pengerjaan pesanan[cite: 1].
* `Pending` ("pending"): Pesanan baru masuk dan menunggu verifikasi[cite: 1].
* `InProgress` ("in progress"): Pesanan sedang dikerjakan oleh penjahit[cite: 1].
* `WaitingPayment` ("waiting payment"): Pesanan menunggu proses pembayaran[cite: 1].
* `Finished` ("finished"): Pesanan telah selesai diproses[cite: 1].

---

## 2. Core Entities

### User
Menyimpan data akun pengguna sistem[cite: 1].
* `ID` (int): Identifikasi unik pengguna[cite: 1].
* `Name` (string): Nama lengkap pengguna[cite: 1].
* `Email` (string): Alamat email unik[cite: 1].
* `Password` (string): Kata sandi terenkripsi[cite: 1].
* `Role` (string): Peran akses pengguna (e.g., admin, customer)[cite: 1].
* `Phone` (string): Nomor telepon[cite: 1].
* `CreatedAt` (time.Time): Waktu pembuatan akun[cite: 1].

### UserMeasurement
Menyimpan profil ukuran tubuh milik pengguna[cite: 1].
* `ID` (int): Identifikasi unik data ukuran[cite: 1].
* `UserID` (int): Foreign key ke model User[cite: 1].
* `Title` (string): Label/nama profil ukuran (misal: "Ukuran Formal")[cite: 1].
* `HeightCM` (float64): Tinggi badan dalam satuan sentimeter[cite: 1].
* `ChestCircumference` (float64): Lingkar dada dalam sentimeter[cite: 1].
* `WaistCircumference` (float64): Lingkar pinggang dalam sentimeter[cite: 1].
* `CreatedAt` (time.Time): Waktu pembuatan data ukuran[cite: 1].

### Fabric
Menyimpan master data bahan/kain[cite: 1].
* `ID` (int): ID unik bahan kain[cite: 1].
* `Name` (string): Nama jenis kain[cite: 1].
* `PricePerCM` (float64): Harga dasar kain per sentimeter[cite: 1].

### Pattern
Menyimpan master data pola atau desain pakaian[cite: 1].
* `ID` (int): ID unik pola[cite: 1].
* `Name` (string): Nama pola/model pakaian[cite: 1].

### FabricPattern
Relasi antara kain dan pola beserta ketersediaan stok dan penyesuaian harga[cite: 1].
* `ID` (int): ID unik relasi kain-pola[cite: 1].
* `FabricID` (int): ID referensi kain[cite: 1].
* `FabricName` (string): Nama kain[cite: 1].
* `PatternID` (int): ID referensi pola[cite: 1].
* `PatternName` (string): Nama pola[cite: 1].
* `StockCM` (int): Sisa stok kain dalam sentimeter[cite: 1].
* `PricePerCM` (float64): Harga per sentimeter untuk kombinasi ini[cite: 1].

### Payment
Menyimpan data catatan pembayaran[cite: 1].
* `ID` (int): ID unik transaksi pembayaran[cite: 1].
* `OrderID` (int): ID pesanan terkait[cite: 1].
* `OrderCode` (string): Kode unik pesanan[cite: 1].
* `CustomerName` (string): Nama pembeli[cite: 1].
* `Amount` (float64): Total nominal pembayaran[cite: 1].
* `Status` (string): Status verifikasi pembayaran[cite: 1].
* `CreatedAt` (time.Time): Waktu pencatatan pembayaran[cite: 1].

---

## 3. Data Transfer Objects (DTO) & Views

### CustomerOrder
Proyeksi data pesanan untuk tampilan sisi pelanggan[cite: 1].
* `ID` (int): ID pesanan[cite: 1].
* `OrderCode` (string): Kode acuan pesanan[cite: 1].
* `DeterminedSize` (string): Ukuran pakaian yang ditentukan[cite: 1].
* `TotalPrice` (float64): Total biaya pesanan[cite: 1].
* `PaymentStatus` (PaymentStatusType): Status pembayaran saat ini[cite: 1].
* `Status` (OrderStatusType): Status progres pesanan[cite: 1].
* `Progress` (string): Deskripsi singkat progres pekerjaan[cite: 1].
* `CreatedAt` (time.Time): Tanggal pembuatan pesanan[cite: 1].

### AdminOrder
Proyeksi data pesanan untuk dashboard manajemen admin[cite: 1].
* `ID` (int): ID pesanan[cite: 1].
* `OrderCode` (string): Kode unik pesanan[cite: 1].
* `CustomerName` (string): Nama pelanggan[cite: 1].
* `DeterminedSize` (string): Ukuran baju yang dikerjakan[cite: 1].
* `CMUsed` (int): Total panjang kain yang digunakan (dalam cm)[cite: 1].
* `AssignedWorkerID` (*int): Pointer ID penjahit yang ditugaskan (opsional/nullable)[cite: 1].
* `Status` (string): Status pesanan[cite: 1].
* `Progress` (string): Catatan perkembangan pengerjaan[cite: 1].
* `CreatedAt` (time.Time): Waktu pembuatan pesanan[cite: 1].

### OrderSummary
Ringkasan perhitungan rincian pesanan[cite: 1].
* `FabricName` (string): Nama bahan yang dipilih[cite: 1].
* `PatternName` (string): Nama pola yang dipilih[cite: 1].
* `Size` (string): Ukuran akhir yang dipilih/dihitung[cite: 1].
* `RequiredCM` (int): Kebutuhan panjang kain (cm)[cite: 1].
* `PricePerCM` (float64): Tarif harga per sentimeter[cite: 1].
* `TotalPrice` (float64): Total harga keseluruhan[cite: 1].

### FabricPatternOption
Struktur data pilihan kain dan pola pada form pemesanan[cite: 1].
* `FabricPatternID` (int), `FabricID` (int), `PatternID` (int)[cite: 1]
* `PatternName` (string), `StockCM` (int), `PricePerCM` (float64)[cite: 1]

### AvailableWorker
Proyeksi singkat untuk mendaftar pekerja/penjahit yang siap ditugaskan[cite: 1].
* `ID` (int): ID penjahit[cite: 1].
* `Name` (string): Nama penjahit[cite: 1].

### SalesReport
Proyeksi data agregat untuk laporan penjualan/statistik[cite: 1].
* `TotalOrders` (int): Total keseluruhan pesanan[cite: 1].
* `TotalRevenue` (float64): Total akumulasi pendapatan[cite: 1].
* `PaidOrders` (int): Jumlah pesanan yang sudah dibayar[cite: 1].
* `UnpaidOrders` (int): Jumlah pesanan yang belum dibayar[cite: 1].
* `FinishedOrders` (int): Jumlah pesanan yang telah selesai[cite: 1].