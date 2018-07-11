df_to_RNN_input <- function(siret_period_pairs,df_m,df_y,date_inf,date_sup,lookback_monthlydata, lookback_yearlydata, na_value){
  batch_size = nrow(siret_period_pairs)
  date_inf <- as.Date(date_inf)
  date_sup <- as.Date(date_sup)
  lookback_monthlydata_m <- months(lookback_monthlydata)
  lookback_yearlydata_m <- months(lookback_yearlydata * 12)

  n <- nrow(siret_period_pairs)

  ##
  ###
  # Define variables of interest
  ###
  ##

  var0 <-  c('siret','periode')

  var_monthly <-  c('effectif',
                    'log_cotisationdue_effectif',
                    'ratio_dettecumulee_cotisation_12m')

  var_yearly <- c('poids_frng',
                  'taux_marge',
                  'delai_fournisseur',
                  'dette_fiscale',
                  'financier_court_terme',
                  'frais_financier')

  # check_input
  assertthat::assert_that('outcome' %in% names(siret_period_pairs))
  assertthat::assert_that(all(var0 %in% names(siret_period_pairs)) & all(var0 %in% names(df_m)))
  assertthat::assert_that(all(var_monthly %in% names(df_m)))
  assertthat::assert_that(all(var_yearly %in% names(df_m)))

  # Replace missing values
  replace_all_na <- function(x,na_value){
    if (is.numeric(x))  x[is.na(x)] = na_value
    return(x)
  }
  df_m[] <- df_m %>%
    lapply(FUN = replace_all_na, na_value = na_value)

  ##
  ###
  # define RNN target
  ###
  ##

  target <- if_else(siret_period_pairs$outcome == 'default',1,0)

  # Prep

  siret_period_pairs <- siret_period_pairs %>%
    select(-outcome,-prob) %>%
    mutate(id = 1:batch_size)

  ##
  ###
  # sample_monthly  ##############################################
  ###
  ##

  var_monthly <- c(var0,var_monthly)

  sample_monthly <- siret_period_pairs %>%
    mutate(periode_lookback = periode %m-% lookback_monthlydata_m %m+% months(1)) %>%
    rename(periode_obs = periode) %>%
    left_join(df_m %>% select(var_monthly), by = 'siret') %>%
    toTimeSeries(batch_size,
                 periodicity = 'months',
                 na_value = na_value)

  ##
  ###
  # sample_yearly ################################################
  ###
  ##

  var_yearly <- c(var0, var_yearly)

  sample_yearly <- siret_period_pairs %>%
    mutate(periode_lookback = floor_index(periode,'1 Y') %m-% lookback_yearlydata_m %m+% months(12)) %>%
    mutate(periode_obs = floor_index(periode,'1 Y')) %>%
    select(-periode) %>%
    left_join(df_y %>% select(var_yearly),by = 'siret') %>%
    toTimeSeries(batch_size,
                 periodicity = 'years',
                 na_value = na_value)

  # df_yearly <- siret_period_pairs %>%
  #   mutate(periode_lookback = floor_index(periode,'1 Y') %m-% lookback_yearlydata_m %m+% months(12)) %>%
  #   mutate(periode_obs = floor_index(periode,'1 Y')) %>%
  #   select(-periode) %>%
  #   left_join(df_m %>% select(var_yearly),by = 'siret') %>%
  #   arrange(periode) %>%
  #   as_tbl_time(periode) %>%
  #   collapse_by("1 year",
  #               side = 'start',
  #               start_date = floor_index(min(.$periode),'1 Y'),
  #               clean = TRUE) %>%
  #   group_by(id, periode) %>%
  #   summarize_all(function(x) last(x[!is.na(x)])) %>%
  #   ungroup()

  # sample_yearly <- df_yearly %>%
  #   toTimeSeries(batch_size,
  #                periodicity = 'years',
  #                na_value = na_value)


  #   for (ii in 1:n){
  #     aux_siret = siret_period_pairs$siret[ii]
  #     aux_period = floor_index(siret_period_pairs$periode[ii],'1 year')
  #
  #     # tic("prep")
  #      aux <- df %>%
  #       filter(siret == aux_siret) %>%
  #       select(vars) %>%
  #       collapse_by("1 year", side = 'start', start_date = floor_index(min(.$periode),'1 Y'),clean = TRUE) %>%
  #       group_by(periode) %>%
  #       summarize_all(function(x)last(x))
  #      # toc()
  #      #tic("TTS")
  #     sample_yearly[ii,,] <- aux %>%
  #       toTimeSeries(
  #         aux_period %m-% lookback_yearlydata_m %m+% months(12),
  #         aux_period,
  #         periodicity = 'years',
  #         na_value = na_value
  #       )
  #     #toc()
  #   }
  #
  #
  ##
  ###
  # Select other data ###############################################################
  ###
  ##

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

    return(list(
      list(
        monthly_input = sample_monthly,
        yearly_input = sample_yearly,
        other_input = sample_other
      ),
      list(target)
    ))


}

##
###
####
#### Auxiliary function
####
###
##

toTimeSeries <- function(data,batch_size, periodicity,na_value){
  # periodicity = 'years' or 'months'


  if (typeof(na_value) == 'list'){
    na_replace_list = na_value
  } else {
    fields <- names(data %>% select(-siret, -periode,-periode_lookback,-periode_obs))
    na_replace_list <-  as.list(array(na_value,dim = c(length(fields),1)))
    names(na_replace_list) <- fields
  }

  assertthat::assert_that(all(c('siret','periode') %in% names(data)))

  time_series <- data %>%
    filter(periode >= periode_lookback & periode <= periode_obs) %>%
    group_by(id) %>%
    tidyr::complete(periode = seq.Date(first(periode_lookback),
                                       first(periode_obs),
                                       by = periodicity),
                    siret,
                    fill = na_replace_list) %>%
    ungroup() %>%
    select(-siret, -periode, -id, -periode_obs, -periode_lookback) %>%
    as.matrix() %>%
    array_reshape(dim = c(batch_size, nrow(.) / batch_size, ncol(.)))

  return(time_series)
}


# toTimeSeries <- function(SP_pairs,data,periodicity,na_value){
#   batch_size <-  nrow(SP_pairs)
#   # periodicity = 'years' or 'months'
#   if (typeof(na_value) == 'list'){
#     na_replace_list = na_value
#   } else {
#     fields <- names(data %>% select(-siret, -periode))
#     na_replace_list <-  as.list(array(na_value,dim = c(length(fields),1)))
#     names(na_replace_list) <- fields
#   }
#
#   assertthat::assert_that(all(c('siret','periode') %in% names(data)))
#
#   time_series <- SP_pairs %>%
#     left_join(data, by = 'siret') %>%
#     filter(periode >= periode_lookback & periode <= periode_obs) %>%
#     group_by(siret, periode_obs) %>%
#     tidyr::complete(periode = seq.Date(first(periode_lookback),
#                                        first(periode_obs),
#                                        by = periodicity),
#                     siret,
#                     fill = na_replace_list) %>%
#     ungroup() %>%
#     select(-siret, -periode, -periode_obs, -periode_lookback) %>%
#     as.matrix() %>%
#     array_reshape(dim = c(batch_size, nrow(.) / batch_size, ncol(.)))
#
#   return(time_series)
# }
