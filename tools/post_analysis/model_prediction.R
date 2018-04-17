model_prediction <-  function(.model, sample, do.print = FALSE) {
  # Moche (car données plus vieilles que l'entraînement)
  # Filtrer pour accélérer le calcul

  ########### which model ?

  ########### Compute prediction
  out_pred_raw <- .model %>%
    predict(newdata = sample)
  out_pred_prob <- .model %>%
    predict(newdata = sample,type = "prob") %>%
    select(default) %>%
    rename(predicted_default = default) %>%
    cbind(out_pred_raw)

  augmented_sample <- sample %>%
    bind_cols(out_pred_prob)


  if (do.print) {
    # Export top-100
    output_prediction_0_12 %>%
      dplyr::anti_join(
        y = compute_filter_proccollectives(db = database_signauxfaibles,
                                           .date = "2018-02-01"),
        by = "numero_compte",
        copy = TRUE
      ) %>%
      dplyr::anti_join(
        y = compute_filter_ccsv(db = database_signauxfaibles,
                                .date = "2018-02-01"),
        by = c("numero_compte"),
        copy = TRUE
      ) %>%
      dplyr::arrange(dplyr::desc(prediction_0_12)) %>%
      dplyr::filter(is.na(prediction_0_12) == FALSE) %>%
      dplyr::select(
        raison_sociale,
        siret,
        numero_compte,
        libelle_naf_niveau1,
        code_departement,
        region,
        code_ape,
        cut_effectif,
        cut_growthrate,
        lag_effectif_missing,
        apart_last12_months,
        apart_consommee,
        apart_share_heuresconsommees,
        log_cotisationdue_effectif,
        log_ratio_dettecumulee_cotisation_12m,
        indicatrice_dettecumulee_12m,
        prediction_0_12,
        apart_effectif_moyen,
        apart_heures_consommees,
        apart_potentiel_effectif,
        growthrate_effectif,
        delai,
        delai_sup_6mois,
        indicatrice_croissance_dettecumulee,
        indicatrice_dettecumulee,
        montant_part_ouvriere,
        montant_part_ouvriere_12m,
        montant_part_patronale,
        montant_part_patronale_12m,
        lag_montant_part_ouvriere,
        lag_montant_part_patronale,
        mean_cotisation_due,
        cotisationdue_effectif,
        ratio_dettecumulee_cotisation,
        ratio_dettecumulee_cotisation_12m,
        nb_debits,
        log_effectif,
        log_growthrate_effectif,
        log_ratio_dettecumulee_cotisation
      ) %>%
      dplyr::slice(1:100) %>%
      write.table(
        row.names = F,
        dec = ',',
        sep = ';',
        file = "output/table_predictions_2018_02_22v4_BDF.csv",
        quote = T,
        append = F
      )


    #### Export manufactured goods
    output_prediction_0_12 %>%
      dplyr::anti_join(
        y = compute_filter_proccollectives(db = database_signauxfaibles,
                                           .date = "2018-02-01"),
        by = "numero_compte",
        copy = TRUE
      ) %>%
      dplyr::anti_join(
        y = compute_filter_ccsv(db = database_signauxfaibles,
                                .date = "2018-02-01"),
        by = c("numero_compte"),
        copy = TRUE
      ) %>%
      dplyr::arrange(dplyr::desc(prediction_0_12)) %>%
      dplyr::filter(is.na(prediction_0_12) == FALSE, code_naf_niveau1 == 'C') %>%
      dplyr::select(
        raison_sociale,
        siret,
        numero_compte,
        libelle_naf_niveau1,
        code_departement,
        region,
        code_ape,
        cut_effectif,
        cut_growthrate,
        lag_effectif_missing,
        apart_last12_months,
        apart_consommee,
        apart_share_heuresconsommees,
        log_cotisationdue_effectif,
        log_ratio_dettecumulee_cotisation_12m,
        indicatrice_dettecumulee_12m,
        prediction_0_12,
        apart_effectif_moyen,
        apart_heures_consommees,
        apart_potentiel_effectif,
        growthrate_effectif,
        delai,
        delai_sup_6mois,
        indicatrice_croissance_dettecumulee,
        indicatrice_dettecumulee,
        montant_part_ouvriere,
        montant_part_ouvriere_12m,
        montant_part_patronale,
        montant_part_patronale_12m,
        lag_montant_part_ouvriere,
        lag_montant_part_patronale,
        mean_cotisation_due,
        cotisationdue_effectif,
        ratio_dettecumulee_cotisation,
        ratio_dettecumulee_cotisation_12m,
        nb_debits,
        log_effectif,
        log_growthrate_effectif,
        log_ratio_dettecumulee_cotisation
      ) %>%
      dplyr::slice(1:100) %>%
      write.table(
        row.names = F,
        dec = ',',
        sep = ';',
        file = "output/table_predictions_2018_02_22v4_im.csv",
        quote = T,
        append = F
      )


    # Export dunno what
    output_prediction_0_12 %>%
      dplyr::anti_join(
        y = compute_filter_proccollectives(db = database_signauxfaibles,
                                           .date = "2018-02-01"),
        by = "numero_compte",
        copy = TRUE
      ) %>%
      dplyr::anti_join(
        y = compute_filter_ccsv(db = database_signauxfaibles,
                                .date = "2018-02-01"),
        by = c("numero_compte"),
        copy = TRUE
      ) %>%
      dplyr::arrange(dplyr::desc(prediction_0_12)) %>% dplyr::select(
        raison_sociale,
        siret,
        numero_compte,
        libelle_naf_niveau1,
        code_departement,
        region,
        code_ape,
        cut_effectif,
        cut_growthrate,
        lag_effectif_missing,
        apart_last12_months,
        apart_consommee,
        apart_share_heuresconsommees,
        log_cotisationdue_effectif,
        log_ratio_dettecumulee_cotisation_12m,
        indicatrice_dettecumulee_12m,
        prediction_0_12,
        apart_effectif_moyen,
        apart_heures_consommees,
        apart_potentiel_effectif,
        growthrate_effectif,
        delai,
        delai_sup_6mois,
        indicatrice_croissance_dettecumulee,
        indicatrice_dettecumulee,
        montant_part_ouvriere,
        montant_part_ouvriere_12m,
        montant_part_patronale,
        montant_part_patronale_12m,
        lag_montant_part_ouvriere,
        lag_montant_part_patronale,
        mean_cotisation_due,
        cotisationdue_effectif,
        ratio_dettecumulee_cotisation,
        ratio_dettecumulee_cotisation_12m,
        nb_debits,
        log_effectif,
        log_growthrate_effectif,
        log_ratio_dettecumulee_cotisation
      ) %>%
      dplyr::filter(is.na(prediction_0_12) == FALSE) %>%
      write.table(
        row.names = F,
        dec = ',',
        sep = ';',
        file = "output/table_predictions_2018_02_22v4_all.csv",
        quote = T,
        append = F
      )
  }

  return(augmented_sample)
}
