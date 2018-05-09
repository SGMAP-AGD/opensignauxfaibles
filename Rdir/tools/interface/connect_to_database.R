connect_to_database <- function(collection){
  dbconnection <- mongo(collection = collection, db = 'opensignauxfaibles', verbose = TRUE, url = 'mongodb://localhost:27017')

  data <- dbconnection$aggregate('[{"$unwind":{"path": "$value"}}]')$value %>%
    mutate(
      cut_growthrate = forcats::fct_relevel(
        cut_growthrate,
        c("stable", "moins_de_20p", "moins_20_a_5p",
          "plus_5_a_20p", "plus_20p", "manquant")),
      cut_effectif = forcats::fct_relevel(cut_effectif),
      outcome_0_12 = factor(outcome_0_12,
                            levels = c(FALSE,TRUE),
                            labels =  c("default", "non_default")
      )
    )

  return(data)

}
