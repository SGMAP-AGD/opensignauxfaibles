generator <-
  function(data,
           var_yearly,
           lookback_yearly_data,
           var_monthly,
           lookback_monthly_data,
           var_other,
           sirets,
           date_inf,
           date_sup,
           seed = 1010,
           batch_size = 128,
           step = 1,
           na_value = 0,
           oversampling = TRUE) {
    set.seed(seed)

    date_inf <- as.Date(date_inf)
    date_sup <- as.Date(date_sup)
    lookback_monthlydata_m <- months(lookback_monthly_data)
    lookback_yearlydata_m <- months(lookback_yearly_data * 12)


    var0 <-  c('siret', 'periode')



    df_m <- data %>%
      filter(siret %in% sirets) %>%
      arrange(periode)

    possible_pairs <- df_m %>%
      select(siret, periode, outcome) %>%
      filter_time(date_inf ~ date_sup)

    if (oversampling) {
      possible_pairs <- possible_pairs %>%
        mutate(prob = if_else(
          outcome == 'default',
          (n() - sum(outcome == 'default')) / n(),
          sum(outcome == 'default') / n()
        ))
    } else {
      possible_pairs <- possible_pairs %>%
        mutate(prob = 1)
    }

    #
    # Replacing all na values by a fix value
    #

    df_y <- monthly_to_yearly(df_m)

    gen <- function() {
      #
      # Sampling the input data by siret and period
      #

      batch_sample <- possible_pairs %>%
        sample_n(size = batch_size,
                 replace = FALSE,
                 weight = .$prob)

      out <- df_to_RNN_input(
        siret_period_pairs = batch_sample,
        df_m = df_m,
        df_y = df_y,
        var_monthly = var_monthly,
        lookback_monthly_data = lookback_monthly_data,
        var_yearly = var_yearly,
        lookback_yearly_data = lookback_yearly_data,
        var_other = var_other,
        na_value = na_value
      )

      return(out)

    }
    return(gen)
  }
