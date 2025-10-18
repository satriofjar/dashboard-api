# dashboard-api

## Endpoint
```
GET /summary?periode=YYYY-MM
```
Contoh:
```
GET 127.0.0.1:8000/summary?periode=2024-01
```

## Deskripsi
Endpoint digunakan untuk menampilkan summary penggunaan ruangan perkantor (office) berdasarkan periode bulan dan tahun yang diminta melalui parameter periode.

## Contoh Response
```
{
        "officeName": "UID JAYA",
        "rooms": [
            {
                "roomName": "Ruang Borobudur",
                "persentasePemakaian": 7.5,
                "nominalKonsumsi": 5940000,
                "snackSiang": 102,
                "makanSiang": 62,
                "snackSore": 102
            },
            {
                "roomName": "Ruang Prambanan",
                "persentasePemakaian": 7.5,
                "nominalKonsumsi": 4200000,
                "snackSiang": 72,
                "makanSiang": 72,
                "snackSore": 30
            },
            {
                "roomName": "Ruang Mendhut",
                "persentasePemakaian": 3.125,
                "nominalKonsumsi": 1650000,
                "snackSiang": 15,
                "makanSiang": 27,
                "snackSore": 27
            }
        ]
    },
``` 