impute_missing_data_BdF <- function(df,number_imputations = 1, number_iterations = 5,seed = 1234){

  df <- df %>% as.data.frame()
  all_variables = colnames(df)


  if (!('annee' %in% all_variables)){
    df <- df %>% mutate(annee = year(periode))
  }

  df2 <- df


  variables_to_impute <-  c(
    'taux_marge' ,
    'financier_court_terme' ,
    'frais_financier' ,
    'delai_fournisseur' ,
    'poids_frng' ,
    'dette_fiscale'
  )

  predictors <-  c(
    'annee',
    'outcome',
    'cut_effectif' ,
    'cut_growthrate' ,
    'lag_effectif_missing' ,
    'apart_last12_months' ,
    'apart_consommee' ,
    'apart_share_heuresconsommees' ,
    'log_cotisationdue_effectif' ,
    'log_ratio_dettecumulee_cotisation' ,
    'indicatrice_dettecumulee' ,
    'indicatrice_croissance_dettecumulee' ,
    'nb_debits' ,
    'delai' ,
    'delai_sup_6mois' ,
    'taux_marge' ,
    'financier_court_terme' ,
    'frais_financier' ,
    'delai_fournisseur' ,
    'poids_frng' ,
    'dette_fiscale'
  )


  excluded_variables <- all_variables[!(all_variables %in% predictors)]
  df2 <- df2 %>%
    group_by(siret,annee) %>%
    arrange(desc(periode)) %>%
    slice(1) %>%
    ungroup() %>%
    group_by(siret) %>%
    arrange(desc(periode)) %>%
    fill(variables_to_impute,.direction = 'up') %>%
    ungroup() %>%
    filter(annee >= 2014)

  # Initialisation
  ini <-  mice(df2,
              maxit = 0,
              printFlag = FALSE)



  # select predictors
  ini$pred[,all_variables[!(all_variables %in% predictors)]] = 0
  # select variables to impute
  ini$meth[all_variables[!(all_variables %in% variables_to_impute)]] = ''
  ini$pred[all_variables[!(all_variables %in% variables_to_impute)],] = 0

  # imputing on a yearly basis

  temp_mids <- mice(df2,
               m = number_imputations,
               maxit = number_iterations,
               method = ini$meth,
               pred = ini$pred,
               printFlag = TRUE,
               seed = seed)

    long_by_year <- mice::complete(temp_mids,'long', include = TRUE)

  # Expanding imputed values at each period
  long_by_periode <- df %>%
    filter(annee >= 2014) %>%
    select(-one_of(variables_to_impute)) %>%
    left_join(
      long_by_year %>%
      select(c('siret','annee','.imp', variables_to_impute)),by=c('siret','annee')) %>%
    arrange(.imp,siret)


    #mids <- as.mids(long_by_periode)

  # mids completion
    #mids$pred[,all_variables[!(all_variables %in% predictors)]] = 0
    #mids$meth[all_variables[!(all_variables %in% variables_to_impute)]] = ''
    #mids$pred[all_variables[!(all_variables %in% variables_to_impute)],] = 0
    #mids$seed <- seed

return(long_by_periode %>% as_tbl_time(periode))
}
