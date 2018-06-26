df_to_RNN_input <- function(siret_period_pairs,df,date_inf,date_sup,lookback_monthlydata, lookback_yearlydata, na_value){

  date_inf <- as.Date(date_inf)
  date_sup <- as.Date(date_sup)
  lookback_monthlydata_m <- months(lookback_monthlydata)
  lookback_yearlydata_m <- months(lookback_yearlydata * 12)

  n <- nrow(siret_period_pairs)
  # Define variables of interest
  var0 <-  c('siret','periode')

  var_monthly <-  c('effectif',
                    'log_cotisationdue_effectif',
                    'log_ratio_dettecumulee_cotisation_12m')

  var_yearly <- c('poids_frng',
                  'taux_marge',
                  'delai_fournisseur',
                  'dette_fiscale',
                  'financier_court_terme',
                  'frais_financier')

  # check_input
  assertthat::assert_that('outcome' %in% names(siret_period_pairs))
  assertthat::assert_that(all(var0 %in% names(siret_period_pairs)) & all(var0 %in% names(df)))
  assertthat::assert_that(all(var_monthly %in% names(df)))
  assertthat::assert_that(all(var_yearly %in% names(df)))

  # Replace missing values
  replace_all_na <- function(x,na_value){
    if (is.numeric(x))  x[is.na(x)] = na_value
    return(x)
  }
  df[] <- df %>%
    lapply(FUN = replace_all_na, na_value = na_value)

  # define RNN target

  target <- if_else(siret_period_pairs$outcome == 'default',1,0)

  ###
  # sample_monthly  ##############################################
  ###

  vars <- c(var0,var_monthly)

  sample_monthly <- array(0,dim = c(n,lookback_monthlydata,length(vars)-2))


  for (ii in 1:n){

    aux_siret <-  siret_period_pairs$siret[ii]
    aux_period <- siret_period_pairs$periode[ii]

    sample_monthly[ii,,] <- df %>%
      filter(siret == aux_siret) %>%
      select(vars) %>%
      toTimeSeries(
        aux_period %m-% lookback_monthlydata_m %m+% months(1),
        aux_period,
        periodicity = 'months',
        na_value = na_value
      )
  }


  #
  # Select yearly data ###############################################################
  #

  vars <- c(var0, var_yearly)

  sample_yearly <- array(0,dim = c(n,lookback_yearlydata,length(vars)-2))


  for (ii in 1:n){
    aux_siret = siret_period_pairs$siret[ii]
    aux_period = floor_index(siret_period_pairs$periode[ii],'1 year')

    sample_yearly[ii,,] <- df %>%
      filter(siret == aux_siret) %>%
      select(vars) %>%
      collapse_by("1 year", side = 'start', start_date = floor_index(min(.$periode),'1 Y'),clean = TRUE) %>%
      group_by(periode) %>%
      summarize_all(function(x)last(x)) %>%
      toTimeSeries(
        aux_period %m-% lookback_yearlydata_m %m+% months(12),
        aux_period,
        periodicity = 'years',
        na_value = na_value
      )
  }


  #
  # Select other data ###############################################################
  #
  vars <- c('siret','periode')
  vars <- c(vars,
            'code_ape')#,
  #           'code_ape_2',
  #           'code_ape_3',
  #           'code_ape_4',
  #           'code_ape_5'
  # )

  sample_other <- array(0,dim = c(n,length(vars)-2))

  for (ii in 1:n){
    aux_siret = siret_period_pairs$siret[ii]

    sample_other[ii,] <- df %>%
      mutate(code_ape = as.numeric(code_ape)-1) %>%
      filter(siret == aux_siret) %>%
      select(vars) %>%
      summarize_all(first) %>%
      select(-siret, -periode) %>%
      as.matrix()
  }

  dimnames(sample_monthly) <- NULL
  dimnames(sample_yearly) <- NULL
  dimnames(sample_other) <- NULL
  dimnames(target) <- NULL


  # ## Provisoire
  # return(list(
  #   sample_monthly,
  #   target
  # ))
  #############
  return(list(
    list(
      monthly_input = sample_monthly,
      yearly_input = sample_yearly,
      other_input = sample_other
    ),
    list(target)
  ))


  }


####
#### Auxiliary function
####

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

