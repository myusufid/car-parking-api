# Car Parking API
## Design Table
```
parking_slots
-------------
block | string
slot_number | string
occupied | bool

cars
---------------
plat_number | string 
color | string 
type | string 
entry_time | timestamp
slot_id | object_id


tickets 
-----------
car_plat_number | string
slot_id | object_id
entry_time | timestamp
exit_time | timestamp
fee | int
```


## How to run server

### Docker
```shell
docker compose up -d 
```

### Air
```shell
air
```

## How to run test 

```shell
go test ./...
```
or 
```shell
make test 
```


## Endpoints 
### REGISTRASI KENDARAAN

CURL
```shell
curl --location --request POST 'localhost:4000/register' --header 'Content-Type: application/json' --data '{
    "plat_nomor": "BG23102PX",
    "warna": "black",
    "tipe": "SUV"
}'
```
Example Response
```shell
{
    "plat_nomor": "BG23102PX",
    "parking_lot": "B5",
    "tanggal_masuk": "2024-06-10 14:52"
}
```



### KENDARAAN KELUAR

CURL
```shell
curl --location --request POST 'localhost:4000/exit' --header 'Content-Type: application/json' --data '{
    "plat_nomor": "BG23102PY"
}'
```
Example Response 
```json
{
    "plat_nomor": "BG23102PY",
    "tanggal_masuk": "2024-06-10 11:46",
    "tanggal_keluar": "2024-06-10 11:46",
    "jumlah_bayar": 25000
}
```

### REPORT JUMLAH MOBIL PER TIPE MOBIL
CURL
```shell
curl --location --request GET 'localhost:4000/total_car?tipe=SUV'
```
Example Response
```json 
{
    "jumlah_kendaraan": 2
}
```

### LIST NOMOR KENDARAAN SESUAI WARNA

CURL 
```shell
curl --location --request GET 'localhost:4000/license_by_color?warna=black'
```
Example Response 
```json
{
    "plat_nomor": [
        "BG23102PY",
        "BG23102PX"
    ]
}
```

## Use MongoDB
Populate data with MongoDB
```mongodb

db.createCollection("parking_slots", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["block", "slot_number", "occupied"],
      properties: {
        block: {
          bsonType: "string",
          description: "must be a string and is required"
        },
        slot_number: {
          bsonType: "string",
          description: "must be a string and is required"
        },
        occupied: {
          bsonType: "bool",
          description: "must be a boolean and is required"
        }
      }
    }
  }
});

db.parking_slots.insertMany([
  { block: "A", slot_number: "A1", occupied: false },
  { block: "A", slot_number: "A2", occupied: false },
  { block: "A", slot_number: "A20", occupied: false },
  { block: "B", slot_number: "B1", occupied: false },
  { block: "B", slot_number: "B2", occupied: false },
  { block: "B", slot_number: "B30", occupied: false }
]);
```