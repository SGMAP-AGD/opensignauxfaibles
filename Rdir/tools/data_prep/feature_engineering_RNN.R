feature_engineering_RNN <- function(train_set, ..., oversampling) {
  arguments <- list(...)

  aux_scale <- train_set %>%
    normalize_df()

  norm_means = aux_scale[['means']]
  norm_stds = aux_scale[['stds']]
  rm(aux_scale)

  aux_fun <- function(data) {

    #1 check vars
    assertthat::assert_that(all(
      c(
        'siret',
        'periode',
        'debut_activite',
        'activite_saisonniere',
        'productif',
        'etat_proc_collective'
      ) %in% names(data)
    ))

    #2 new computed fields
    data <- data %>%
      mutate(age = year(periode) - debut_activite)

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

    # Normalization
    #
    aux <- data %>%
      normalize_df(means = norm_means, stds = norm_stds)

    data_norm <- aux[['data']]

    return(data_norm)

  }

  if (oversampling) {
    train_set <- oversample(train_set)
  }

  data_list <- lapply(c(list(train_set), arguments), aux_fun)

  return(data_list)

}
