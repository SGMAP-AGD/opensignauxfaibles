add_past_trends <- function(data,variables, lookback_months, slope = FALSE){

  assertthat::assert_that(all(c('siret','periode', variables) %in% names(data)))
  assertthat::assert_that(length(variables) == 1 | length(variables) == length(lookback_months))

  if (length(variables) == 1) {
    variables <- rep(variables, length(lookback_months))
  }
  data_nested <- data %>%
    arrange(siret,periode) %>%
    group_by(siret) %>%
    nest()

  if (slope) {
    aux_slope <- function(x,y){
      n <- length(x)
      sdx <- sqrt((n^2 - 1)/12)*sqrt(n/(n-1))
      slope <- cor(x,y) * sd(y)/sdx
      return(slope)
    }

    aux_past <- function(data, variable, last_n) {
      y <- data[[variable]]
      slopes <- array(dim = length(y))
      if (last_n <= length(y))
        for (i in last_n:length(y)) {
          sub_y <- tail(y[1:i], n = last_n)
          x <- 1:length(sub_y)
          slopes[i] <- aux_slope(x, sub_y)
        }

      if (length(slopes) != length(y))
        browser()
      assertthat::assert_that(length(slopes) == length(y))

      return(slopes)
    }
  } else {
    aux_past <-  function(data, variable, last_n) {
      y <- data[[variable]]
      y_lag <- lag(y, last_n)
      return(y / y_lag[1:length(y)] - 1)
    }
  }

  cat('/!\ FIX ME: How are slopes on 6 months computed with only  3 available months? + Sensitive to wholes','\n')

  for (i in seq_along(lookback_months)){
      data_nested[[paste0(variables[i],'_variation_',lb)]] <-  purrr::map(data_nested$data, aux_past, variable= variables[i], last_n = lookback_months[i])
  }


  data <- data_nested %>%
    unnest()

  return(data)

}
