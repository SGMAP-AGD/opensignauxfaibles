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
           crossvalidation = TRUE,
           frac_train = 0.60,
           frac_val = 0.20,
           frac_eyeball = 0.05,
           seed = 1010) {



    assertthat::assert_that(
      frac_train > 0 ,
      frac_val > 0,
      frac_eyeball >= 0,
      frac_train < 1,
      frac_train + frac_val + frac_eyeball <= 1,
      msg = "Fractions must be positive and not exceed 1"
    )
    assertthat::assert_that(
      nrow(data) == n_distinct(data %>% select(siret, periode))
    )


    set.seed(seed)

    date_inf <- as.Date(date_inf)
    date_sup <- as.Date(date_sup)

    data <- data %>%
      arrange(periode) %>%
      filter_time( date_inf ~ date_sup) %>%
      select(siret,periode,outcome)

    sample_sirets_train <- sample_frac(tbl = unique(data['siret']), size = frac_train + frac_val)

    sample_train <- data %>%
      semi_join(sample_sirets_train,by = "siret")

    cat("Fraction of positive outcomes in sample_train:", sum(sample_train$outcome == 'default') / nrow(sample_train))
    cat('\n')
    cat("Fraction of sirets in sample_train:", n_distinct(sample_train$siret) / n_distinct(data$siret))
    cat('\n')

    sample_train <- sample_train %>%
      select(siret, periode)


    remaining <- data %>%
      anti_join(sample_sirets_train,by = "siret") %>%
      select(siret, periode)

    cv_folds <- list()
    sample_val <- data.frame()

    if (crossvalidation){
      cv_folds <- groupKFold(sample_train$siret, k = 5)
    } else {
      val_siret <- sample_frac(tbl = unique(sample_train['siret']), size = frac_val / (frac_val + frac_train))
      sample_val <- sample_train %>% filter(siret %in% val_siret$siret)
      sample_train <- sample_train %>% filter(! siret %in% val_siret$siret)
    }


    if (frac_eyeball > 0) {
      sample_eyeball <- remaining %>%
        sample_frac(frac_eyeball / (1 - frac_train - frac_val))
    } else {
      sample_eyeball <-
        tibble(siret = character(), periode = as.Date(character()))
    }

    sample_test <- remaining %>%
      anti_join(sample_eyeball, by = 'siret')

    return(list("train" =  sample_train,"cv_folds" = cv_folds,"validation" = sample_val, "eyeball" = sample_eyeball,  "test" = sample_test))
  }
