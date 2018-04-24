export_top100 <- function(table_wholesample, prediction,export_filter = TRUE) {

  if(is.character(export_filter) & export_filter == 'IM'){
    bool_filter <- table_wholesample$code_naf_niveau1 == 'C'
  } else {
    bool_filter <- export_filter
  }


  der_periode <- max(prediction$periode)

  comptes_proc_collectives <-
    compute_filter_proccollectives( db = database_signauxfaibles,
          .date = as.Date(der_periode) %m+% months(1)) %>%
    collect() %>%
    .$numero_compte



  comptes_ccsf <-
    compute_filter_ccsv(db = database_signauxfaibles,
              .date = as.Date(der_periode) %m+% months(1)) %>%
    collect() %>%
    .$numero_compte

  temp_sample <-  table_wholesample %>%
    filter(bool_filter) %>%
    dplyr::inner_join(prediction %>% select(siret,periode, prob), by = c('siret','periode')) %>%
    dplyr::filter(periode == der_periode) %>%
    dplyr::mutate(
      proc_collective  = (numero_compte %in% comptes_proc_collectives)
      ) %>%
    dplyr::mutate(
      CCSF = (numero_compte %in% comptes_ccsf)
      ) %>%
    dplyr::arrange(dplyr::desc(prob))

  temp_sample %>%
    #dplyr::filter(is.na(prediction_0_12) == FALSE) %>%
    dplyr::select(
      siret,
      raison_sociale,
      prob,
      CCSF,
      proc_collective,
      cut_effectif,
      code_departement,
      region,
      libelle_naf_niveau1,
      code_ape,
      montant_part_ouvriere,
      montant_part_patronale,
      nb_debits,
      cut_growthrate,
      lag_effectif_missing,
      apart_last12_months,
      apart_consommee,
      apart_share_heuresconsommees,
      log_cotisationdue_effectif,
      log_ratio_dettecumulee_cotisation_12m,
      indicatrice_dettecumulee_12m,
      apart_effectif_moyen,
      apart_heures_consommees,
      apart_potentiel_effectif,
      growthrate_effectif,
      delai,
      delai_sup_6mois,
      indicatrice_croissance_dettecumulee,
      indicatrice_dettecumulee,
      montant_part_ouvriere_12m,
      montant_part_patronale_12m,
      lag_montant_part_ouvriere,
      lag_montant_part_patronale,
      mean_cotisation_due,
      cotisationdue_effectif,
      ratio_dettecumulee_cotisation,
      ratio_dettecumulee_cotisation_12m,
      log_effectif,
      log_growthrate_effectif,
      log_ratio_dettecumulee_cotisation,
      numero_compte
    ) %>%
    dplyr::slice(1:100) %>%
    write.table(
      row.names = F,
      dec = ',',
      sep = ';',
      file = paste0('output/algo1_',
                    der_periode,
                    if_else(is.character(first(export_filter)), first(export_filter), ''),
                    '.csv'),
      quote = T,
      append = F
    )
}
