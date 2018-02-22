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
compute_sample_effectif <- function(db, .date, .periode) {

  .periode <- ifelse(missing(.periode), .date, .periode)
  initial_date_year_month <- format(lubridate::ymd(.date), "%Y-%m")
  initial_date_minus_one_year <- format(
    (lubridate::ymd(.date) - lubridate::dyears(1)),
    "%Y-%m")

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
      ),
      periode = as.character(.periode)
    ) %>%
    dplyr::rename(numero_compte = compte) %>%
    dplyr::select(siret, numero_compte, raison_sociale, periode,
                  code_departement, effectif, lag_effectif, lag_effectif_missing, growthrate_effectif, code_ape
                  )

}

#' Collect sample effectif
#'
#'
#' @param db database
#' @param .date a date
#'
#' @return a table
#' @export
#'
#' @examples
#' \dontrun{
#' collect_sample_effectif(db = database_signauxfaibles, .date = "2017-01-01")
#' }
#'
collect_sample_effectif <- function(db, .date, .periode) {

  compute_sample_effectif(db = db, .date = .date, .periode = .periode) %>%
    dplyr::collect(n = Inf) %>%
    dplyr::distinct(siret, .keep_all = TRUE) %>%
    dplyr::mutate_(
      .dots = list(
        "log_effectif" = ~ log(x = effectif),
        "log_growthrate_effectif" = ~ dplyr::if_else(
          condition = (lag_effectif_missing == FALSE),
          true = log(growthrate_effectif),
          false = 0),
      "region" = ~ forcats::fct_collapse(
        code_departement,
        bourgogne = c("21", "58", "71", "89"),
        franche_comte = c("25", "39", "70", "90")
        )
      )
    ) %>%
    dplyr::select(siret, numero_compte, raison_sociale, periode,
           code_departement, region,
           effectif, log_effectif,
           growthrate_effectif, log_growthrate_effectif,
           lag_effectif_missing, code_ape)
}

#' Compute whole sample effectif
#'
#'
#' @param db database
#' @param name name of the table
#'
#' @return a table in the database
#' @export
#'
#' @examples
#'
#' \dontrun{
#' compute_wholesample_effectif(db = database_signauxfaibles, name = "wholesample_effectif")
#' }
#'
compute_wholesample_effectif <- function(db, name, start, end, last) {

  periods <- make_sequence(start = start, end = end)
  periods2 <- make_fake_sequence(start = start, end = end, last = last)

  db_drop_table_ifexist(db = db, table = name)

  plyr::llply(
    .data = periods,
    .fun = function(x) {
      collect_sample_effectif(
        db = db,
        .date = return_lastperiod(x = x, sequence1 = periods, sequence2 = periods2),
        .periode = x)
      }
    ) %>%
    dplyr::bind_rows() %>%
    dplyr::filter(effectif >= 1) %>%
    insert_multi(
      dest = db,
      name = name,
      slices = 30,
      indexes = list("siret", "numero_compte", "periode")
    )

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
  .dateplus6 <- .date %m+% months(6)
  .dateplus12 <- .date %m+% months(12)
  .dateplus18 <- .date %m+% months(18)
  .dateplus24 <- .date %m+% months(24)

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
    dplyr::mutate_(.dots = list("periode" = ~ as.character(.date))) %>%
    dplyr::select_(
      .dots = list(~ siret, ~ periode, ~ date_effet)
    ) %>%
    dplyr::ungroup() %>%
    dplyr::mutate(
      outcome_0_12 = ifelse((date_effet <= .dateplus12), "default", "non-default"),
      outcome_12_24 = ifelse((date_effet > .dateplus12 & date_effet <= .dateplus24), "default", "non-default"),
      outcome_6_18 = ifelse((date_effet > .dateplus6 & date_effet <= .dateplus18), "default", "non-default")
      )

}

