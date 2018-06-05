generator <-
  function(data,
           lookback_yearlydata,
           lookback_monthlydata,
           delay,
           sirets,
           date_inf,
           date_sup,
           seed = 1010,
           batch_size = 128,
           step = 1) {

    library(reshape2)

    date_inf <- as.Date(date_inf)
    date_sup <- as.Date(date_sup)
    lookback_monthlydata_m <- months(lookback_monthlydata)
    lookback_yearlydata_y <- years(lookback_yearlydata)
    delay <- months(delay)

    # month_seq <-
    #   seq(date_inf %m-% lookback_monthlydata, date_sup, by = 'months')

    # Samples
    batch_sample <- data %>%
      select(siret,periode,outcome_any) %>%
      filter(siret %in% sirets) %>%
      arrange(periode) %>%
      filter_time( date_inf ~ date_sup) %>%
      mutate( prob = if_else(outcome_any,(n()-sum(outcome_any)) / n(),sum(outcome_any)/n())) %>%
      sample_n(size = batch_size,replace = TRUE, weight = .$prob)






    #
    # Replacing all na values by a fix value
    #
    na_value = -100;

    replace_all_na <- function(x,na_value){
      if (is.numeric(x))  x[is.na(x)] = na_value
      return(x)
    }
    batch_sample[] <- batch_sample %>%
      lapply(replace_all_na,na_value = na_value)




    # Select monthly data

    num_vars = 3

    sample_monthly <- array(0,dim = c(batch_size,lookback_monthlydata,num_vars))
    target_monthly <- array(0,dim = c(batch_size))

    for (ii in 1:batch_size){
      auxSiret = batch_sample$siret[ii]
      auxPeriod = batch_sample$periode[ii]


      # encapsuler cette transformation dans une fontion déportée FIX ME.
      data_monthly <- data %>%
        filter(siret == auxSiret) %>%
        filter_time(auxPeriod %m-% lookback_monthlydata_m %m+% months(1) ~ auxPeriod) %>%
        select(siret,
               periode,
               effectif,
               log_cotisationdue_effectif,
               log_ratio_dettecumulee_cotisation_12m) %>%
        tidyr::complete(
          periode =  seq(
            auxPeriod %m-% lookback_monthlydata_m %m+% months(1),
            auxPeriod,
            by = 'months'),
          siret,
          fill = list(
            effectif = -100,
            log_cotisationdue_effectif = -100,
            log_ratio_dettecumulee_cotisation_12m = -100
          )
        ) %>%
        melt(id.vars = c("siret", "periode")) %>%
        dcast(siret + variable ~ periode) %>%
        # Ou plus simple, direct as.matrix !!
        # Potentiellement besoin de trier les périodes dans l'autre sens
        # FIX ME.
        select(-siret,-variable) %>%
        as.matrix() %>%
        t()

      sample_monthly[ii,,] <- data_monthly;
      target_monthly[[ii]] <- data %>%
        filter(siret == auxSiret, periode == auxPeriod %m+% delay) %>%
        .$outcome
    }
    return(list(sample_monthly,target_monthly))

  }
