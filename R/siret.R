#' Get SIRENE
#'
#' get a list of information from the SIRENE database
#'
#' @param db name of the database
#' @param .siret a siret number
#'
#' @return a list
#' @export
#'
#' @examples
#'
#'
get_sirene <- function(db, .siret) {
  dplyr::tbl(src = db, from = "table_sirene") %>%
    dplyr::filter_(.dots = ~ siret == .siret) %>%
    dplyr::select_(.dots = list(~ siret, ~ siren, ~ date_creation_etablissement, ~ siege, ~ libelle_naf_niveau1)) %>%
    dplyr::collect() %>%
    as.list()
}

#' Get siret
#'
#' @param db database connexion
#' @param .numero_compte account number
#'
#' @return a string
#' @export
#'
#' @examples
#'
#'
get_siret <- function(db, .numero_compte) {
  dplyr::tbl(src = db, "table_effectif") %>%
    dplyr::filter_(.dots = list(~ compte == .numero_compte)) %>%
    dplyr::select_(.dots = ~ siret) %>%
    dplyr::distinct_() %>%
    dplyr::collect() %>%
    .$siret
}

#' Get accountnumber
#'
#' @param db a database connexion
#' @param .siret a siret number
#'
#' @return a string
#' @export
#'
#' @examples
#'
get_accountnumber <- function(db, .siret) {
  dplyr::tbl(src = db, from = "table_effectif") %>%
    dplyr::filter_(.dots = list(~ siret == .siret)) %>%
    dplyr::select_(.dots = ~ compte) %>%
    dplyr::distinct_() %>%
    dplyr::collect() %>%
    .$compte
}


#' is ccsf
#'
#' @param db a database connexion
#' @param siret a valid siret number
#'
#' @return a boolean
#' @export
#'
#' @examples
#'

is_ccsf <- function(db, siret) {

  account_number <- get_accountnumber(
    db = database_signauxfaibles,
    .siret = siret
  )

  dplyr::tbl(src = db, from = "table_ccsv") %>%
    dplyr::filter(numero_compte == account_number) %>%
    dplyr::collect() %>%
    nrow() > 0

}

#' Get CCSV
#'
#' Get information in the CCSV table for a siret number
#'
#' @param db name of the database
#' @param .siret a siret number
#'
#' @return a tibble
#' @export
#'
#' @examples
#'
get_ccsv <- function(db, .siret) {
  .compte <- get_accountnumber(db = db, .siret = .siret)
  dplyr::tbl(src = db, from = "table_ccsv") %>%
    dplyr::filter_(.dots = list(~ compte == .compte)) %>%
    dplyr::collect()
}


#' Get raison sociale
#'
#' @param db a database connexion
#' @param .siret a siret number
#'
#' @return a string
#' @export
#'
#' @examples
#'
#'
get_raisonsociale <- function(db, .siret) {

  dplyr::tbl(src = db, from = "table_effectif") %>%
    dplyr::filter_(.dots = ~ siret == .siret) %>%
    dplyr::distinct_(.dots = ~ raison_sociale) %>%
    dplyr::collect() %>%
    magrittr::extract2("raison_sociale")

}

#' Get effecitf
#'
#' Get the effectif table for a siret number
#'
#' @param db a database connexion
#' @param .siret a siret number
#'
#' @return a tibble
#' @export
#'
#' @examples
#'
get_effectif <- function(db, .siret) {
  dplyr::tbl(src = db, from = "table_effectif") %>%
    dplyr::filter_(.dots = ~ siret == .siret) %>%
    dplyr::select_(.dots = list(~ siret, ~ compte, ~ period, ~ effectif)) %>%
    dplyr::collect()
}