#' Compute whole sample altares
#'
#' @param db a database
#' @param name name of the table
#' @param start start date
#' @param end end date
#'
#' @return a table in the database
#' @export
#'
#' @examples
#'
#' \dontrun{
#' compute_wholesample_altares(
#' db = database_signauxfaibles,
#' name = "wholesample_altares",
#' start = "2013-01-01",
#' end = "2017-03-01")
#' }
#'
compute_wholesample_altares <- function(db, name, start, end) {

  periods <- as.character(seq(
    from = lubridate::ymd(start),
    to = lubridate::ymd(end),
    by = "month")
  )

  db_drop_table_ifexist(db = db, table = name)

  plyr::llply(
    .data = periods,
    .fun = function(x) {
      compute_sample_altares(db = db, .date = x) %>%
      dplyr::collect()
    }
  ) %>%
    dplyr::bind_rows() %>%
    insert_multi(
      dest = db,
      name = name,
      df = .,
      slices = 50,
      indexes = list("siret", "periode")
    )
}

#' Compute prefilter Altares
#'
#' @param db database
#' @param .date date
#'
#' @return a table in the database
#' @export
#'
#' @examples
#' \dontrun{
#' compute_prefilter_altares(
#' db = database_signauxfaibles,
#' .date = "2017-01-01")
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
    dplyr::distinct(siret)

}


