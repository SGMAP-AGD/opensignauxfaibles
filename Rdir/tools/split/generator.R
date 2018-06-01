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

    library(reshape2)
    date_inf <- as.Date(date_inf)
    date_sup <- as.Date(date_sup)
    lookback_monthlydata <- months(lookback_monthlydata)
    lookback_yearlydata <- years(lookback_yearlydata)

    month_seq <-
      seq(date_inf %m-% lookback_monthlydata, date_sup, by = 'months')




   #
   # Replacing all na values by a fix value
   #
    na_value = -100;

    replace_all_na <- function(x,na_value){
      if (is.numeric(x))  x[is.na(x)] = na_value
      return(x)
    }
    data[] <- data %>%
      lapply(replace_all_na,na_value = na_value)


    # Select monthly data
    data_monthly <- data %>%
      select(siret,
             periode,
             effectif,
             log_cotisationdue_effectif,
             log_ratio_dettecumulee_cotisation_12m)

    # Expand monthly data
    data_monthly <- data_monthly %>%
      filter(periode >=  date_inf %m-% lookback_monthlydata & periode <= date_sup) %>%
      tidyr::complete(periode = month_seq,
                      siret,
                      fill = list(
                        effectif = -100,
                        log_cotisationdue_effectif = -100,
                        log_ratio_dettecumulee_cotisation_12m = -100
                      )) %>%
      group_by(siret) %>%
      arrange(siret,periode) %>%
      melt(id.vars = c("siret", "periode")) %>%
      dcast(siret ~ periode)

  }
