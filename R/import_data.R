
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
