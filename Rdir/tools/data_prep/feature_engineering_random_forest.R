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



    data <- replace_na_by("montant_echeancier",data,0)
    data <- replace_na_by("delai",data,0)
    data <- replace_na_by("duree_delai",data,0)
    data <- replace_na_by("cotisation",data,0)

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

    data <- add_past_trends(data,
                            c('effectif',
                              'apart_heures_consommees',
                              'montant_part_ouvriere',
                              'montant_part_patronale'),
                            c(1,3,6,12),
                            slope = FALSE)

    names_with_na <- names(data %>% select(contains('variation')))
    for (name in names_with_na)
      data <- replace_na_by(name,data,0)

    seed <- 1234
    imputed_data_long <-  impute_missing_data_BdF(data,number_imputations,number_iterations,seed) %>%
      as_tbl_time(periode)




    return(imputed_data_long)

  }

  if (oversampling) {
    train_set <- oversample(train_set)
  }


  data_list <- c(list(aux_fun(train_set,1,5)),lapply(arguments, aux_fun, number_imputations = 5, number_iterations = 5))

  return(data_list)

}
