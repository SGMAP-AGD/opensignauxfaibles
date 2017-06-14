
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

#' Import table debits
#'
#' @param path the path to the debit files
#'
#' @return a tibble
#' @export
#'
#' @examples
#' \dontrun{
#' import_table_debit(path = "raw-data/urssaf/Debit_bourgogne_31_07_2016.csv")
#' }
#'
import_table_debit <- function(path) {
  table_debit <- readr::read_csv2(
    file = path,
    col_names = c("numero_compte", "numero_ecart_negatif",
                  "date_traitement_ecart_negatif",
                  "montant_part_ouvriere", "montant_part_patronale",
                  "numero_historique_ecart_negatif",
                  "date_immatriculation_1", "etat_compte",
                  "code_procedure_collective", "siren",
                  "periode", "date_immatriculation",
                  "code_operation_ecart_negatif", "code_motif_ecart_negatif"),
    col_types = readr::cols(
      numero_compte = readr::col_character(),
      numero_ecart_negatif = readr::col_character(),
      date_traitement_ecart_negatif = readr::col_character(),
      montant_part_ouvriere = readr::col_number(),
      montant_part_patronale = readr::col_number(),
      numero_historique_ecart_negatif = readr::col_integer(),
      date_immatriculation_1 = readr::col_character(),
      etat_compte = readr::col_integer(),
      code_procedure_collective = readr::col_character(),
      siren = readr::col_character(),
      periode = readr::col_character(),
      date_immatriculation = readr::col_date(),
      code_operation_ecart_negatif = readr::col_character(),
      code_motif_ecart_negatif = readr::col_character()
    ),
    skip = 1,
    progress = FALSE
  ) %>%
    dplyr::mutate(
      date_traitement_ecart_negatif = convert_urssaf_date(date_traitement_ecart_negatif),
      montant_part_ouvriere = montant_part_ouvriere / 100,
      montant_part_patronale = montant_part_patronale / 100
    ) %>%
    dplyr::filter(is.na(date_traitement_ecart_negatif) == FALSE)

  table_debit <- convert_urssaf_periods_(
    .data = table_debit,
    .variable = ~ periode,
    format = "yyyyqm")

  return(table_debit)
}

