
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

#' Import table cotisation
#'
#' @param path path to the data
#'
#' @return a tibble
#' @export
#'
#' @examples
#' \dontrun{
#' import_table_cotisation(path = "raw-data/urssaf/Urssaf_bourgogne_Cotis_dues_histo_31_08_2016.txt")
#' }
#'
import_table_cotisation <- function(path) {
  table <- readr::read_fwf(
    file = path,
    readr::fwf_empty(
      file = path,
      skip = 2,
      col_names = c(
        "numero_compte", "periode_debit", "cotisation_mise_en_recouvrement",
        "cotisation_encaissee_directement", "periode",
        "cotisation_due", "numero_ecart_negatif"
      )
    ),
    skip = 2,
    col_types = readr::cols(
      numero_compte = readr::col_character(),
      periode_debit = readr::col_character(),
      cotisation_mise_en_recouvrement = readr::col_number(),
      cotisation_encaissee_directement = readr::col_number(),
      periode = readr::col_character(),
      cotisation_due = readr::col_number(),
      numero_ecart_negatif = readr::col_character()
    )
  )

  table <- convert_urssaf_periods_(.data = table, .variable = ~ periode)

  table <- table %>%
    dplyr::select(
      numero_compte, periodicity, period,
      numero_ecart_negatif,
      cotisation_mise_en_recouvrement,
      cotisation_encaissee_directement,
      cotisation_due
    )

  return(table)
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

#' Import table effectif
#'
#' @param path path of the source files. Source files are stored in raw-data/urssaf directory
#'
#' @return a tibble with columns siret, raison_sociale, period, effectif, code_departement
#' @export
#'
#' @examples
#' \dontrun{
#' import_table_effectif(
#' path = "raw-data/urssaf/Urssaf_emploi_BFC_200401_201606.csv"
#' )
#' }
#'
import_table_effectif <- function(path) {
  table_effectif <- readr::read_csv2(
    file = path
  ) %>%
    tidyr::gather(
      key = yyyyqm,
      value = effectif,
      - c(rais_soc,UR_EMET, SIRET, dep)
    ) %>%
    dplyr::rename(
      raison_sociale = rais_soc,
      siret = SIRET,
      code_departement = dep
    ) %>%
    dplyr::mutate(
      code_departement = as.character(code_departement),
      yyyyqm = stringr::str_replace(
        string = yyyyqm,
        pattern = "eff([[:digit:]]{6})",
        replacement = "\\1"
      )
    ) %>%
    convert_urssaf_periods_(
      .variable = ~ yyyyqm,
      format = "yyyyqm"
    ) %>%
    dplyr::select(siret, raison_sociale, code_departement, period, effectif) %>%
    dplyr::filter(is.na(effectif) == FALSE) %>%
    dplyr::group_by(siret, period, raison_sociale, code_departement) %>%
    dplyr::summarise(effectif = sum(effectif))
  return(table_effectif)
}


#' Import the new version of table effectif
#'
#' @param path the path to the CSV file
#'
#' @return a tibble
#' @export
#'
#' @examples
#' \dontrun{
#' table_effectif2 <- import_table_effectif2(path = "raw-data/urssaf/Urssaf_emploi_BFC_200401_201609.csv")
#' }
import_table_effectif2 <- function(path) {
  table_effectif <- readr::read_csv2(
    file = path,
    col_types = readr::cols(
      compte = readr::col_character(),
      SIRET = readr::col_character(),
      dep = readr::col_character(),
      ape_ins = readr::col_character()
    )
  ) %>%
    tidyr::gather(
      key = yyyyqm,
      value = effectif,
      - c(rais_soc,UR_EMET, SIRET, dep, compte, ape_ins)
    ) %>%
    dplyr::rename(
      raison_sociale = rais_soc,
      siret = SIRET,
      code_departement = dep,
      code_ape = ape_ins
    ) %>%
    dplyr::mutate(
      yyyyqm = stringr::str_replace(
        string = yyyyqm,
        pattern = "eff([[:digit:]]{6})",
        replacement = "\\1"
      )
    ) %>%
    convert_urssaf_periods_(
      .variable = ~ yyyyqm,
      format = "yyyyqm"
    ) %>%
    dplyr::select(siret, compte, raison_sociale, code_departement, period, effectif, code_ape) %>%
    dplyr::filter(is.na(effectif) == FALSE) %>%
    dplyr::group_by(siret, compte, period, raison_sociale, code_departement, code_ape) %>%
    dplyr::mutate(effectif = as.numeric(effectif)) %>%
    dplyr::summarise(effectif = sum(effectif))
  return(table_effectif)
}

#' Import table naf
#'
#' This function imports the NAF code http://www.insee.fr/fr/methodes/default.asp?page=nomenclatures/naf2008/naf2008.htm
#'
#' @return a tibble with code nav level 5, code naf level 1 and label naf level 1
#' @export
#'
#' @examples
#' \dontrun{
#' import_table_naf(path = "data-raw/insee/naf/naf2008_5_niveaux.xls")
#' }
#'
import_table_naf <- function(path) {
  output_table <- readxl::read_excel(
    path = path,
    sheet = "naf2008_5_niveaux",
    skip = 1,
    col_names = c("code_naf_niveau5", "code_naf_niveau4", "code_naf_niveau3", "code_naf_niveau2", "code_naf_niveau1")
  ) %>%
    dplyr::select(code_naf_niveau5, code_naf_niveau1) %>%
    dplyr::mutate(
      code_naf_niveau5 = stringr::str_replace(
        string = code_naf_niveau5,
        pattern = "([[:digit:]]{2})\\.([[:digit:]]{2}[[:upper:]]{1})",
        replacement = "\\1\\2")
    ) %>%
    dplyr::left_join(
      y = readxl::read_excel(
        path = "data-raw/insee/naf/naf2008_liste_n1.xls",
        sheet = "Feuil1",
      skip = 3,
      col_names = c("code_naf_niveau1", "libelle_naf_niveau1")
    ),
    by = "code_naf_niveau1"
    )
  return(output_table)
}

#' Import table CCSV
#'
#' Cette fonction permet d'importer les données CCSV dans la base de données.
#'
#' @param path chemin vers le fichier CSV
#'
#' @return a tibble : table de données
#' @export
#'
#' @examples
#'
#' \dontrun{
#' "raw-data/urssaf/Bourgogne_ccsf.csv" %>%
#' import_table_ccsv(path = .)
#' }
#'
#' \dontrun{
#' "raw-data/urssaf/FRC_ccsf.csv" %>%
#' import_table_ccsv(path = .)
#' }
#'
import_table_ccsv <- function(path) {
  readr::read_csv2(
    file = path,
    col_types = readr::cols(
      Compte = readr::col_character(),
      `Date de traitement` = readr::col_character(),
      `Code externe du stade` = readr::col_character(),
      `Code externe de l'action` = readr::col_character()
    ),
    locale = readr::locale(decimal_mark = ",")
    ) %>%
    tricky::set_standard_names(.data = .) %>%
    dplyr::mutate_(
      .dots = list(
        "date_de_traitement" = lazyeval::interp(
          ~ convert_urssaf_date(weird_date = x),
          x = quote(date_de_traitement)
        )
      )
    ) %>%
    rename("numero_compte" = compte)
}

#' Import table délais
#'
#' @param path path to the sources files
#'
#' @return a tibble
#' @export
#'
#' @examples
#'
#' \dontrun{
#' import_table_delais(path = "raw-data/urssaf/delais_bourgogne_01_2013_01_2017_utf8.csv")
#' import_table_delais(path = "raw-data/urssaf/delais_franchecomte_01_2013__01_2017.csv")
#' }
#'
import_table_delais <- function(path) {
  readr::read_csv2(
    file = path,
    col_names = c("numero_compte", "numero_contentieux", "date_creation", "date_echeance",
                  "duree_delai", "denomination_premiere_ligne", "indic_6m", "annee_date_creation",
                  "montant_global_de_l_echeancier", "numero_de_structure", "code_externe_du_stade",
                  "code_externe_de_l_action"),
    col_types = readr::cols(
      numero_compte = readr::col_character(),
      numero_contentieux = readr::col_character(),
      numero_de_structure = readr::col_character()
    ),
    skip = 1,
    locale = readr::locale(decimal_mark = ",")
  )
}



#' Import table activite partielle
#'
#' Import de la table d'activité partielle augmentée
#'
#' @param path path
#' @param hta heures totales autorisees
#' @param mta montant total autorise
#' @param effectif_autorise effectif autorise
#'
#' @return a tibble
#' @export
#'
#' @examples
#'
#' \dontrun{
#'   import_table_activite_partielle(
#'   path = "data-raw/activite_partielle/act_partielle_ddes_2012_juin2017.xlsx",
#'   hta = "hta",
#'   mta = "mta",
#'   effectif_autorise = "eff_auto"
#'   )
#' }
#'
import_table_activite_partielle <- function(
  path,
  hta = "hta_heures_totales_autorisees",
  mta = "mta_montant_total_autorise",
  effectif_autorise = "eff_auto_effectif_autorise_a_chomer") {
  table_temp <- readxl::read_excel(path = path) %>%
    tricky::set_standard_names()

  table_temp <- table_temp[, c("id_da", "etab_siret", "eff_etab", hta, mta, "motif_recours_se", effectif_autorise)] %>%
    magrittr::set_colnames(
      value = c("id_da", "siret", "effectif", "hta", "mta", "motif_recours_se", "effectif_autorise")
    ) %>%
    dplyr::mutate(
      motif_label = factor(
        motif_recours_se,
        levels = c(1,2,3,4,5),
        labels = c("conjoncture", "approvisionnement", "intemperies", "restructuration", "autres"))
    )

  return(table_temp)
}

#' Import Apart heures consommées
#'
#' @param path path to the file
#' @param sheet name of the sheet
#' @param skip number of lines to skip
#' @param date variable date
#' @param effectif_concerne  variable effectif
#' @param heures_consommees variable number of hours
#' @param montants variable montants
#'
#' @return a tibble
#' @export
#'
#' @examples
#'
#' \dontrun{
#' import_apart_heuresconsommees(
#' path = "data-raw/activite_partielle/act_partielle_consommée.xlsx",
#' sheet = "INDEMNITE",
#' skip = 0,
#' date = "mois_concerne_par_le_paiement",
#' effectif_concerne = "effectifs_concernes_par_le_paiement",
#' heures_consommees = "heures_consommees",
#' montants = "montants_des_heures_consommees")
#' }
#'
import_apart_heuresconsommees <- function(path, sheet, skip, date, effectif_concerne, heures_consommees, montants) {

  table_temp <- readxl::read_excel(path = path, sheet =  sheet, skip = skip) %>%
    tricky::set_standard_names()

  table_temp <- table_temp[, c("id_da", "etab_siret", date, "eff_etab", effectif_concerne, heures_consommees, montants, "source", "date_payement_annee_mois")] %>%
    magrittr::set_colnames(
      value = c("id_da", "siret", "date", "effectif_etablissement", "effectif_concerne",  "heures_consommees", "montants", "source", "date_payement_annee_mois")
    )

  return(table_temp)

}

# import_apart_heuresconsommees <- function(path) {
#   readxl::read_excel(path = path) %>%
#     tricky::set_standard_names() %>%
#     dplyr::select(
#       id_da,
#       siret = etab_siret,
#       date = mois_concerne_par_le_paiement,
#       effectif_etablissement = eff_etab,
#       effectif_concerne = effectifs_concernes_par_le_paiement,
#       heures_consommees,
#       montants_des_heures_consommees,
#       source,
#       date_payement_annee_mois
#     )
# }

#' Import table SIREN
#'
#' @param path the path of the siren database in a SAS format
#' @param db a database connexion
#'
#' @return a tibble with siren variables
#' @export
#'
#' @examples
#' \dontrun{
#' import_table_sirene(path = "data-raw/raw-data/direccte/bfc.sas7bdat")
#' }
#'
import_table_sirene <- function(path) {

  table_sirene <- haven::read_sas(data_file = path) %>%
    dplyr::select(
      siren = SIREN,
      nic = NIC,
      code_departement = DEPET,
      code_ancienne_region = RPET,
      siege = SIEGE,
      tranche_effectif_salarie = TEFET,
      date_effectif_salarie = DEFET,
      date_creation_etablissement = DCRET,
      code_apet_700 = APET700
    ) %>%
    dplyr::mutate(
      siret = paste0(siren, nic),
      date_creation_etablissement = format(
        zoo::as.yearmon(
          date_creation_etablissement,
          format = "%Y%m"
        ),
        format = "%Y-%m"
      )
    )

  return(table_sirene)

}


