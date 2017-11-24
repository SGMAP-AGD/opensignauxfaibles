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

#' Get table last n months
#'
#' @param .date a date as a string
#' @param .n_months an integer
#'
#' @return a tibble
#' @export
#'
#' @examples
#' get_table_last_n_months(.date = "2017-01-01", .n_months = 12)
#'
get_table_last_n_months <- function(.date, .n_months) {
  tibble::tibble(
    period = format(
      x = seq(
        from = lubridate::ymd(.date) %m-% months(.n_months),
        to = lubridate::ymd(.date) %m-% months(1),
        by = "month"
      ),
      format = "%Y-%m"
    )
  )
}


#' Count NA
#'
#' count_na count the number of NA in a vector
#'
#' @param x a vector
#'
#' @return a tibble
#' @export
#'
#' @examples
#'
#' table_training$code_departement %>% count_na()
#'
#' plyr::ldply(.data = table_training, .fun = count_na)
#'
count_na <- function(x) {
  x %>%
    is.na() %>%
    factor() %>%
    forcats::fct_count()
}


#' Detect missing values
#'
#' Detect na takes a table as input and returns a table with the number and the share of missing values by column
#'
#' @param table name of the input variable
#'
#' @return a table with 3 columns (variable, n_missing, share_missing)
#' @export
#'
#' @examples
#'
#' tibble::tibble(col1 = c(1,NA,3,4), col2 = c(NA,NA,NA,NA)) %>%
#' detect_na()
#'
detect_na <- function(table) {
  table %>%
    plyr::ldply(.data = ., .fun = count_na, .id = "variable") %>%
    dplyr::group_by_(~ variable) %>%
    dplyr::mutate_(
      .dots = list(
        "share_missing" = lazyeval::interp(~ 100 * x / sum(x), x = quote(n))
      )
    ) %>%
    dplyr::filter_(
      .dots = list(~ f == TRUE)
    ) %>%
    dplyr::select_(
      .dots = list(~variable, "n_missing" = ~n, ~share_missing)
    )

}


#' Count infinite
#'
#' @param x a vector
#'
#' @return a tibble
#' @export
#'
#' @examples
#'
#' count_infinite(x = c(Inf, 1, 2, 3, 0, -Inf))
#'
count_infinite <- function(x) {
  x %>%
    is.infinite() %>%
    factor() %>%
    forcats::fct_count()
}


#' Detect infinite values
#'
#' @param table a table
#'
#' @return a tibble
#' @export
#'
#' @examples
#'
#' tibble::tibble(col1 = c(1,Inf,3,4), col2 = c(-Inf,Inf,-Inf,0)) %>%
#' detect_infinite()
#'
detect_infinite <- function(table) {

  table %>%
    plyr::ldply(.data = ., .fun = count_infinite, .id = "variable") %>%
    dplyr::group_by_(~ variable) %>%
    dplyr::mutate_(
      .dots = list(
        "share_infinite" = lazyeval::interp(~ 100 * x / sum(x), x = quote(n))
      )
    ) %>%
    dplyr::filter_(
      .dots = list(~ f == TRUE)
    ) %>%
    dplyr::select_(
      .dots = list(~variable, "n_infinite" = ~n, ~share_infinite)
    )

}

#' Make sequence
#'
#' @param start start date
#' @param end end date
#'
#' @return a vector of dates as a string
#' @export
#'
#' @examples
#'
#' make_sequence(start = "2013-01-01", end = "2017-03-01")
#'
make_sequence <- function(start, end) {
  as.character(
    seq(
      from = lubridate::ymd(start),
      to = lubridate::ymd(end),
      by = "month")
  )
}

#' Make fake sequence
#'
#' Cette fonction permet de générer une fausse séquence.
#' C'est utile pour générer l'échantillon sur les effectifs
#' lorsque la date de dernière mise à jour est en retard
#' par rapport aux autres datasets.
#'
#' @param start a start date
#' @param end an end date
#' @param last the last update
#'
#' @return a character vector
#' @export
#'
#' @examples
#'
#' make_fake_sequence(start = "2013-01-01", end = "2017-03-01", last = "2017-01-01")
#'
make_fake_sequence <- function(start, end, last) {
  sequence_ <- make_sequence(start = start, end = end)
  sequence_[sequence_ %in% make_sequence(start = start, end = last) == FALSE] <- last
  return(sequence_)
}

#' Return lastperiods
#'
#' @param x a date
#' @param sequence1 a sequence
#' @param sequence2 a sequence
#'
#' @return a date
#' @export
#'
#' @examples
#'
#' return_lastperiod(
#' x = "2017-02-01",
#' sequence1 = make_sequence(start = "2013-01-01", end = "2017-03-01"),
#' sequence2 = make_fake_sequence(start = "2013-01-01", end = "2017-03-01", last = "2017-01-01")
#' )
#'
return_lastperiod <- function(x, sequence1, sequence2) {
  tibble::tibble(sequence1, sequence2) %>%
    dplyr::filter(sequence1 == x) %>%
    .$sequence2
}
