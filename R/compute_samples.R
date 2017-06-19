#' Compute sample effectif
#'
#' Cette fonction permet de créer une table avec les effectifs d'un établissement à une date donnée.
#'
#' @param db the name of the database
#' @param .date a date in the form "yyyy-mm-dd"
#'
#' @return a table in the database
#' @export
#'
#' @examples
#'
#' \dontrun{
#' compute_sample_effectif(
#' db = database_signauxfaibles,
#' .date = "2013-01-01")
#' }
#'
compute_sample_effectif <- function(db, .date) {

  initial_date_year_month <- format(lubridate::ymd(.date), "%Y-%m")
  initial_date_minus_one_year <- format(
    (lubridate::ymd(.date) - lubridate::dyears(1)
    ), "%Y-%m")

  dplyr::tbl(db, "table_effectif") %>%
    dplyr::filter(
      period %in% c(initial_date_minus_one_year, initial_date_year_month)
    ) %>%
    dplyr::mutate(
      period_date = sql("to_date(period || -01, 'YYYY-MM-DD')")
    ) %>%
    dplyr::arrange(siret, period_date) %>%
    dplyr::mutate(
      lag_effectif = sql("lag(effectif) over(PARTITION BY siret ORDER BY period_date)")
    ) %>%
    dplyr::mutate(
      lag_effectif_missing = (is.na(lag_effectif) | lag_effectif == 0)
    ) %>%
    dplyr::filter(
      period == initial_date_year_month
    ) %>%
    dplyr::mutate(
      effectif = sql("cast(effectif as float)"),
      lag_effectif = sql("cast(lag_effectif as float)")
    ) %>%
    dplyr::mutate(
      growthrate_effectif = ifelse(
        lag_effectif_missing == TRUE,
        0,
        effectif / lag_effectif
      )
    ) %>%
    dplyr::rename(numero_compte = compte)

}