#' Compute wholesample prefilter Altares
#'
#' @param db database
#' @param name name of the table
#' @param start start date
#' @param end end date
#'
#' @return a table in the database
#' @export
#'
#' @examples
#' \dontrun{
#' compute_wholesample_prefilter_altares(
#' db = database_signauxfaibles,
#' name = "wholesample_prefilter_altares",
#' start = "2013-01-01",
#' end = "2017-03-01")
#' }
#'
compute_wholesample_prefilter_altares <- function(db, name, start, end) {

  db_drop_table_ifexist(db = db, table = name)

  periods <- as.character(seq(
    from = lubridate::ymd(start),
    to = lubridate::ymd(end),
    by = "month")
  )

  plyr::llply(
    .data = periods,
    .fun = function(x) {
      compute_prefilter_altares(
        db = db,
        .date = x) %>%
        collect(n = Inf) %>%
        mutate(periode = x)
    }
  ) %>%
    dplyr::bind_rows() %>%
    insert_multi(
      dest = db,
      name = name,
      slices = 50,
      indexes = list("siret", "periode")
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
    dplyr::select(
      siren, siret, siege, date_creation_etablissement,
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

  dplyr::tbl(db, "table_activitepartielle") %>%
    dplyr::filter(
      date_debut_periode_autorisee >= start_date,
      date_debut_periode_autorisee < end_date
    ) %>%
    dplyr::group_by(siret) %>%
    dplyr::summarise(apart_last12_months = ifelse(n() >= 1, 1, 0)) %>%
    dplyr::mutate(periode = as.character(.date))

}

#' Compute wholesample Activité partielle
#'
#' @param db database
#' @param name name of the database
#' @param start start date
#' @param end end date
#'
#' @return a table in the database
#' @export
#'
#' @examples
#'
#' \dontrun{
#' compute_wholesample_apart(
#' db = database_signauxfaibles,
#' name = "wholesample_apart",
#' start = "2013-01-01",
#' end = "2017-03-01")
#' }
#'
compute_wholesample_apart <- function(db, name, start, end) {
  db_drop_table_ifexist(db = db, table = name)
  periods <- as.character(seq(
    from = lubridate::ymd(start),
    to = lubridate::ymd(end),
    by = "month")
  )
  plyr::llply(
    .data = periods,
    .fun = function(x) {
      compute_sample_apart(
        db = db,
        .date = x
      ) %>%
        dplyr::collect()
    }
  ) %>%
    dplyr::bind_rows() %>%
    dplyr::copy_to(
      dest = db,
      name = name,
      indexes = list("siret", "periode"),
      temporary = FALSE
    )
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
#' \dontrun{
#' compute_sample_meancotisation(db = database_signauxfaibles, .date = "2017-01-01")
#' }
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
    ) %>%
    dplyr::mutate(periode = as.character(.date)) %>%
    dplyr::select(numero_compte, periode, mean_cotisation_due)

}

#' Compute wholesample meancotisation
#'
#' @param db database
#' @param name name of the output table
#' @param start start date
#' @param end end date
#'
#' @return table in the database
#' @export
#'
#' @examples
#' \dontrun{
#' compute_wholesample_meancotisation(
#' db = database_signauxfaibles,
#' name = "wholesample_meancotisation",
#' start = "2013-01-01",
#' end = "2017-03-01"
#' )
#' }
#'
compute_wholesample_meancotisation <- function(db, name, start, end) {
  db_drop_table_ifexist(db = db, table = name)
  periods <- as.character(seq(
    from = lubridate::ymd(start),
    to = lubridate::ymd(end),
    by = "month")
  )
  plyr::llply(
    .data = periods,
    .fun = function(x) {
      compute_sample_meancotisation(
        db = db,
        .date = x) %>%
        dplyr::collect()
    }
  ) %>%
    dplyr::bind_rows() %>%
    insert_multi(
      df = .,
      dest = db,
      name = name,
      slices = 50,
      indexes = list("numero_compte", "periode")
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

  periode <- .date
  .date <- lubridate::ymd(.date)

  dplyr::tbl(db, from = "table_debit") %>%
    dplyr::filter_(.dots = list(
      ~ periodicity == "monthly",
      ~ date_traitement_ecart_negatif <= .date)
    ) %>%
    dplyr::group_by_(~ numero_compte, ~ period, ~ numero_ecart_negatif) %>%
    dplyr::filter_(
      .dots = list(
        ~numero_historique_ecart_negatif == max(numero_historique_ecart_negatif)
      )
    ) %>%
    dplyr::ungroup() %>%
    dplyr::group_by_(~ numero_compte) %>%
    dplyr::summarise(
      montant_part_ouvriere = sum(montant_part_ouvriere),
      montant_part_patronale = sum(montant_part_patronale)
    ) %>%
    mutate(periode = as.character(periode))

}

#' Compute wholesample dettecumulee
#'
#' @param db database
#' @param name name of the output table
#' @param start start date
#' @param end end date
#'
#' @return a table in the database
#' @export
#'
#' @examples
#'
#' \dontrun{
#' compute_wholesample_dettecumulee(
#' db = database_signauxfaibles,
#' name = "wholesample_dettecumulee",
#' start = "2013-01-01",
#' end = "2017-03-01"
#' )
#' }
#'
compute_wholesample_dettecumulee <- function(db, name, start, end) {
  db_drop_table_ifexist(db = db, table = name)
  periods <- as.character(seq(
    from = lubridate::ymd(start),
    to = lubridate::ymd(end),
    by = "month")
  )
  plyr::llply(
    .data = periods,
    .fun = function(x) {
      compute_sample_dettecumulee(
        db = db,
        .date = x) %>%
        dplyr::collect()
    }
  ) %>%
  dplyr::bind_rows() %>%
  insert_multi(
      df = .,
      dest = db,
      name = name,
      slices = 50,
      indexes = list("numero_compte", "periode")
    )
}


#' Calcul de la dette cumulée à 12 mois
#'
#' @param db a database
#' @param .date a date
#'
#' @return a table in the database
#' @export
#'
#' @examples
#'
#'  \dontrun{
#'  compute_sample_dettecumulee12M(db = database_signauxfaibles, .date = "2017-01-01")
#'  }
compute_sample_dettecumulee_12m <- function(db, .date) {

  periode <- .date
  .date <- lubridate::ymd(.date)

  dplyr::tbl(db, from = "table_debit") %>%
    dplyr::filter_(.dots = list(
      ~ periodicity == "monthly",
      ~ date_traitement_ecart_negatif <= .date)
    ) %>%
    dplyr::semi_join(
      y = get_table_last_n_months(.date = .date, .n_months = 12),
      by = "period",
      copy = TRUE
    ) %>%
    dplyr::group_by_(~ numero_compte, ~ period, ~ numero_ecart_negatif) %>%
    dplyr::filter_(
      .dots = list(
        ~numero_historique_ecart_negatif == max(numero_historique_ecart_negatif)
      )
    ) %>%
    dplyr::ungroup() %>%
    dplyr::group_by_(~ numero_compte) %>%
    dplyr::summarise(
      montant_part_ouvriere_12m = sum(montant_part_ouvriere),
      montant_part_patronale_12m = sum(montant_part_patronale)
    ) %>%
    mutate(periode = as.character(periode))

}

#' compute wholesample dettecumulee 12m
#'
#' @param db a database conexion
#' @param name name of the table
#' @param start start date
#' @param end end date
#'
#' @return a table in the database
#' @export
#'
#' @examples
#'
compute_wholesample_dettecumulee_12m <- function(db, name, start, end) {
  db_drop_table_ifexist(db = db, table = name)
  periods <- as.character(seq(
    from = lubridate::ymd(start),
    to = lubridate::ymd(end),
    by = "month")
  )
  plyr::llply(
    .data = periods,
    .fun = function(x) {
      compute_sample_dettecumulee_12m(
        db = db,
        .date = x) %>%
        dplyr::collect()
    }
  ) %>%
    dplyr::bind_rows() %>%
    insert_multi(
      df = .,
      dest = db,
      name = name,
      slices = 50,
      indexes = list("numero_compte", "periode")
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
      ) %>%
    dplyr::mutate(periode = as.character(.date))

}


#' Compute wholesample lag dette cumulee
#'
#' @param db database
#' @param name name of the table
#' @param start start date
#' @param end end date
#'
#' @return a table in the database
#' @export
#'
#' @examples
#'
#' \dontrun{
#' compute_wholesample_lagdettecumulee(
#' db = database_signauxfaibles,
#' name = "wholesample_lagdettecumulee",
#' start = "2013-01-01",
#' end = "2017-03-01")
#' }
#'
compute_wholesample_lagdettecumulee <- function(db, name, start, end) {
  db_drop_table_ifexist(db = db, table = name)
  periods <- as.character(seq(
    from = lubridate::ymd(start),
    to = lubridate::ymd(end),
    by = "month")
  )
  plyr::llply(
    .data = periods,
    .fun = function(x) {
      compute_sample_lag_dettecumulee(
        db = db,
        .date = x,
        lag = 12) %>%
        dplyr::collect()
    }
  ) %>%
    dplyr::bind_rows() %>%
    insert_multi(
      df = .,
      dest = db,
      name = name,
      slices = 50,
      indexes = list("numero_compte", "periode")
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

  periode <- .date
  .date <- lubridate::ymd(.date)

  dplyr::tbl(src = database_signauxfaibles, from = "table_ccsv") %>%
    dplyr::filter_(.dots = ~ date_de_traitement < .date) %>%
    dplyr::select_(.dots = ~ numero_compte ) %>%
    dplyr::mutate_(.dots = list("periode" = ~ as.character(periode))) %>%
    dplyr::distinct_()

}

#' Compute wholesample CCSV
#'
#' @param db database
#' @param name name of the table
#' @param start start date
#' @param end end date
#'
#' @return a table in the database
#' @export
#'
#' @examples
#'
#' compute_wholesample_ccsv(
#' db = database_signauxfaibles,
#' name = "wholesample_ccsv",
#' start = "2013-01-01",
#' end = "2017-03-01"
#' )
#'
compute_wholesample_ccsv <- function(db, name, start, end) {
  db_drop_table_ifexist(db = db, table = name)
  periods <- as.character(seq(
    from = lubridate::ymd(start),
    to = lubridate::ymd(end),
    by = "month")
  )
  plyr::llply(
    .data = periods,
    .fun = function(x) {
      compute_filter_ccsv(
        db = db,
        .date = x
      ) %>%
        dplyr::collect()
    }
  ) %>%
    dplyr::bind_rows() %>%
    dplyr::copy_to(
      dest = db,
      name = name,
      indexes = list("numero_compte", "periode"),
      temporary = FALSE
    )
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

  .date2 <- lubridate::ymd(.date)

  dplyr::tbl(db, from = "table_debit") %>%
    dplyr::filter_(
      .dots = list(
        ~ periodicity == "monthly",
        ~ code_operation_ecart_negatif == "1"),
        ~  date_traitement_ecart_negatif < .date2
        ) %>%
    dplyr::semi_join(
      y =  get_table_last_n_months(.date = .date, .n_months = n_months),
      by = "period",
      copy = TRUE
    ) %>%
    dplyr::distinct_(.dots = list(~ numero_compte, ~ period)) %>%
    dplyr::group_by_(.dots = ~ numero_compte) %>%
    dplyr::count_() %>%
    dplyr::rename_(.dots = list("nb_debits" = ~ n)) %>%
    dplyr::mutate(periode = as.character(.date))

}

#' Compue wholesample nbdebits
#'
#' @param db database
#' @param name name of the table
#' @param start start date
#' @param end end date
#'
#' @return a table in the database
#' @export
#'
#' @examples
#'
#' \dontrun{
#' compute_wholesample_nbdebits(
#' db = database_signauxfaibles,
#' name = "wholesample_nbdebits",
#' start = "2013-01-01",
#' end = "2017-03-01")
#' }
#'
compute_wholesample_nbdebits <- function(db, name, start, end) {
  db_drop_table_ifexist(db = db, table = name)
  periods <- as.character(seq(
    from = lubridate::ymd(start),
    to = lubridate::ymd(end),
    by = "month")
  )
  plyr::llply(
    .data = periods,
    .fun = function(x) {
      compute_sample_nbdebits(
        db = db,
        .date = x,
        n_months = 12
      ) %>%
        dplyr::collect()
    }
  ) %>%
    dplyr::bind_rows() %>%
    insert_multi(
      df = .,
      dest = db,
      name = name,
      slices = 50,
      indexes = list("numero_compte", "periode")
    )
}

#' Compute filter proccollectives
#'
#' @param db database
#' @param .date a date
#'
#' @return a table in the database
#' @export
#'
#' @examples
#' \dontrun{
#' compute_filter_proccollectives(db = database_signauxfaibles, .date = "2017-01-01")
#' }
#'
compute_filter_proccollectives <- function(db, .date) {

  .date <- lubridate::ymd(.date)

  dplyr::tbl(src = db, from = "table_debit") %>%
    dplyr::filter_(.dots = ~ periodicity == "monthly") %>%
    dplyr::select_(.dots = ~ numero_compte, ~ code_procedure_collective, ~ period) %>%
    dplyr::mutate_(.dots = list("period_date" = ~ sql("to_date(period || -01, 'YYYY-MM-DD')"))) %>%
    dplyr::filter_(
      .dots = list(
        ~ period_date <= .date,
        ~ code_procedure_collective != "0"
      )
    ) %>%
    dplyr::distinct_(.dots = ~ numero_compte)

}

#' Compute sample delais
#'
#' @param db database
#' @param .date date
#'
#' @return a table in the database
#' @export
#'
#' @examples
#' \dontrun{
#' compute_sample_delais(db = database_signauxfaibles, .date = "2017-01-01")
#' }
#'
compute_sample_delais <- function(db, .date) {

  periode <- .date
  .date <- lubridate::ymd(.date)

  dplyr::tbl(src = db, from = "table_delais") %>%
    dplyr::filter(
      date_creation <= .date,
      date_echeance >= .date
    ) %>%
    dplyr::select(numero_compte, "delai_sup_6mois" = indic_6m) %>%
    dplyr::mutate(
      periode = as.character(periode),
      delai = 1,
      delai_sup_6mois = ifelse(delai_sup_6mois == "SUP", 1, 0)
      )

}


#' Compute wholesample delais
#'
#' @param db database
#' @param name name of the table
#' @param start start date
#' @param end end date
#'
#' @return a table in the database
#' @export
#'
#' @examples
#'
#' \dontrun{
#' compute_wholesample_delais(
#' db = database_signauxfaibles,
#' name = "wholesample_delais",
#' start = "2013-01-01",
#' end = "2017-03-01"
#' )
#' }
#'
compute_wholesample_delais <- function(db, name, start, end) {
  db_drop_table_ifexist(db = db, table = name)
  periods <- as.character(seq(
    from = lubridate::ymd(start),
    to = lubridate::ymd(end),
    by = "month")
  )
  plyr::llply(
    .data = periods,
    .fun = function(x) {
      compute_sample_delais(
        db = db,
        .date = x
      ) %>%
        dplyr::collect()
    }
  ) %>%
    dplyr::bind_rows() %>%
    copy_to(
      dest = db,
      name = name,
      temporary = FALSE
    )
}


#' Compute sample activite partielle consommee
#'
#' @param db database
#' @param .date date
#'
#' @return a table in the databaase
#' @export
#'
#' @examples
#'
#' \dontrun{
#' compute_sample_apart_consommee(db = database_signauxfaibles, .date = "2017-03-01")
#' }
#'
compute_sample_apart_consommee <- function(db, .date) {

  periode <- .date
  .date <- lubridate::ymd(.date)
  .lag_date <- .date %m-% months(12)

  tbl_effectif <- dplyr::tbl(
    src = db,
    from = "table_effectif") %>%
    dplyr::semi_join(
      y = get_table_last_n_months(.date = .date, .n_months = 12),
      by = "period",
      copy = TRUE
    ) %>%
    dplyr::filter(effectif > 0) %>%
    dplyr::select(siret, period, effectif)

  tbl_consommee <- dplyr::tbl(
    src = db, from = "table_apart_consommee") %>%
    dplyr::semi_join(
      y = dplyr::tbl(
        src = db, from = "table_activitepartielle") %>%
        dplyr::filter(motif_label != " intemperies") %>%
        dplyr::select(id_da, motif_label),
      by = "id_da"
    ) %>%
    dplyr::mutate(
      period = sql("to_char(date, 'YYYY-MM')")
    ) %>%
    dplyr::select(siret, period, heures_consommees) %>%
    dplyr::filter(heures_consommees > 0) %>%
    dplyr::distinct() %>%
    dplyr::group_by(siret, period) %>%
    dplyr::summarise(heures_consommees = sum(heures_consommees))

  tbl_effectif %>%
    dplyr::left_join(
      y = tbl_consommee,
      by = c("siret", "period")
    ) %>%
    dplyr::mutate(
      heures_consommees = ifelse(is.na(heures_consommees), 0, heures_consommees)
    ) %>%
    dplyr::group_by(siret) %>%
    dplyr::summarise(
      apart_consommee = ifelse(sum(heures_consommees) > 0, 1, 0),
      apart_share_heuresconsommees = 100 * sum(heures_consommees) / sum(effectif * 1607/12)
    ) %>%
    mutate(periode = as.character(periode))

}

#' Compute wholesample activite partielle consommmee
#'
#' @param db name of the database
#' @param name name of the table
#' @param start starting date
#' @param end end date
#'
#' @return a table in the database
#' @export
#'
#' @examples
#'
#' \dontrun{
#' compute_wholesample_apartconsommee(
#' db = database_signauxfaibles,
#' name = "wholesample_apartconsommee",
#' start = "2013-01-01",
#' end = "2017-03-01"
#' )
#' }
#'
compute_wholesample_apartconsommee <- function(db, name, start, end) {
  db_drop_table_ifexist(db = db, table = name)
  periods <- as.character(seq(
    from = lubridate::ymd(start),
    to = lubridate::ymd(end),
    by = "month")
  )
  plyr::llply(
    .data = periods,
    .fun = function(x) {
      compute_sample_apart_consommee(
        db = db,
        .date = x
      ) %>%
        dplyr::collect()
    }
  ) %>%
    dplyr::bind_rows() %>%
    insert_multi(
      df = .,
      dest = db,
      name = name,
      indexes = list("siret", "periode"),
      slices = 50
    )
}

#' Compute whole sample
#'
#' @param db database
#' @param name name of the table
#'
#' @return a table in the database
#' @export
#'
#' @examples
#' \dontrun{
#' compute_wholesample(
#' db = database_signauxfaibles,
#' name = "wholesample")
#' }
#'
compute_wholesample <- function(db, name) {

  db_drop_table_ifexist(db = db, table = name)
  dplyr::left_join(
    x = dplyr::tbl(src = db, from = "wholesample_effectif"),
    y = dplyr::tbl(src = db, from = "wholesample_altares"),
    by = c("siret", "periode")
  ) %>%
    dplyr::anti_join(
      y = tbl(src = db, from = "wholesample_prefilter_altares"),
      by = c("siret", "periode")
    ) %>%
    dplyr::inner_join(
      y = dplyr::tbl(src = db, from = "wholesample_meancotisation"),
      by = c("numero_compte", "periode")
    ) %>%
    dplyr::left_join(
      y = dplyr::tbl(src = db, from = "wholesample_dettecumulee"),
      by = c("numero_compte", "periode")
    ) %>%
    dplyr::left_join(
      y = dplyr::tbl(src = db, from = "wholesample_lagdettecumulee"),
      by = c("numero_compte", "periode")
    ) %>%
    dplyr::left_join(
      y = dplyr::tbl(src = db, from = "wholesample_dettecumulee_12m"),
      by = c("numero_compte", "periode")
      ) %>%
    dplyr::left_join(
      y = dplyr::tbl(src = db, "wholesample_nbdebits"),
      by = c("numero_compte", "periode")
    ) %>%
    dplyr::left_join(
      y = dplyr::tbl(src = db, "wholesample_delais"),
      by = c("numero_compte", "periode")
    ) %>%
    dplyr::left_join(
      y = dplyr::tbl(src = db, "wholesample_apartconsommee"),
      by = c("siret", "periode")
    ) %>%
    dplyr::left_join(
      y = dplyr::tbl(src = db, "wholesample_apart"),
      by = c("siret", "periode")
      ) %>%
    dplyr::left_join(
      y = import_table_naf(path = "data-raw/naf/naf2008_5_niveaux.xls"),
      by = c("code_ape" = "code_naf_niveau5"),
      copy = TRUE) %>%
    dplyr::compute(name = name, temporary = FALSE)

}

#' Collect wholesample
#'
#' @param db database
#' @param table name of the table
#'
#' @return a table
#' @export
#'
#' @examples
#'
#' \dontrun{
#' collect_wholesample(db = database_signauxfaibles, table = "wholesample")
#' }
#'
collect_wholesample <- function(db, table) {
  dplyr::tbl(
    src = db,
    from = "wholesample"
  ) %>%
    dplyr::collect(n = Inf) %>%
    tidyr::replace_na(
      replace = list(
        "montant_part_ouvriere" = 0,
        "montant_part_patronale" = 0,
        "montant_part_ouvriere_12m" = 0,
        "montant_part_patronale_12m" = 0,
        "lag_montant_part_ouvriere" = 0,
        "lag_montant_part_patronale" = 0,
        "nb_debits" = 0,
        "delai" = 0,
        "delai_sup_6mois" = 0,
        "apart_last12_months" = 0
      )
    ) %>%
    dplyr::mutate(
      cut_effectif = forcats::fct_relevel(
        dplyr::case_when(
          .$effectif < 20 ~ "10-20",
          .$effectif < 50 ~ "21-50",
          TRUE ~ "Plus de 50"
        )
      ),
      cut_growthrate = forcats::fct_relevel(
        dplyr::case_when(
          .$lag_effectif_missing == TRUE ~ "manquant",
          .$growthrate_effectif < .80 ~ "moins 20%",
          .$growthrate_effectif < .95 ~ "moins 20 à 5%",
          .$growthrate_effectif < 1.05 ~ "stable",
          .$growthrate_effectif < 1.20 ~ "plus 5 à 20%",
          TRUE ~ "plus 20%"
        ),
        c("stable", "moins 20%", "moins 20 à 5%", "plus 5 à 20%", "plus 20%", "manquant")
      ),
      outcome_0_12 = forcats::fct_explicit_na(
        factor(outcome_0_12, levels = c("non-default", "default")),
        na_level = "non-default"),
      outcome_6_18 = forcats::fct_explicit_na(
        factor(outcome_6_18, levels = c("non-default", "default")),
        na_level = "non-default"),
      outcome_12_24 = forcats::fct_explicit_na(
        factor(outcome_6_18, levels = c("non-default", "default")),
        na_level = "non-default"),
      cotisationdue_effectif = (mean_cotisation_due) / effectif,
      log_cotisationdue_effectif = log(cotisationdue_effectif),
      ratio_dettecumulee_cotisation = (montant_part_ouvriere + montant_part_patronale) / mean_cotisation_due,
      indicatrice_dettecumulee = (montant_part_ouvriere + montant_part_patronale > 0),
      log_ratio_dettecumulee_cotisation = dplyr::if_else(
        condition = indicatrice_dettecumulee == TRUE,
        true = log(ratio_dettecumulee_cotisation),
        false = 0
      ),
      ratio_dettecumulee_cotisation_12m = (montant_part_ouvriere_12m + montant_part_patronale_12m) / mean_cotisation_due,
      indicatrice_dettecumulee_12m = (montant_part_ouvriere_12m + montant_part_patronale_12m > 0),
      log_ratio_dettecumulee_cotisation_12m = dplyr::if_else(
        condition = indicatrice_dettecumulee_12m == TRUE,
        true = log(ratio_dettecumulee_cotisation_12m),
        false = 0
      ),
      indicatrice_croissance_dettecumulee = (
        montant_part_ouvriere + montant_part_patronale > lag_montant_part_ouvriere + lag_montant_part_patronale
      )
    )
}
