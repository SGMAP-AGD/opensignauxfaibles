library("opensignauxfaibles")
library("dplyr")
database_signauxfaibles <- database_connect()
src_tbls(database_signauxfaibles)
periods <- as.character(seq(
  from = lubridate::ymd("2013-01-01"),
  to = lubridate::ymd("2017-03-01"),
  by = "month")
)

## Compute whole sample

tbl(src = database_signauxfaibles, from = "wholesample") %>%
  dplyr::collect(n = Inf) %>%
  tidyr::replace_na(
    replace = list(
      "montant_part_ouvriere" = 0,
      "montant_part_patronale" = 0,
      "lag_montant_part_ouvriere" = 0,
      "lag_montant_part_patronale" = 0,
      "nb_debits" = 0,
      "delai" = 0,
      "delai_sup_6mois" = 0
    )
  )


  mutate(
    outcome_0_12 = factor(
        dplyr::if_else(
          condition = date_effet <= lubridate::ymd(.date) %m+% months(12),
          true = "default",
          false = "non_default",
          missing = "non_default"),
        levels = c("non_default", "default")
      ),
    outcome_12_24 = factor(
      dplyr::if_else(
          condition = (date_effet > lubridate::ymd(.date) %m+% months(12) & date_effet <= lubridate::ymd(.date) %m+% months(24)),
          true = "default",
          false = "non_default",
          missing = "non_default"),
        levels = c("non_default", "default")),
    outcome_6_18 = factor(
      dplyr::if_else(
          condition = (date_effet > lubridate::ymd(.date) %m+% months(6) & date_effet <= lubridate::ymd(.date) %m+% months(18)),
          true = "default",
          false = "non_default",
          missing = "non_default"),
        levels = c("non_default", "default")),
    cotisationdue_effectif = (mean_cotisation_due) / effectif,
    log_cotisationdue_effectif = log(cotisationdue_effectif),
    ratio_dettecumulee_cotisation = (montant_part_ouvriere + montant_part_patronale) / mean_cotisation_due
  ) %>%
  detect_na()

