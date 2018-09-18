count_etab_entr <- function(df){
  nb_etab   <- n_distinct(df %>% select(siret))
  nb_entr <- n_distinct(df %>% select(siren))

  cat(nb_etab,' établissements dans ', nb_entr, ' entreprises','\n')


}
