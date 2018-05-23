# data_prep_RJLJ_12m_OSbyMonths
# outcome: TRUE si RJ ou LJ ou sauvegarde constatee dans les 12 mois
# Oversampling en faisant une copie des entreprises à chaque mois avant la faillite

set.seed(1409)

database_signauxfaibles <- database_connect()
table_wholesample <-
    collect_wholesample(db = database_signauxfaibles, table = "wholesample")

# Objectif
table_wholesample_prep <- objective_RJ_LJ_PS(table_wholesample)

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


# outcome as default
table_wholesample_prep <- table_wholesample_prep %>%
  mutate(outcome = factor(outcome, levels = c(TRUE,FALSE), labels =  c("default", "non_default")))

### Feature selection

table_wholesample_sel <-
  table_wholesample_prep %>% select(
    siret,
    periode,
    outcome,
    outcome_any,
    date_defaillance,
    poids_frng,
    taux_marge,
    delai_fournisseur,
    dette_fiscale,
    financier_ct,
    financier,
    montant_part_patronale,
    montant_part_ouvriere,
    starts_with("log_effectif"),
    starts_with("ratio_dettecumulee"),
    starts_with('nb_debits'),
    starts_with('apart_share_heuresconsommees'),
    outcome,-ratio_dettecumulee_cotisation_12m
  )

# Préparation séries temporelles

# MOCHE pour l'instant fait des calculs pour les dropper
# Plus ne remplit pas les NA des lags -.-'
ts_wholesample <- table_wholesample_sel %>%
  spread_relative_periods("log_effectif",11) %>%
  spread_relative_periods("ratio_dettecumulee_cotisation",11) %>%
  spread_relative_periods("nb_debits",11) %>%
  spread_relative_periods("apart_share_heuresconsommees",11) %>%
  spread_relative_periods("montant_part_patronale",11) %>%
  spread_relative_periods("montant_part_ouvriere",11)

ts_wholesample <- ts_wholesample %>%
  filter(!is.na(outcome))


### Check NAs, Infinites etc.
detect_na(ts_wholesample)
detect_infinite(ts_wholesample)

# SPLIT
samples <-
  split_snapshot_each_month(
    ts_wholesample,
    date_inf = as.Date("2015-01-01"),
    date_sup = as.Date("2016-12-01")
  )

sample_train <- samples$train
cv_folds <- samples$cv_fold
sample_test <- samples$test



# filtrage is.na etonnant: devrait se faire au niveau de spread_relative_periods
sample_train <- sample_train %>%
  filter(!(outcome == "non_default"  & outcome_any == "default"))

# Impute missing data
sample_train <- sample_train %>%
  impute_missing_data()

sample_test <- sample_test %>%
  impute_missing_data()

