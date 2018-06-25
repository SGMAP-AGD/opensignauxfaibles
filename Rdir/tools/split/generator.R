generator <-
  function(data,
           lookback_yearlydata,
           lookback_monthlydata,
           sirets,
           date_inf,
           date_sup,
           seed = 1010,
           batch_size = 128,
           step = 1,
           na_value = -5
           ) {

    set.seed(seed)

    date_inf <- as.Date(date_inf)
    date_sup <- as.Date(date_sup)
    lookback_monthlydata_m <- months(lookback_monthlydata)
    lookback_yearlydata_m <- months(lookback_yearlydata * 12)


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

    sub_data <- data %>%
      filter(siret %in% sirets) %>%
      arrange(periode)

    possible_pairs <- sub_data %>%
      select(siret,periode,outcome) %>%
      filter_time( date_inf ~ date_sup) %>%
      mutate( prob = if_else(outcome=='default',
                             (n()-sum(outcome=='default')) / n(),
                             sum(outcome == 'default')/n()))

    #
    # Replacing all na values by a fix value
    #

    # replace_all_na <- function(x,na_value){
    #   if (is.numeric(x))  x[is.na(x)] = na_value
    #   return(x)
    # }
    # sub_data[] <- sub_data %>%
    #   lapply(replace_all_na,na_value = na_value)





    gen <- function(){
      #
      # Sampling the input data by siret and period
      #

      batch_sample <- possible_pairs %>%
        sample_n(size = batch_size,replace = TRUE, weight = .$prob)



      return(
        df_to_RNN_input(
          siret_period_pairs = batch_sample,
          df= sub_data,
          date_inf = date_inf,
          date_sup = date_sup,
          lookback_monthlydata = lookback_monthlydata,
          lookback_yearlydata = lookback_yearlydata,
          na_value = na_value
        )
      )

    }
    return(gen)
  }
