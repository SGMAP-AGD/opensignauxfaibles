feature_engineering_random_forest <- function(train_set, ..., oversampling) {
  arguments <- list(...)

  aux_fun <- function(data, number_imputations, number_iterations) {   #1 check vars
    assertthat::assert_that(all(
      c(
        'siret',
        'periode',
        'debut_activite',
        'activite_saisonniere',
        'productif',
        'etat_proc_collective',
        'libelle_naf_niveau1'
      ) %in% names(data)
    ))

    #2 new computed fields
    data <- data %>%
      mutate(age = year(periode) - debut_activite)
    replace_na_by_0 <- function(name,data) {
      data[is.na(data[,name]), name] = 0
      return(data)
    }

    data <- replace_na_by_0("montant_echeancier",data)
    data <- replace_na_by_0("delai",data)
    data <- replace_na_by_0("duree_delai",data)
    data <- replace_na_by_0("cotisation",data)

    data <- data %>%
      mutate(
        activite_saisonniere = ifelse(activite_saisonniere == "S", 1, 0),
        productif = ifelse(productif == "O", 1, 0),
        etat_proc_collective = as.numeric(factor(
          etat_proc_collective,
          levels = c(
            'in_bonis',
            'cession',
            'sauvegarde',
            'continuation',
            'plan_sauvegarde',
            'plan_redressement',
            'liquidation'
          )
        ))
      )

    libelle_naf_simplifie <- data$libelle_naf_niveau1
    libelle_naf_simplifie[libelle_naf_simplifie %in% c(
      'Enseignement',
      'Activités extra-territoriales',
      'Administration publique',
      'Agriculture, sylviculture et pêche')] <-
      'autre'
    data <- data %>%
      mutate(libelle_naf_simplifie = libelle_naf_simplifie)

    seed <- 1234
    imputed_data_long <-  impute_missing_data_BdF(data,number_imputations,number_iterations,seed) %>%
      as_tbl_time(periode)
    #imputed_data <- imputed_data_long %>% filter(.imp==1) %>% select(-.imp) %>% as_tbl_time(periode)

    return(imputed_data_long)

  }

  if (oversampling) {
    train_set <- oversample(train_set)
  }


  data_list <- c(list(aux_fun(train_set,1,5)),lapply(arguments, aux_fun, number_imputations = 5, number_iterations = 5))

  return(data_list)

}
