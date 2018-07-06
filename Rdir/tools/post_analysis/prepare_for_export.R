prepare_for_export <- function(data, additional_names = NULL){

  export_names <-  c(
  'siret',
  'raison_sociale',
  'departement',
  'region',
  'prob',
  'date_ccsf',
  'proc_collective',
  'cut_effectif',
  'libelle_naf_niveau1',
  'libelle_naf_niveau5',
  'code_ape',
  'montant_part_ouvriere',
  'montant_part_patronale',
  'poids_frng',
  'taux_marge',
  'frais_financier',
  'financier_court_terme',
  'delai_fournisseur',
  'dette_fiscale',
  'apart_consommee',
  'apart_share_heuresconsommees',
  'mean_cotisation_due'
  #indicatrice_dettecumulee_12m,
  #indicatrice_croissance_dettecumulee,
  #apart_effectif_moyen,
  #apart_heures_consommees,
  #apart_potentiel_effectif,
  #ratio_dettecumulee_cotisation
  )



  export_names <- c(additional_names,export_names)



  cat("Préparation à l'export ... \n")
  cat(paste0('Dernière période connue: ',max(data$periode, na.rm = TRUE)))


  # Report des dernières infos financieres connues
  derniers_bilans_connus <- data %>%
      group_by(siret) %>%
      dplyr::arrange(periode) %>%
      summarize(poids_frng = last(na.omit(poids_frng)),
                taux_marge = last(na.omit(taux_marge)),
                frais_financier = last(na.omit(frais_financier)),
                financier_court_terme = last(na.omit(financier_court_terme)),
                delai_fournisseur = last(na.omit(delai_fournisseur)),
                dette_fiscale = last(na.omit(dette_fiscale)))

    temp_sample <-  data %>%
      select(
        -poids_frng,
        -taux_marge,
        -frais_financier,
        -financier_court_terme,
        -delai_fournisseur,
        -dette_fiscale
      ) %>%
      left_join(derniers_bilans_connus, by = 'siret') %>%
      dplyr::mutate(CCSF = date_ccsf ) %>%
      dplyr::arrange(dplyr::desc(prob))

    temp_sample <- temp_sample %>%
      mark_known_sirets(name = 'sirets_connus.csv')

    export_names <- export_names[export_names %in% names(temp_sample)]
    cat('Les variables suivantes sont absentes du dataframe:','\n')
    cat(!(export_names %in% names(temp_sample)))

    #if (is.emp)
    toExport <- temp_sample %>%
      dplyr::select(one_of(export_names))
}

