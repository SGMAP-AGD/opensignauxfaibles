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


#' Compute sample apart
#'
#' @param db name of the database
#' @param .date initial date
#' @param n_months number of months to consider
#'
#' @return a table in the database
#' @export
#'
#' @examples
#' \dontrun{
#' compute_sample_apart(
#' db = database_signauxfaibles,
#' initial_date = "2013-01-01",
#' n_months = 12
#' )
#' }
#'

compute_sample_apart <- function(db, .date, n_months = 12) {

  start_date <- lubridate::ymd(.date) %m-% months(n_months)
  end_date <- lubridate::ymd(.date)

  dplyr::tbl(db, "table_apart") %>%
    dplyr::filter(
      date_debut_periode_autorisee >= start_date,
      date_debut_periode_autorisee < end_date
    ) %>%
    dplyr::group_by(siret) %>%
    dplyr::summarise(apart_last12_months = ifelse(n() >= 1, 1, 0))

}


#' Compute sample cotisations
#'
#' @param db name of the database
#' @param .date a date to be considered
#'
#' @return a table in the database
#' @export
#'
#' @examples
#'
#'
compute_sample_meancotisation <- function(db, .date) {

  dplyr::tbl(db, "table_cotisation") %>%
    dplyr::filter(periodicity == "monthly") %>%
    dplyr::select(numero_compte, period, numero_ecart_negatif, cotisation_due) %>%
    dplyr::semi_join(
      y = get_table_last_n_months(.date = .date, .n_months = 12),
      by = "period",
      copy = TRUE)  %>%
    dplyr::group_by(numero_compte, period) %>%
    dplyr::summarise(cotisation_due = sum(cotisation_due)) %>%
    dplyr::ungroup() %>%
    dplyr::group_by(numero_compte) %>%
    dplyr::summarise(
      mean_cotisation_due = mean(cotisation_due)
    )

}

#' Compute sample dette cumulee
#'
#' @param db name of the database
#' @param .date a date
#'
#' @return a table in the database
#' @export
#'
#' @examples
#' \dontrun{
#' compute_sample_dettecumulee(db = database_signauxfaibles, .date = "2013-01-01")
#' }
#'
compute_sample_dettecumulee <- function(db, .date) {

  .date <- lubridate::ymd(.date)

  dplyr::tbl(db, from = "table_debit") %>%
    dplyr::filter_(.dots = list(~ periodicity == "monthly"))  %>%
    dplyr::select_(
      ~ numero_compte, ~ period, ~ numero_ecart_negatif,
      ~ numero_historique_ecart_negatif, ~ date_traitement_ecart_negatif,
      ~ montant_part_ouvriere, ~montant_part_patronale
    ) %>%
    dplyr::filter_(.dots = list(~ date_traitement_ecart_negatif <= .date)) %>%
    dplyr::group_by_(~ numero_compte, ~ period, ~ numero_ecart_negatif) %>%
    dplyr::filter_(.dots = list(
      ~numero_historique_ecart_negatif == max(numero_historique_ecart_negatif)
    )
    ) %>%
    dplyr::ungroup() %>%
    dplyr::select_(~ numero_compte, ~ montant_part_ouvriere, ~ montant_part_patronale) %>%
    dplyr::group_by_(~ numero_compte) %>%
    dplyr::summarise_(
      .dots = list(
        "montant_part_ouvriere" = lazyeval::interp(~ sum(x), x = quote(montant_part_ouvriere)),
        "montant_part_patronale" = lazyeval::interp(~ sum(x), x = quote(montant_part_patronale))
      )
    )

}

#' DEPRECATED Compute sample growth dette cumulee
#'
#' @param db name of the database
#' @param .date date
#' @param lag number of lag months
#'
#' @return a table in the database
#' @export
#'
#' @examples
#' \dontrun{
#' compute_sample_growth_dettecumulee(
#' db = database_signauxfaibles,
#' date = "2013-01-01",
#' lag = 12,
#' )
#' }

