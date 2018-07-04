compare_prediction_probabilities <-  function(sample_start,sample_end){

    assertthat::assert_that('prob' %in% names(sample_start) && 'prob' %in% names(sample_end))
    assertthat::assert_that(!anyDuplicated(sample_start$siret) && !anyDuplicated(sample_end$siret))

    sample_start <- sample_start %>% rename(prob_old = prob)

    joined_samples <- sample_start %>%
      full_join(sample_end,by = c('siret','periode')) %>%
      mutate(diff = prob - prob_old)

    return(joined_samples)
}
