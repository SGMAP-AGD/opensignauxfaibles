  library("dplyr")
  dbconnection <- mongolite::mongo(collection = "Entreprise", db = 'opensignauxfaibles', verbose = TRUE, url = 'mongodb://localhost:27017')

  data <- dbconnection$aggregate('[{"$addFields": {"batch": {"$objectToArray": "$value.batch"}}}]')
