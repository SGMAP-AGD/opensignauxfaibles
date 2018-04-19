# Libraries
library(tidyverse)
library(lubridate)
library(assertthat)
library(opensignauxfaibles)
library(mice)
library(caret)
library(broom)
library(randomForest)
library(MLmetrics)

# Sources
source("./tools/data_prep/impute_missing_data_BdF.R")
source("./tools/objective/objective_RJ_LJ_PS.R")
source("./tools/split/split_snapshot_each_month.R")
source("./tools/utilities/elapsed_months.R")

# seed

seed <- 10011
set.seed(seed)

# Actual date
actual <- "2018-03-01"

# Collecting data

database_signauxfaibles <- database_connect()
table_wholesample <-
  collect_wholesample(db = database_signauxfaibles, table = "wholesample") %>%
  as.data.frame()

# Objective
# TODO move to javascript !!

table_wholesample_prep <- objective_RJ_LJ_PS(table_wholesample)

table_wholesample_prep <- table_wholesample_prep %>%
  mutate(outcome = factor(outcome, levels = c(TRUE,FALSE), labels =  c("default", "non_default")))

# TODO Corriger à la source, javascript !!
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

# Feature selection

table_wholesample_sel <- table_wholesample_prep %>%
  select(siret,periode,outcome,outcome_any,date_effet, cut_effectif,cut_growthrate, lag_effectif_missing,
         apart_last12_months, apart_consommee, apart_share_heuresconsommees,
         log_cotisationdue_effectif,
         log_ratio_dettecumulee_cotisation, indicatrice_dettecumulee,
         indicatrice_croissance_dettecumulee,
         nb_debits,
         delai, delai_sup_6mois, taux_marge, financier_ct, financier, delai_fournisseur, poids_frng, dette_fiscale)

# Impute missing data from BdF

mids <-  impute_missing_data_BdF(table_wholesample_sel,seed)

tw_complete <- mids::complete(mids,1) # provisoire. Continuer avec l'objet mids: with(mids, ...)

# Check for NAs, infinites

# detect_na(table_wholesample_sel %>% select(-taux_marge, -financier_ct, -financier, -delai_fournisseur, -poids_frng, -dette_fiscale))
# detect_infinite(table_wholesample_sel %>% select(-taux_marge, -financier_ct, -financier, -delai_fournisseur, -poids_frng, -dette_fiscale))


  #provisoire: retrait NA et inf
  tw_complete <-tw_complete %>% filter(!is.infinite(log_cotisationdue_effectif))

# Split

samples <-
  split_snapshot_each_month(
    tw_complete,
    date_inf = as.Date("2015-01-01"),
    date_sup = as.Date("2016-12-01"),
    frac_train = 0.7,
    frac_cross = 0.3
  )

sample_train <- samples$train
cv_folds <- samples$cv_fold

sample_actual <- tw_complete %>% filter(periode == actual)

# apply algorithm

formula <-  (
  outcome ~ cut_effectif + cut_growthrate + lag_effectif_missing +
    apart_last12_months + apart_consommee + apart_share_heuresconsommees +
    log_cotisationdue_effectif +
    log_ratio_dettecumulee_cotisation + indicatrice_dettecumulee +
    indicatrice_croissance_dettecumulee +
    nb_debits +   delai + delai_sup_6mois +
    taux_marge + financier_ct + financier + delai_fournisseur + poids_frng + dette_fiscale
)


ctrl <-
  trainControl(
    method = "cv",
    classProbs = TRUE,
    summaryFunction = prSummary,
    savePredictions = "all",
    index = cv_folds
  )

randomForest <- train(formula,
                 data = sample_train,
                 method = 'rf',
                 metric = 'AUC',
                 trControl = ctrl,
                 tuneLength = 10,
                 na.action = "na.omit")

prob <- predict(randomForest, newdata = sample_actual,type = "prob")

prediction <- sample_actual %>%
  cbind(prob = prob$default) %>%
  arrange(prob)

