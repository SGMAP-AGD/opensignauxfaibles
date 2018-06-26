# Documentation API open signaux faibles

## Lots d'intégration
### POST /api/batch
Création d'un batch

```json
input: {
    "batch": "1802",
    "description": "Lot de février 2018",
    "files": {
        "bdf": ["/1802/bdf/fichier.xls"],
        "admin_urssaf": ["listBou.csv", "listFC.csv"]
    }
}
```