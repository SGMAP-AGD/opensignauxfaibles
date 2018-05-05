spread_relative_periods <- function(wholesample, feature,n_period, by = 1){
  # COMPLETE MISSING VALUES
  # IN COMPUTATION RELATIVE MONTHS TO CURRENT DATE
  # Ajouter evolution par annee !



  # fill missing periods
  wholesample <- wholesample %>%
    mutate(periode = as.Date(periode)) %>%
    tidyr::complete(siret, periode = seq.Date(from = min(as.Date(periode)),to = max(as.Date(periode)),by = "month")) %>%
    group_by(siret) %>%
    arrange(siret, periode) %>%
    mutate_(temp_var = feature)



  for (i in 1:n_period) {
    wholesample <- wholesample %>%
      mutate(!!paste0(feature, i) := lag(temp_var, i))
  }



    wholesample <- wholesample %>%
    select(-temp_var) %>%
    ungroup()

  # Spread relative periods
  # wholesample <- wholesample %>%
  #   rowid_to_column("ID") %>% #to avoid that spread merges rows
  #   mutate(tempkey = months_before_default) %>%
  #   spread(key = tempkey, value = feature, sep = '_') %>%
  #   rename_at(vars(starts_with("tempkey")),funs(paste0(feature,str_match(.,pattern = '[0-9]+'))))%>%
  #   select(-ID)

  return(wholesample)
}
