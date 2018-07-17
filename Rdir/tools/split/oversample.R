oversample <- function(df) {
  assertthat::assert_that(all(c('siret', 'periode', 'outcome') %in% names(df)))


  n_to_oversample <- abs(diff(table(df$outcome)))

  cat('Balance before oversampling:',
      table(df$outcome)['default'] / nrow(df),
      '\n')

  oversamples <- df %>%
    filter(outcome == 'default') %>%
    sample_n(n_to_oversample, replace = TRUE)

  result <- df %>% rbind(oversamples)

  cat('Balance after oversampling:',
      table(result$outcome)['default'] / nrow(result),
      '\n')

  return(result)

}
