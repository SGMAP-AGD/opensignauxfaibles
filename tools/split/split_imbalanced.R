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
#'
split_imbalanced <-
  function(wholesample,
           date_inf,
           date_sup,
           frac_train = 0.60,
           frac_cross = 0.20) {
    library(assertthat)
    library(dplyr)

    assert_that(
      frac_train > 0 ,
      frac_train < 1,
      frac_cross > 0,
      frac_cross < 1,
      frac_train + frac_cross <= 1,
      msg = "Fractions must be positive and strictly not exceed 1"
    )

    # donnees desequilibrees et donnees surechantillonnees

    imbalanced_subsample <- wholesample %>%
      filter(periode == date_sup)

    #


    sample_sirets_train <- imbalanced_subsample %>%
      group_by(siret) %>%
      summarize(a = first(siret)) %>%
      select(-a) %>%
      sample_frac(frac_train)


    sample_train <- imbalanced_subsample %>%
      inner_join(sample_sirets_train,by = "siret")

    cv_folds <- groupKFold(sample_train$siret, k = 5)

    sample_test <- imbalanced_subsample %>%
      anti_join(sample_sirets_train,by = "siret")

    return(list("train" =  sample_train,"cv_folds" = cv_folds, "test" = sample_test))
  }
