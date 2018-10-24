# fakeData
Fabrication de dataset anonymisé pour Signaux Faibles

## Objectif: 
- Fournir un environnement de démonstration comportant un jeu de données
- Fabrication d'un dataset de test

## Méthode d'anonymisation
- Remplacement des siret avec des siret aléatoires
- Randomisation des informations sirène à l'exception de la région et du code APE
- Application d'un coefficient aléatoire différent sur chaque type de données chiffrées (urssaf et financières)

## Utilisation
- Créer config.toml à partir de config.toml.example
- go run *.go
- profit