

#
# sample_train <- data %>%
#   filter(periode == periode_train) %>%
#   filter(effectif >= 10)
# sample_test <- data %>%
#   filter(periode == periode_test) %>%
#   filter(effectif >= 10)
# sample_actual <- data %>%
#   filter(periode == periode_actual) %>%
#   filter(effectif >= 10)
#
#
# output_prediction_0_12 <- (outcome_0_12 ~ cut_effectif + cut_growthrate +
#   lag_effectif_missing + apart_last12_months + apart_consommee +
#   apart_share_heuresconsommees + log_cotisationdue_effectif +
#   log_ratio_dettecumulee_cotisation_12m + indicatrice_dettecumulee_12m) %>%
#   glm(
#     formula = .,
#     data = sample_train,
#     family = "binomial"
#   ) %>%
#   broom::augment(
#     newdata = sample_actual,
#     type.predict = "response"
#   ) %>%
#   dplyr::rename(prediction_0_12 = .fitted)
#
# output_prediction_0_12 %>% select(siret, prediction_0_12) %>% jsonlite::toJSON()
