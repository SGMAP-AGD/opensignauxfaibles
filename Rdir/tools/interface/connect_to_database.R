connect_to_database <- function(collection){

  library(mongolite)
  library(dplyr)

  dbconnection <- mongo(collection = collection, db = 'opensignauxfaibles', url = 'mongodb://localhost:27017')


  data <- dbconnection$aggregate('[{"$unwind":{"path": "$value"}}]')$value %>%
    mutate(
      cut_growthrate = forcats::fct_relevel(
        cut_growthrate,
        c("stable", "moins de 20%", "moins 20 à 5%",
          "plus 5 à 20%", "plus 20%", "manquant")),
      cut_effectif = forcats::fct_relevel(cut_effectif),
      outcome_0_12 = factor(outcome_0_12,
                            levels = c("non-default", "default")
      )
    )
}
