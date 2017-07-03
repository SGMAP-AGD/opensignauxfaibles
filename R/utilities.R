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
        from = lubridate::ymd(.date) %m-% months(.n_months - 1),
        to = lubridate::ymd(.date),
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


#' Detect NA
#'
#' Detect NA takes a tibble and return the number and share of missing values for each variable
#'
#' @param table name of the input variable
#'
#' @return a tibble with 3 columns (variable, n_missing, share_missing)
#' @export
#'
#' @examples
#' dplyr::tbl(src = database_signauxfaibles, from = "table_training") %>%
#' dplyr::collect() %>%
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
#' count_infinite(x = c(Inf, -))
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
