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


