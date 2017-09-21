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
#' \dontrun{
#' get_sirene(db = database_signauxfaibles, .siret = '40094678600011')
#' }
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
#' \dontrun{
#' get_siret(
#' db = database_signauxfaibles,
#' .numero_compte = "267000001600093120"
#' )
#' }
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
#' \dontrun{
#' get_accountnumber(
#' db = database_signauxfaibles,
#' .siret = "33289883200040"
#' )
#' }
#'
get_accountnumber <- function(db, .siret) {
  dplyr::tbl(src = db, "table_effectif") %>%
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
#' \dontrun{
#' get_ccsv(db = database_signauxfaibles, .siret = '40094678600011')
#' }
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
#' \dontrun{
#' get_raisonsociale(db = database_signauxfaibles, .siret = "03578024600027")
#' }
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
#' \dontrun{
#' get_effectif(db = database_signauxfaibles, .siret = "40094678600011")
#' }
#'
get_effectif <- function(db, .siret) {
  dplyr::tbl(src = db, from = "table_effectif") %>%
    dplyr::filter_(.dots = ~ siret == .siret) %>%
    dplyr::select_(.dots = list(~ siret, ~ period, ~ effectif)) %>%
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
#' \dontrun{
#' plot_effectif(db = database_signauxfaibles, .siret = "40094678600011")
#' }
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
