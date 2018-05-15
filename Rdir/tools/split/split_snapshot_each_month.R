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


split_snapshot_each_month <-
  function(wholesample,
           date_inf,
           date_sup,
           frac_train = 0.60,
           frac_cross = 0.20,
           frac_eyeball = 0.05) {

    assert_that(
      frac_train > 0 ,
      frac_cross > 0,
      frac_eyeball >= 0,
      frac_train < 1,
      frac_train + frac_cross + frac_eyeball <= 1,
      msg = "Fractions must be positive and not exceed 1"
    )

    wholesample <- wholesample %>%
      arrange(periode) %>%
      select(siret,periode,date_defaillance,outcome_any)
    # donnees desequilibrees et donnees surechantillonnees

    imbalanced_sample <- wholesample %>%
      filter_time(~date_sup)

    default_oversample <- wholesample %>%
      filter(outcome_any) %>%
      filter_time(date_inf ~ date_sup) %>%
      filter(elapsed_months(periode,date_defaillance) <= 0,
            elapsed_months(periode,date_defaillance) > -12)

    oversampled_sample <- wholesample %>%
      filter(!outcome_any) %>%
      filter_time( ~ date_sup) %>%
      bind_rows(default_oversample) %>%
      as_tbl_time(periode)

    #


    sample_sirets_train <- oversampled_sample %>%
      as_tibble() %>%
      group_by(siret) %>%
      summarize() %>%
      sample_frac(frac_train)


    sample_train <- oversampled_sample %>%
      inner_join(sample_sirets_train,by = "siret") %>%
      select(siret, periode)

    cv_folds <- groupKFold(sample_train$siret, k = 5)

    remaining <- imbalanced_sample %>%
      anti_join(sample_sirets_train,by = "siret") %>%
      select(siret, periode)

    if (frac_eyeball > 0) {
      sample_eyeball <- remaining %>%
        sample_frac(frac_eyeball / (1 - frac_train - frac_cross))
    } else {
      sample_eyeball <-
        tibble(siret = character(), periode = as.Date(character()))
    }

    sample_test <- remaining %>%
      anti_join(sample_eyeball, by = 'siret')

    return(list("train" =  sample_train,"cv_folds" = cv_folds, "eyeball" = sample_eyeball,  "test" = sample_test))
  }
