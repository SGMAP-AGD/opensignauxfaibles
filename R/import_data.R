
#' Import table Altares
#'
#' @param path the path to the data
#'
#' @return a tibble with all variables in the table altares
#' @export
#'
#' @examples
#' \dontrun{
#' table_altares <- import_table_altares(path = "raw-data/direccte/RECAP_ALTARES.csv")
#' }
#'

import_table_altares <- function(path) {
  names_table_altares <- readr::read_csv(
    file = path,
    n_max = 10) %>%
    names() %>%
    stringr::str_to_lower() %>%
    stringr::str_replace_all(., pattern = "[éèê]", replacement = "e") %>%
    stringr::str_replace_all(., pattern = "[[:blank:]]", replacement = "_") %>%
    stringr::str_replace_all(., pattern = "'", replacement = "_") %>%
    stringr::str_replace_all(., pattern = "n°", replacement = "numero") %>%
    stringr::str_replace_all(., pattern = "[\\(\\)\\-]", replacement = "") %>%
    stringr::str_replace_all(., pattern = "ô", replacement = "o")

  table <- readr::read_csv(
    file = path,
    col_names = names_table_altares,
    progress = FALSE
  ) %>%
    dplyr::filter(
      is.na(siren) == FALSE,
      siren != "Siren",
      is_siren(siren) == TRUE
    ) %>%
    dplyr::rename(
      date_effet = date_d_effet
    ) %>%
    dplyr::mutate(
      pays = forcats::fct_recode(
        factor(pays), France = "FRANCE", France = "France"
      ),
      date_de_creation_entreprise = lubridate::ymd(date_de_creation_entreprise),
      date_effet = lubridate::ymd(date_effet),
      effectif = as.numeric(effectif)
    )
  return(table)
}


#' Import table apart  (ie Activite partielle)
#'
#' @param path path to the excel file
#'
#' @return a tibble
#' @export
#'
#' @examples
#' \dontrun{
#' import_table_apart(path = "raw-data/direccte/201609_POP_Bourg Fr Comté de janv2008àaoût2016.xls")
#' }
#'
import_table_apart <- function(path) {
  table <- readxl::read_excel(
    path = path,
    sheet = "Liste des établissements ",
    skip = 4,
    col_names = c(
      "code_departement", "libelle_departement",
      "code_commune", "libelle_commune",
      "siret", "raison_sociale",
      "code_naf_700", "libelle_naf_700",
      "date_decision", "numero_decision",
      "numero_avenant",
      "date_debut_periode_autorisee", "date_fin_periode_autorisee",
      "heures_autorisees", "montants_autorises",
      "heures_consommees", "montants_consommes",
      "effectif_concerne", "effectif_etablissement"
    )
  ) %>%
    dplyr::filter(
      libelle_departement != "Total"
    ) %>%
    dplyr::mutate(
      date_decision = lubridate::ymd(date_decision),
      date_debut_periode_autorisee = lubridate::ymd(date_debut_periode_autorisee),
      date_fin_periode_autorisee = lubridate::ymd(date_fin_periode_autorisee)
    )
  return(table)
}


#' Import table cotisation in a CSV format
#'
#'
#' @param path path to the CSV file
#'
#' @return a tibble
#' @export
#'
#' @examples
#'
#' \dontrun{
#' import_table_cotisation_csv(
#' path = "raw-data/urssaf/Urssaf_bourgogne_Cotis_dues_09_2016_01_2017.csv"
#' )
#' }
#'
import_table_cotisation_csv <- function(path) {
  readr::read_csv2(
    file = path,
    col_names = c(
      "numero_compte", "periode_debit", "cotisation_mise_en_recouvrement",
      "cotisation_encaissee_directement", "periode",
      "cotisation_due", "numero_ecart_negatif"
    ),
    col_types = readr::cols(
      numero_compte = readr::col_character(),
      periode_debit = readr::col_character(),
      cotisation_mise_en_recouvrement = readr::col_number(),
      cotisation_encaissee_directement = readr::col_number(),
      periode = readr::col_character(),
      cotisation_due = readr::col_number(),
      numero_ecart_negatif = readr::col_character()
    ),
    skip = 1
  ) %>%
    convert_urssaf_periods_(.data = ., .variable = ~ periode) %>%
    dplyr::select_(
      ~ numero_compte, ~ periodicity, ~ period,
      ~ numero_ecart_negatif, ~ cotisation_mise_en_recouvrement,
      ~ cotisation_encaissee_directement, ~cotisation_due
    )
}
