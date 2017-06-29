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



#' Compute sample ALTARES
#'
#' @param db name of the database
#' @param .date date
#'
#' @return a table in the database
#' @export
#'
#' @examples
#' \dontrun{
#' compute_sample_altares(
#' db = database_signauxfaibles,
#' .date = "2013-01-01"
#' )
#' }
#'
compute_sample_altares <- function(db, .date) {
  .date <- lubridate::ymd(.date)

  dplyr::semi_join(
    x = dplyr::tbl(src = db, from = "table_altares"),
    y = dplyr::tbl(src = db, from = "table_code_rj_lj"),
    by = c("code_de_la_nature_de_l_evenement" = "code")
    ) %>%
    dplyr::select_(
      .dots = list(
        ~ siret,
        ~ code_du_journal,
        ~ code_de_la_nature_de_l_evenement,
        ~ date_effet
      )
    ) %>%
    dplyr::filter_(
      .dots = list(
        ~ code_du_journal == "001",
        ~ date_effet >= .date
      )
    ) %>%
    dplyr::group_by_(.dots = ~ siret) %>%
    dplyr::filter_(
      .dots = list(
        ~ date_effet == min(date_effet)
      )
    ) %>%
    dplyr::select_(
      .dots = list(~ siret, ~ date_effet)
    )

}


#' Compute prefilter altares
#'
#' Cette fonction a pour but de pouvoir retirer toutes les entreprises qui sont déjà en RJ/LJ à la date considérée.
#'
#' @param db name of the database
#' @param .date a date
#'
#' @return a table in the database
#' @export
#'
#' @examples
#' \dontrun{
#' compute_prefilter_altares(
#' db = database_signauxfaibles,
#' .date = "2013-01-01")
#' }
#'
compute_prefilter_altares <- function(db, .date) {

  .date <- lubridate::ymd(.date)

  dplyr::semi_join(
    x = dplyr::tbl(src = db, from = "table_altares"),
    y = dplyr::tbl(src = db, from = "table_code_rj_lj"),
    by = c("code_de_la_nature_de_l_evenement" = "code")
  ) %>%
    dplyr::select_(
      .dots = list(
        ~ siret,
        ~ code_du_journal,
        ~ code_de_la_nature_de_l_evenement,
        ~ date_effet
      )
    ) %>%
    dplyr::filter_(
      .dots = list(
        ~ code_du_journal == "001",
        ~ date_effet < .date
      )
    ) %>%
    dplyr::group_by_(.dots = ~ siret) %>%
    dplyr::filter_(
      .dots = list(
        ~ date_effet == min(date_effet)
      )
    ) %>%
    dplyr::mutate_(
      .dots = list("row_number" = ~ sql("row_number() over(PARTITION BY siret ORDER BY date_effet)"))
    ) %>%
    dplyr::filter_(
      ~ row_number == 1
    ) %>%
    dplyr::select_(
      .dots = list(~ siret, ~ date_effet)
    )

}

#' Compute sample sirene
#'
#' @param db name of the database
#'
#' @return a table in the database
#' @export
#'
#' @examples
#' \dontrun{
#' compute_sample_sirene(
#' db = database_signauxfaibles,
#' )
#' }
#'
compute_sample_sirene <- function(db) {

  dplyr::tbl(db, "table_sirene") %>%
    dplyr::select(siren, siret, siege, date_creation_etablissement,
                  libelle_naf_niveau1, code_naf_niveau1)

}
