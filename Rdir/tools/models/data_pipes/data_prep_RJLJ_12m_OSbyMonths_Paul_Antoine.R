# data_prep_RJLJ_12m_OSbyMonths
# outcome: TRUE si RJ ou LJ ou sauvegarde constatee dans les 12 mois
# Oversampling en faisant une copie des entreprises à chaque mois avant la faillite


set.seed(1409)

database_signauxfaibles <- database_connect()
table_wholesample <-
    collect_wholesample(db = database_signauxfaibles, table = "wholesample")

# Objectif
table_wholesample_prep <- objective_RJ_LJ_PS(table_wholesample)
# outcome as default
table_wholesample_prep <- table_wholesample_prep %>%
  mutate(outcome = factor(outcome, levels = c(TRUE,FALSE), labels =  c("default", "non_default")))


### Nettoyage features

# Corriger à la source
table_wholesample_prep <- table_wholesample_prep %>%
  mutate(cut_growthrate = fct_recode(cut_growthrate,
                                     "moins_20pourcent" =  "moins 20%",
                                     "moins_20_a_5_pourcent" = "moins 20 à 5%",
                                     "plus_5_a_20_pourcent" = "plus 5 à 20%",
                                     "plus_20_pourcent" = "plus 20%")) %>%
  mutate(cut_effectif = fct_recode(cut_effectif,
                                     "10_20" = "10-20",
                                     "21_50" = "21-50",
                                     "Plus_de_50" = "Plus de 50"
                                    ))




### Feature selection

table_wholesample_sel <- table_wholesample_prep %>%
  select(siret,periode,outcome,outcome_any,date_effet, cut_effectif,cut_growthrate, lag_effectif_missing,
           apart_last12_months, apart_consommee, apart_share_heuresconsommees,
           log_cotisationdue_effectif,
           log_ratio_dettecumulee_cotisation, indicatrice_dettecumulee,
           indicatrice_croissance_dettecumulee,
           nb_debits,
           delai, delai_sup_6mois, taux_marge, financier_ct, financier, delai_fournisseur, poids_frng, dette_fiscale)

### Check NAs, Infinites etc.
detect_na(table_wholesample_sel)
detect_infinite(table_wholesample_sel)

table_wholesample_sel <- table_wholesample_sel %>%
  filter(!is.infinite(log_cotisationdue_effectif))




### SPLIT
samples <-
  split_snapshot_each_month(
    table_wholesample_sel,
    date_inf = as.Date("2015-01-01"),
    date_sup = as.Date("2016-12-01")
  )

sample_train <- samples$train
cv_folds <- samples$cv_fold
sample_test <- samples$test

# filtrage is.na etonnant: devrait se faire au niveau de spread_relative_periods
sample_train <- sample_train %>%
  filter(!(outcome == "non_default"  & outcome_any == "default")) %>%
  filter(!is.na(outcome))

sample_test <- sample_test %>%
  filter(!is.na(outcome))


# Impute missing data
sample_train <- sample_train %>%
  impute_missing_data()

sample_test <- sample_test %>%
  impute_missing_data()

