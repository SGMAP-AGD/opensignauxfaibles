generator <-
  function(data,
           lookback_yearlydata,
           lookback_monthlydata,
           delay,
           sirets,
           date_inf,
           date_sup,
           seed,
           batch_size = 128,
           step = 1) {
    date_inf <- as.Date(date_inf)
    date_sup <- as.Date(date_sup)
    lookback_monthlydata <- months(lookback_monthlydata)
    lookback_yearlydata <- years(lookback_yearlydata)
    date_seq <-
      seq(date_inf %m-% lookback_monthlydata, date_sup, by = 'months')

    ## FIX ME
    data <- data %>%
      filter(is.na(log_cotisationdue_effectif))

    data <- data %>%
      replace_na(
        replace = list(
          poids_frng = -100,
          taux_marge = -100,
          delai_fournisseur = -100,
          dette_fiscale = -100,
          financier_court_terme = -100,
          frais_financier = -100
        )
      )

    data_monthly <- data %>%
      select(siret,
             periode,
             effectif,
             log_cotisationdue_effectif,
             log_ratio_dettecumulee_cotisation_12m) %>%
      tidyr::complete(periode = date_seq,
                      siret,
                      fill = list(
                        effectif = -100,
                        log_cotisationdue_effectif = -100,
                        log_ratio_dettecumulee_cotisation_12m = -100
                      )) %>%
      arrange(siret,periode) %>%
      group_by(siret) %>%
      spread(key = nested(siret,periode), value = effectif)


  }