#' Plot effectif
#'
#' @param db a database connexion
#' @param .siret siret number
#'
#' @return a ggplot
#' @export
#'
#' @examples
#'
plot_effectif <- function(db, .siret) {
  get_effectif(db = db, .siret = .siret) %>%
    dplyr::mutate_(
      .dots = list(
        "date" = ~ lubridate::ymd(paste0(period, "-01")),
        "yearmon" = ~ zoo::as.yearmon(date)
      )) %>%
    ggplot2::ggplot() +
    ggplot2::geom_col(
      mapping = ggplot2::aes(x = yearmon, y = effectif),
      color = "white"
    ) +
    zoo::scale_x_yearmon()
}

#' Get meancotisation
#'
#' @param db database connexion
#' @param siret a valid siret number
#'
#' @return a table
#' @export
#'
#' @examples
#'
get_meancotisation <- function(db, siret) {

  account_number <- get_accountnumber(db = db, .siret = siret)
  dplyr::tbl(src = db, from = "wholesample_meancotisation") %>%
    dplyr::filter_(
      .dots = list(~ numero_compte ==  account_number)
    ) %>%
    dplyr::collect() %>%
    dplyr::mutate(periode = lubridate::ymd(periode))

}

#' Plot mean cotisation
#'
#' @param db database connexion
#' @param siret a siret number
#'
#' @return a ggplot
#' @export
#'
#' @examples
#'
plot_meancotisation <- function(db, siret) {
  get_meancotisation(
    db = db,
    siret = siret) %>%
    ggplot2::ggplot() +
    ggplot2::geom_col(
      mapping = ggplot2::aes(
        x = periode,
        y = mean_cotisation_due)
    ) +
    ggplot2::scale_x_date()
}


#' Get ratio cotisation effectif
#'
#' @param db database connexion
#' @param siret siret number
#'
#' @return a table
#' @export
#'
#' @examples
#'
get_ratio_cotisation_effectif <- function(db, siret) {

  get_effectif(db = db, .siret = siret) %>%
    dplyr::mutate(
      periode = paste0(period, "-01")
    ) %>%
    dplyr::inner_join(
      y = get_meancotisation(db = db, siret = siret) %>%
        dplyr::mutate(periode = as.character(periode)),
      by = c("compte" = "numero_compte", "periode")
    ) %>%
    dplyr::mutate_(
      .dots = list("cotisation_effectif" = ~ mean_cotisation_due / effectif)
    ) %>%
    dplyr::select_(
      .dots = list(~ siret, ~ periode, ~ cotisation_effectif)
    )

}

#' PLot ratio cotisation effectif
#'
#' @param db a database connexion
#' @param siret a siret number
#'
#' @return a ggplot
#' @export
#'
#' @examples
#'
plot_ratio_cotisation_effectif <- function(db, siret) {

  get_ratio_cotisation_effectif(db = db, siret = siret) %>%
    dplyr::mutate_(.dots = list("yearmon" = ~ zoo::as.yearmon(periode))) %>%
    ggplot2::ggplot() +
    ggplot2::geom_col(
      mapping = ggplot2::aes_string(
        x = "yearmon",
        y = "cotisation_effectif"
      ),
      color = "white"
    ) +
    zoo::scale_x_yearmon() +
    ggplot2::scale_y_continuous(
      labels = tricky::french_formatting
    )
}

#' Plot dettecumulee
#'
#' @param db a database connexion
#' @param siret a siret number
#' @param variable a variable name ("montant_part_patronale", "dettecumulee", "montant_part_patronale")
#'
#' @return a ggplot
#' @export
#'
#' @examples
#'
plot_dettecumulee <- function(db, siret, variable) {

  account_number <- get_accountnumber(db = database_signauxfaibles, .siret = siret)
  dplyr::tbl(src = db, "wholesample_dettecumulee") %>%
    dplyr::filter(numero_compte == account_number) %>%
    dplyr::collect() %>%
    dplyr::mutate(
      "yearmon" = zoo::as.yearmon(periode),
      "dettecumulee" = montant_part_ouvriere + montant_part_patronale
    ) %>%
    ggplot2::ggplot() +
    ggplot2::geom_col(
      mapping = ggplot2::aes_string(
        x = "yearmon",
        y = variable)
    ) +
    zoo::scale_x_yearmon()

}


