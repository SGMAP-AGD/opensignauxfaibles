#' Is siren
#'
#' @param siren a string vector which is suspected to include siren number
#'
#' @return a boolean vector with value TRUE if the string includes exactly 9 digits.
#' @export
#'
#' @examples
#' is_siren("2015")
#' is_siren("201512125")

is_siren <- function(siren) {
  stringr::str_detect(string = siren, pattern = "^[[:digit:]]{9}$")
}

#' convert urssaf periods to standard format
#'
#' Urssaf store dates in a very specific format.
#' This can be years + 62 for annual values,
#' years + quarter + 0 for quarterly values
#' years + quarter + number of the month in the quarter for monthly values.
#'
#' @param .data a tibble with a period
#' @param .variable the name of the variable with urssaf periods
#' @param format can bee yyqm or yyyyqm
#'
#' @return an new tibble with two more columns : periodicity and period
#' @export
#'
#' @examples
#' convert_urssaf_periods_(.variable = ~ yyyyqm, format = "yyyyqm")
#'

convert_urssaf_periods_ <- function(.data, .variable, format = "yyqm") {
  pattern_ = "([[:digit:]]{2,4})([[:digit:]]{1})([[:digit:]]{1})"
  year_ = stringr::str_replace(
    string = lazyeval::f_eval(~ uq(.variable), data = .data),
    pattern = pattern_,
    replacement = "\\1") %>%
    as.numeric()
  if (format == "yyqm") {
    year_ <- ifelse(year_ <= 20, year_ + 2000, year_ + 1900)
  }
  quarter_ <- stringr::str_replace(
    string = lazyeval::f_eval(~ uq(.variable), data = .data),
    pattern = pattern_,
    replacement = "\\2"
  )
  month_ <- stringr::str_replace(
    string = lazyeval::f_eval(~ uq(.variable), data = .data),
    pattern = pattern_,
    replacement = "\\3"
  )
  periodicity_ <- ifelse(
    stringr::str_detect(
      string = lazyeval::f_eval(~ uq(.variable), data = .data),
      pattern = "[[:digit:]]{2,4}62$") == TRUE,
    "yearly",
    ifelse(
      month_ == "0",
      "quarterly",
      "monthly"
    )
  )
  period_ <- ifelse(
    periodicity_ == "yearly",
    as.character(year_),
    ifelse(
      periodicity_ == "quarterly",
      paste0(year_, "-Q", quarter_),
      paste0(year_, "-",
             stringr::str_pad(
               string = (as.numeric(quarter_) - 1) * 3 + as.numeric(month_),
               side = "left",
               width = 2,
               pad = "0")
      )
    )
  )
  return(dplyr::bind_cols(.data, tibble::tibble(periodicity = periodicity_, period = period_)))
}

