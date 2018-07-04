predict_at <- function(actual_period, whole_data, prob_fun, rolling_mean = 3, clear_after = 2) {


  sample_actual <- whole_data %>%
    arrange(periode) %>%
    filter_time(actual_period %m-% months(rolling_mean-1) ~ actual_period)

  assertthat::assert_that(n_distinct(sample_actual$periode) == rolling_mean)
  assertthat::assert_that(any(sample_actual$periode == actual_period))


  prob <- prob_fun(sample_actual)

  prediction <- sample_actual %>%
      cbind(prob = prob) %>%
      group_by(siret) %>%
      summarize(prob = mean(prob), periode = max(periode)) %>%
      # Companies that are absent of the base for more than 'clear_after' months are cleared.
      filter(elapsed_months(start_date = periode, end_date  = actual_period) < clear_after) %>%
      arrange(desc(prob)) %>%
      as.data.frame()

  assertthat::assert_that('prob' %in% names(prediction) )
  assertthat::assert_that(!anyDuplicated(prediction$siret))

  return(prediction)
}
