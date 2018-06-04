prepare_for_export <- function(data){


  cat("Préparation à l'export ... \n")
  cat(paste0('Dernière période connue: ',max(data$periode)))

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

    # Import libelles naf. Ne devrait pas être ici ! FIX ME
    libelle_naf <- readxl::read_excel(
          path = rprojroot::find_rstudio_root_file(file.path('..','data-raw','naf','naf2008_5_niveaux.xls')),
          sheet = "naf2008_5_niveaux",
          skip = 1,
          col_names = c("code_naf_niveau5", "code_naf_niveau4", "code_naf_niveau3", "code_naf_niveau2", "code_naf_niveau1")
        ) %>%
          dplyr::select(code_naf_niveau5, code_naf_niveau1) %>%
          dplyr::left_join(
            y = readxl::read_excel(
              path = rprojroot::find_rstudio_root_file(file.path('..','data-raw','naf','naf2008_liste_n5.xls')),
              sheet = "Feuil1",
              skip = 3,
              col_names = c("code_naf_niveau5", "libelle_naf_niveau5")
            ),
            by = "code_naf_niveau5"
          ) %>%
          dplyr::left_join(
            y = readxl::read_excel(
              path = rprojroot::find_rstudio_root_file(file.path('..','data-raw','naf','naf2008_liste_n1.xls')),
              sheet = "Feuil1",
              skip = 3,
              col_names = c("code_naf_niveau1", "libelle_naf_niveau1")
            ),
            by = "code_naf_niveau1"
          )  %>%
          dplyr::mutate(
            code_naf_niveau5 = stringr::str_replace(
              string = code_naf_niveau5,
              pattern = "([[:digit:]]{2})\\.([[:digit:]]{2}[[:upper:]]{1})",
              replacement = "\\1\\2")
          )

    temp_sample <- temp_sample %>%
      left_join(libelle_naf, by = c('code_ape'= 'code_naf_niveau5'))


    toExport <- temp_sample %>%
      dplyr::select(
        siret,
        raison_sociale,
        departement,
        region,
        prob,
        date_ccsf,
        proc_collective,
        cut_effectif,
        libelle_naf_niveau1,
        libelle_naf_niveau5,
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
      )
}

