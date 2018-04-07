library("opensignauxfaibles")
library("dplyr")
library("mongolite")

db <- mongo(collection = "algo1", db = "jason")

data <- db$aggregate('[{"$unwind":{"path": "$value"}}]')$value

sample_train <- data %>%
  filter(periode == "2015-01-01 01:00:00") %>%
  filter(effectif >= 10)
sample_test <- data %>%
  filter(periode == "2016-01-01 01:00:00") %>%
  filter(effectif >= 10)
sample_actual <- data %>%
  filter(periode == "2018-02-01 01:00:00") %>%
  filter(effectif >= 10)

formulas_0_12 <- list(
    "dettecumulee" = outcome_0_12 ~ cut_effectif + cut_growthrate + 
    lag_effectif_missing + apart_last12_months + apart_consommee + 
    apart_share_heuresconsommees + log_cotisationdue_effectif +
    log_ratio_dettecumulee_cotisation_12m + indicatrice_dettecumulee_12m
)

sample_train %>% dplyr::slice(1:100)

# contrasts(sample_train) = NULL
# output_prediction_0_12 <- formulas_0_12$dettecumulee %>%
#   glm(
#     formula = .,
#     data = sample_train,
#     family = "binomial"
#     ) %>%
#   broom::augment(
#     newdata = sample_actual,
#     type.predict = "response"
#   ) %>%
#   dplyr::rename(prediction_0_12 = .fitted) 

# ```{r export-top100}
# output_prediction_0_12 %>%
#   dplyr::anti_join(
#     y = compute_filter_proccollectives(
#       db = database_signauxfaibles,
#       .date = "2017-10-01"),
#     by = "numero_compte",
#     copy = TRUE
#   ) %>%
#   dplyr::anti_join(
#     y = compute_filter_ccsv(
#       db = database_signauxfaibles,
#       .date = "2017-10-01"),
#     by = c("numero_compte"),
#     copy = TRUE
#   ) %>%
#   dplyr::arrange(
#     dplyr::desc(
#       prediction_0_12
#     )
#   ) %>%
#   dplyr::filter(is.na(prediction_0_12) == FALSE) %>%
#   dplyr::select(
#         raison_sociale, siret, 
#         libelle_naf_niveau1, 
#         code_departement, 
#         prediction_0_12,
#         effectif,
#         apart_last12_months, apart_consommee, apart_share_heuresconsommees,
#         nb_debits, delai, delai_sup_6mois, everything()
#   ) %>%
#   dplyr::slice(1:100) %>%
#   readr::write_csv(path = "output/table_predictions_2017_11_29.csv")
# ```

# ```{r}
# output_prediction_0_12 %>%
#   dplyr::anti_join(
#     y = compute_filter_proccollectives(
#       db = database_signauxfaibles,
#       .date = "2017-10-01"),
#     by = "numero_compte",
#     copy = TRUE
#   ) %>%
#   dplyr::anti_join(
#     y = compute_filter_ccsv(
#       db = database_signauxfaibles,
#       .date = "2017-10-01"),
#     by = c("numero_compte"),
#     copy = TRUE
#   ) %>%
#   dplyr::arrange(
#     dplyr::desc(
#       prediction_0_12
#     )
#   ) %>%
#   dplyr::filter(is.na(prediction_0_12) == FALSE) %>%
#   readr::write_csv(path = "output/table_predictions_2017_11_29_all.csv")
# ```