compute_sample_growth_dettecumulee <- function(db, .date, lag) {

  lag_date <- lubridate::ymd(.date) %m-% months(lag)

  dplyr::left_join(
    x = compute_sample_dettecumulee(db = db, .date = .date),
    y = compute_sample_dettecumulee(db = db, .date = lag_date) %>%
      dplyr::select(
        numero_compte,
        montant_part_ouvriere_old = montant_part_ouvriere,
        montant_part_patronale_old = montant_part_patronale
      ),
    by = "numero_compte") %>%
    dplyr::mutate(
      croissance_dettecumulee_bool = ((montant_part_ouvriere + montant_part_patronale) > (montant_part_patronale_old + montant_part_ouvriere_old)),
      croissance_dettecumulee_new = ((montant_part_ouvriere + montant_part_patronale > 0) & (montant_part_ouvriere_old + montant_part_patronale_old == 0))
    ) %>%
    dplyr::select(numero_compte, croissance_dettecumulee_bool, croissance_dettecumulee_new)

}


#' Compute sample lag dette cumulee
#'
#' @param db database
#' @param .date a date
#' @param lag number of months
#'
#' @return a table in the database
#' @export
#'
#' @examples
#'
#' \dontrun{
#' compute_sample_lag_dettecumulee(db = database_signauxfaibles, lag = 12, .date = "2017-01-01")
#' }
#'
compute_sample_lag_dettecumulee <- function(db, .date, lag) {

  lag_date <- lubridate::ymd(.date) %m-% months(lag)
  compute_sample_dettecumulee(db = db, .date = lag_date) %>%
      dplyr::select(
        numero_compte,
        lag_montant_part_ouvriere = montant_part_ouvriere,
        lag_montant_part_patronale = montant_part_patronale
      )

}


#' Compute filter CCSV
#'
#' @param db a database connexion
#' @param .date a date
#'
#' @return a table in the database
#' @export
#'
#' @examples
#'
#' \dontrun{
#' compute_filter_ccsv(db = database_signauxfaibles, .date = "2013-01-01")
#' compute_filter_ccsv(db = database_signauxfaibles, .date = "2014-01-01")
#' compute_filter_ccsv(db = database_signauxfaibles, .date = "2015-01-01")
#' }
#'
compute_filter_ccsv <- function(db, .date) {

  .date <- lubridate::ymd(.date)

  dplyr::tbl(src = db, from = "table_ccsv") %>%
    dplyr::filter_(.dots = ~ date_de_traitement < .date) %>%
    dplyr::distinct_(.dots = ~ compte)

}

#' Compute sample nbdebits
#'
#' @param db database
#' @param .date date
#' @param n_months n_months
#'
#' @return a table in the database
#' @export
#'
#' @examples
#'
#' \dontrun{
#' compute_sample_nbdebits(
#' db = database_signauxfaibles,
#' .date = "2017-01-01",
#' n_months = 12)
#' }
#'
compute_sample_nbdebits <- function(db, .date, n_months) {

  dplyr::tbl(db, from = "table_debit") %>%
    dplyr::filter_(.dots = list(~ periodicity == "monthly")) %>%
    dplyr::select_(
      ~ numero_compte, ~ period, ~ numero_ecart_negatif,
      ~ numero_historique_ecart_negatif, ~ date_traitement_ecart_negatif,
      ~ montant_part_ouvriere, ~montant_part_patronale
    ) %>%
    dplyr::semi_join(
      y =  get_table_last_n_months(.date = .date, .n_months = n_months),
      by = "period",
      copy = TRUE
    ) %>%
    dplyr::distinct_(.dots = list(~ numero_compte, ~ period)) %>%
    dplyr::group_by_(.dots = ~ numero_compte) %>%
    dplyr::count_() %>%
    dplyr::rename_(.dots = list("nb_debits" = ~ n))

}
