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
