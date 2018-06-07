generator <-
  function(data,
           lookback_yearlydata,
           lookback_monthlydata,
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
    lookback_yearlydata_m <- months(lookback_yearlydata * 12)
    # delay: le délai est fixé dans la fonction objectif !!

    #
    # Sampling the input data by siret and period
    #
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



    #
    # Select monthly data ###############################################################
    #

    vars <- c('siret','periode')
    vars <- c(vars,
              'effectif',
              'log_cotisationdue_effectif',
              'log_ratio_dettecumulee_cotisation_12m')


    sample_monthly <- array(0,dim = c(batch_size,lookback_monthlydata,length(vars)-2))
    target <- array(0,dim = c(batch_size))

    for (ii in 1:batch_size){
      aux_siret = batch_sample$siret[ii]
      aux_period = batch_sample$periode[ii]

      sample_monthly[ii,,] <- data %>%
        filter(siret == aux_siret) %>%
        select(vars) %>%
        toTimeSeries(
          aux_period %m-% lookback_monthlydata_m %m+% months(1),
          aux_period,
          periodicity = 'months',
          na_value = -100
        )

      target[[ii]] <- data %>%
        filter(siret == aux_siret, periode == aux_period) %>%
        select(outcome)
    }

    #
    # Select yearly data ###############################################################
    #

    vars <- c('siret','periode')
    vars <- c(vars,
              'poids_frng',
              'taux_marge',
              'delai_fournisseur',
              'dette_fiscale',
              'financier_court_terme',
              'frais_financier')

    sample_yearly <- array(0,dim = c(batch_size,lookback_yearlydata,length(vars)-2))


    for (ii in 1:batch_size){
      aux_siret = batch_sample$siret[ii]
      aux_period = floor_index(batch_sample$periode[ii],'1 year')

      sample_yearly[ii,,] <- data %>%
        filter(siret == aux_siret) %>%
        select(vars) %>%
       collapse_by("1 year", side = 'start', start_date = floor_index(min(.$periode),'1 Y'),clean = TRUE) %>%
       group_by(periode) %>%
       summarize_all(first) %>%
        toTimeSeries(
          aux_period %m-% lookback_yearlydata_m %m+% months(12),
          aux_period,
          periodicity = 'years',
          na_value = -100
        )
    }

    #
    # Select other data ###############################################################
    #
    vars <- c('siret','periode')
    vars <- c(vars,
              'code_ape_1',
              'code_ape_2',
              'code_ape_3',
              'code_ape_4',
              'code_ape_5'
    )

    sample_other <- array(0,dim = c(batch_size,length(vars)-2))

    for (ii in 1:batch_size){

      aux_siret = batch_sample$siret[ii]

      sample_other[ii,] <- data %>%
        filter(siret == aux_siret) %>%
        mutate(#code_ape_1 = FIX ME,
               code_ape_2 = str_sub(code_ape,1,2),
               code_ape_3 = str_sub(code_ape,1,3),
               code_ape_4 = str_sub(code_ape,1,4),
               code_ape_5 = str_sub(code_ape,1,5)) %>%
        select(vars) %>%
        summarize_all(first) %>%
        select(-siret, -periode) %>%
        as.matrix()
    }

    return(list(
      list(sample_monthly = sample_monthly,
           sample_yearly = sample_yearly,
           sample_other),
      list( target = target)
      )
    )

  }


toTimeSeries <- function(data,low_date,sup_date,periodicity,na_value){
# periodicity = 'years' or 'months'
  if (typeof(na_value) == 'list'){
    na_replace_list = na_value
  } else {
    fields <- names(data %>% select(-siret, -periode))
    na_replace_list <-  as.list(array(na_value,dim = c(length(fields),1)))
    names(na_replace_list) <- fields
  }

  assertthat::assert_that(all(c('siret','periode') %in% names(data)))

  time_series <- data %>%
    filter_time(low_date ~ sup_date) %>%
    tidyr::complete(
      periode = seq(
        low_date,
        sup_date,
        by = periodicity),
      siret,
      fill = na_replace_list
    ) %>%
    select(-siret,-periode) %>%
    as.matrix()


  return(time_series)
}
