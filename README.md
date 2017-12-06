
# Open Signaux Faibles

## Configuration 

### Installer le projet


    $ git clone git@github.com:SGMAP-AGD/opensignauxfaibles.git
    

### Installer la librairie opensignauxfaibles

* Ouvrir RStudio
* Ouvrir le projet opensignauxfaibles et taper : 
    
    library("packrat")
    on()
    library("devtools")
    build()
    install()
    
### PostGreSQL

Les données sont chargées dans une base de données PostGreSQL

Le fichier keys.json à la racine du dossier contient les paramètres de connexion la base PostGre.

    {
      "host": [""],
      "dbname": [""],
      "port": [""],
      "id":[""],
      "pw":[""]
    }

### Charger les données brutes

* Les données brutes sont dans le répertoire [data-raw/](data-raw/).
* Le programme [import_data.Rmd](import_data.Rmd) importe toutes les données brutes dans la base PostGre avec le préfixe `table_`.
* Les fonctions appelées dans import_data.Rmd sont définies dans [R/import_data.R](R/import_data.R).

### Calcul des variables

* Le programme [compute_samples.Rmd](compute_samples.Rmd) permet de calculer l'ensemble des variables
* Les fonctions appelées dans [compute_samples.Rmd](compute_samples.Rmd) sont définies dans [R/compute_samples.R](R/compute_samples.R) 

### Calcul des prédictions

* Le programme [compute_model_0_12.Rmd]compute_model_0_12.Rmd) permet de calculer les prédictions à 12 mois.

