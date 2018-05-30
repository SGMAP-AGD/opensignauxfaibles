export_top <- function(prediction,table_wholesample, top = 100,export_filter = TRUE,name = 'test') {

  if(is.character(export_filter) & length(export_filter == 1) & export_filter == 'IM'){
    bool_filter <- table_wholesample$code_naf_niveau1 == 'C'
  } else if (is_character(export_filter)) {
    bool_filter = table_wholesample$siret %in% export_filter
  } else {
    bool_filter <- export_filter
  }


  der_periode <- max(prediction$periode)

  # comptes_proc_collectives <-
  #   compute_filter_proccollectives( db = database_signauxfaibles,
  #         .date = as.Date(der_periode) %m+% months(1)) %>%
  #   collect() %>%
  #   .$numero_compte
  #
  # comptes_ccsf <-
  #   compute_filter_ccsv(db = database_signauxfaibles,
  #             .date = as.Date(der_periode) %m+% months(1)) %>%
  #   collect() %>%
  #   .$numero_compte

  # Report des dernières infos financieres connues
  derniers_bilans_connus <- table_wholesample %>%
    group_by(siret) %>%
    dplyr::arrange(periode) %>%
    summarize(poids_frng = last(na.omit(poids_frng)),
              taux_marge = last(na.omit(taux_marge)),
              frais_financier = last(na.omit(frais_financier)),
              financier_court_terme = last(na.omit(financier_court_terme)),
              delai_fournisseur = last(na.omit(delai_fournisseur)),
              dette_fiscale = last(na.omit(dette_fiscale)))

  temp_sample <-  table_wholesample %>%
    filter(bool_filter) %>%
    dplyr::inner_join(prediction %>% select(siret, periode, prob),
                      by = c('siret', 'periode')) %>%
    dplyr::filter(periode == der_periode) %>%
    select(
      -poids_frng,-taux_marge,-frais_financier,-financier_court_terme,-delai_fournisseur,-dette_fiscale
    ) %>%
    left_join(derniers_bilans_connus, by = 'siret') %>%
   # dplyr::mutate(proc_collective  = (numero_compte %in% comptes_proc_collectives)) %>%
   # dplyr::mutate(CCSF = (numero_compte %in% comptes_ccsf)) %>%
    dplyr::mutate(proc_collective = date_defaillance) %>%
    dplyr::mutate(CCSF = date_ccsf ) %>%
    dplyr::arrange(dplyr::desc(prob))

  toExport <- temp_sample %>%
    #dplyr::filter(is.na(prediction_0_12) == FALSE) %>%
    dplyr::select(
      siret,
      raison_sociale,
      departement,
      region,
      prob,
      date_ccsf,
      date_defaillance,
      cut_effectif,
      #libelle_naf_niveau1,
      #libelle_naf_niveau5,
      code_ape,
      montant_part_ouvriere,
      montant_part_patronale,
      poids_frng,
      taux_marge,
      frais_financier,
      financier_court_terme,
      delai_fournisseur,
      dette_fiscale,
      apart_consommee,
      apart_share_heuresconsommees,
      mean_cotisation_due
      #indicatrice_dettecumulee_12m,
      #indicatrice_croissance_dettecumulee,
      #apart_effectif_moyen,
      #apart_heures_consommees,
      #apart_potentiel_effectif,
      #ratio_dettecumulee_cotisation
    ) %>%
    dplyr::slice(1:top)


    write.table(toExport,
      row.names = F,
      dec = ',',
      sep = ';',
      file = paste0('../output/table_predictions_algo2_',
                    der_periode %m+% months(1),
                    name,
                    '.csv'),
      quote = T,
      append = F
    )


    file = paste0(
      '../output/extraction_algo2_',
      der_periode %m+% months(1),
      '.xlsx'
    )

    wb <- XLConnect::loadWorkbook( file, create = TRUE)
    XLConnect::removeSheet(wb,sheet = name)
    XLConnect::createSheet(wb, name=name)

    XLConnect::writeWorksheet(
      wb,
      as.data.frame(toExport),
      name,
      header = TRUE,
      rownames = NULL
    )

    saveWorkbook(wb, file)
}
