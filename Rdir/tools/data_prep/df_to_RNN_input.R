df_to_RNN_input <- function(siret_period_pairs,df_m,df_y,var_monthly,lookback_monthly_data, var_yearly,lookback_yearly_data, var_other,na_value, verbose = FALSE){
  if (verbose){
    cat('Préparation du format pour le RNN','\n')
  }
  batch_size <-  nrow(siret_period_pairs)
  lookback_monthlydata_m <- months(lookback_monthly_data)
  lookback_yearlydata_m <- months(lookback_yearly_data * 12)

  #n <- nrow(siret_period_pairs)

  ##
  ###
  ######################
  # Prepare variables ##
  ######################
  ###
  ##

  # check_input
  assertthat::assert_that('outcome' %in% names(siret_period_pairs))
  assertthat::assert_that(all(c('siret','periode') %in% names(siret_period_pairs)) & all(c('siret','periode')  %in% names(df_m)))
  assertthat::assert_that(all(var_monthly %in% names(df_m)))
  assertthat::assert_that(all(var_yearly %in% names(df_m)))

  assertthat::assert_that(nrow(df_m) == n_distinct(df_m %>% select(siret,periode)))


  ##
  ###
  ######################
  # define RNN target ##
  ######################
  ###
  ##

  target <- if_else(siret_period_pairs$outcome == 'default',1,0)

  # Prep

  siret_period_pairs <- siret_period_pairs %>%
    select(siret, periode) %>%
    mutate(id = 1:batch_size)

  ##
  ###
  ####################
  # sample_monthly  ##
  ####################
  ###
  ##
  if (verbose){
    #tic('Monthly')
    cat('1/4: Input monthly ...','\n')
  }
  var_monthly <- c('siret','periode',var_monthly)

  sample_monthly <- siret_period_pairs %>%
    mutate(periode_lookback = periode %m-% lookback_monthlydata_m %m+% months(1)) %>%
    rename(periode_obs = periode) %>%
    left_join(df_m %>% select(var_monthly), by = 'siret') %>%
    toTimeSeries(batch_size,
                 periodicity = 'months',
                 na_value = na_value)
  if (verbose) {
    #toc()
    cat('done','\n')
  }
  ##
  ###
  ##################
  # sample_yearly ##
  ##################
  ###
  ##

  if (verbose){
    #tic('yearly')
    cat('2/4: Input yearly ...','\n')
  }

  var_yearly <- c('siret','periode', var_yearly)

  sample_yearly <- siret_period_pairs %>%
    mutate(periode_lookback = floor_index(periode,'1 Y') %m-% lookback_yearlydata_m %m+% months(12)) %>%
    mutate(periode_obs = floor_index(periode,'1 Y')) %>%
    select(-periode) %>%
    left_join(df_y %>% select(var_yearly),by = 'siret') %>%
    toTimeSeries(batch_size,
                 periodicity = 'years',
                 na_value = na_value)

  if (verbose) {
    #toc()
    cat('done','\n')
  }

  ##
  ###
  ####################
  # Select ape data ##
  ####################
  ###
  ##

  if (verbose){
    cat('3/4: Input APE ...','\n')
  }

  vars <- c('siret','periode','code_ape')

  sample_ape <- array(0,dim = c(batch_size,length(vars)-2))

  sample_ape <- siret_period_pairs %>%
    left_join(df_m,by=c('siret','periode')) %>%
    mutate(code_ape = as.numeric(code_ape)-1) %>%
    select(vars) %>%
    select(-siret, -periode) %>%
    as.matrix()

  ##
  ###
  ######################
  # Select other data ##
  ######################
  ###
  ##

  if (verbose){
    cat('4/4: Input other ...','\n')
  }

  vars <- c('siret','periode',    var_other)

  sample_other <- array(0,dim = c(batch_size,length(vars)-2))

  sample_other <- siret_period_pairs %>%
    left_join(df_m,by=c('siret','periode')) %>%
    select(vars) %>%
    select(-siret, -periode) %>%
    as.matrix()

  ##
  ###
  ###
  ##

  dimnames(sample_monthly) <- NULL
  dimnames(sample_yearly) <- NULL
  dimnames(sample_ape) <- NULL
  dimnames(sample_other) <- NULL
  dimnames(target) <- NULL

  return(list(
    list(
      monthly_input = sample_monthly,
      yearly_input = sample_yearly,
      ape_input = sample_ape,
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

toTimeSeries <- function(myData,batch_size, periodicity,na_value){
  # periodicity = 'years' or 'months'

  # if (typeof(na_value) == 'list'){
  #   na_replace_list = na_value
  # } else {
  #   fields <- names(myData %>% select(-siret, -periode,-periode_lookback,-periode_obs))
  #   na_replace_list <-  as.list(array(na_value,dim = c(length(fields),1)))
  #   names(na_replace_list) <- fields
  # }
  #
  # assertthat::assert_that(all(c('siret','periode') %in% names(myData)))

  time_series <- myData %>%
    arrange(id)

  siret_period_pairs <- time_series[!diff(c(0,time_series$id))==0,] %>%
    select(siret,periode_lookback,periode_obs,id)

  period_list_aux <- siret_period_pairs %>%
    select(periode_lookback, periode_obs) %>%
    plyr::mlply(.fun = function(periode_lookback,periode_obs)
                                          seq.Date(periode_lookback,periode_obs,by=periodicity)
                                 )
  time_series <-  siret_period_pairs %>%
    mutate(period_list = period_list_aux) %>%
    unnest(periode =period_list) %>%
    left_join(myData, by = c('siret', 'periode','periode_lookback','periode_obs','id'))


  time_series <- time_series %>%
    select(-siret, -periode, -id, -periode_obs, -periode_lookback) %>%
    as.matrix() %>%
    array_reshape(dim = c(batch_size, nrow(.) / batch_size, ncol(.)))

  ### Replace all na
  time_series[is.na(time_series)] <- na_value
  return(time_series)
}
