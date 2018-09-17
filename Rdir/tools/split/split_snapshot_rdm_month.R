#' Splits the wholesample into a train, crossvalidation and test samples.
#' Oversamples positive outcomes in the training set by making snapshots at different periods.
#' Only considers negative outcomes at time \code{date_sup}
#'
#'
#' @param wholesample A dataframe, with fields \code{periode} and \code{outcome}.
#' @param date_inf A Date
#' @param date_sup A Date
#' @return Three disjoint dataframes: train, crossvalidation and test samples.
#' @examples
#'


split_snapshot_rdm_month <-
  function(data,
           date_inf,
           date_sup,
           frac_train = 0.60,
           frac_val = 0.20,
           frac_eyeball = 0) {

    assertthat::assert_that(
      frac_train > 0 ,
      frac_val > 0,
      frac_train < 1,
      frac_train + frac_val <= 1,
      msg = "Fractions must be positive and not exceed 1"
    )

    frac_test = 1 - (frac_train + frac_val)


    assertthat::assert_that(
      nrow(data) == n_distinct(data %>% select(siret, periode))
    )

    # For reproducibility ###
    raw_data <- raw_data %>%
      arrange(siret,periode)
    #########################


    date_inf <- as.Date(date_inf)
    date_sup <- as.Date(date_sup)

    data <- data %>%
      arrange(periode) %>%
      filter_time( date_inf ~ date_sup) %>%
      select(siret,periode)

    ## TRAIN
    sirets  <- unique(data['siret'])

    sirets_groups <- split(
      sirets,
      sample(1:3, nrow(sirets), prob = c(frac_train,frac_val, frac_test), replace=T)
    )

    partition_names <- c('train','validation','test')
    non_empty <- c(frac_train >0, frac_val > 0, frac_test > 0)

    result <- lapply(1:sum(non_empty), function(x){
      cat("Fraction of sirets (", partition_names[non_empty][x], "):", nrow(sirets_groups[[x]]) / n_distinct(data$siret))
      cat('\n')
      data %>% semi_join(sirets_groups[[x]], by = "siret")})

    names(result) <- partition_names[non_empty]

    result[[1]] <- result[[1]] %>%
      groupdata2::fold(5, id_col = 'siret') %>%
      rename(fold_column = .folds) %>%
      mutate(fold_column = as.numeric(fold_column))



    return(result)
  }
