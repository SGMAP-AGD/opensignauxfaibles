connect_to_database <- function(collection){
  dbconnection <- mongo(collection = collection, db = 'opensignauxfaibles', verbose = TRUE, url = 'mongodb://localhost:27017')

  data <- dbconnection$aggregate('[{"$unwind":{"path": "$value"}}]')$value %>%
    mutate(
      cut_growthrate = forcats::fct_relevel(
        cut_growthrate,
        c("moins_de_20p", "moins_20_a_5p","stable",
          "plus_5_a_20p", "plus_20p", "manquant")),
      cut_effectif = forcats::fct_relevel(cut_effectif),
      outcome_0_12 = forcats::fct_relevel(outcome_0_12,
                              c("default", "non_default")
      )
    )

  return(data %>%
           mutate(periode = as.Date(periode)) %>%
           arrange(periode) %>%
           as_tbl_time(periode))
}