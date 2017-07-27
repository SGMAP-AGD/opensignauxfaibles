library("opensignauxfaibles")
library("dplyr")
library("purrr")
library("forcats")
library("corrr")
library("modelr")

database_signauxfaibles <- database_connect()
src_tbls(database_signauxfaibles)
periods <- as.character(seq(
  from = lubridate::ymd("2013-01-01"),
  to = lubridate::ymd("2017-03-01"),
  by = "month")
)
wholesample <- collect_wholesample(db = database_signauxfaibles, table = "wholesample")

model_matrix(
  data = sample_actual,
  formula = ~ cut_effectif + cut_growthrate + lag_effectif_missing +
               apart_last12_months + apart_consommee + apart_share_heuresconsommees +
               log_cotisationdue_effectif +
               log_ratio_dettecumulee_cotisation + indicatrice_dettecumulee +
               indicatrice_croissance_dettecumulee +
               nb_debits +
               delai + delai_sup_6mois +
               libelle_naf_niveau1) %>%
  correlate() %>%
  fashion() %>%
  shave()



model_matrix(
  data = sample_actual,
  formula = ~ cut_effectif + cut_growthrate + lag_effectif_missing +
    apart_last12_months + apart_consommee + apart_share_heuresconsommees +
    log_cotisationdue_effectif +
    log_ratio_dettecumulee_cotisation + indicatrice_dettecumulee +
    indicatrice_croissance_dettecumulee +
    nb_debits +
    delai + delai_sup_6mois +
    libelle_naf_niveau1) %>%
  correlate() %>%
  network_plot(min_cor = .2)



