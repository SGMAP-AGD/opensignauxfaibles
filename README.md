
# Open Signaux Faibles

## Configuration

- Le répertoire data-raw/ contient les données brutes
- Le répertoire output/ contient les résultats
- Le répertoire R/ contient l'ensemble des fonctions
- keys.json contient les paramètres de connexion la base PostGre

## Code

- Les fonctions compute_sample_ calculent les valeurs à une date données
- Les fonctions compute_wholesample_ calculent les variables sur une période donnée
- La fonction collect_wholesample() calcule l'ensemble des variables sur toute la période

- import_data.Rmd : importe les données brutes dans la base PostGre
- compute_samples.Rmd : calcule les variables pertinentes et les enregistre dans une table
- describe_samples.Rmd : calcule des statistiques descriptives sur l'échantillon
- compute_model_0_12.Rmd : calcule les prédiction pour le modèle à 12 mois

