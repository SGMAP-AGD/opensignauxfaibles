#############
## Imports ##
#############

# Libraries
library(tidyverse)
library(lubridate)
library(assertthat)
library(mongolite)
library(mice)
library(caret)
library(broom)
library(randomForest)
library(MLmetrics)

# Sources
source("./tools/interface/connect_to_database.R")
source("./tools/data_prep/impute_missing_data_BdF.R")
source("./tools/objective/objective_RJ_LJ_PS.R")
source("./tools/split/split_snapshot_each_month.R")
source("./tools/utilities/elapsed_months.R")
source("./tools/post_analysis/export_top.R")

##################
## Last Periode ##
##################

actual_period <- as.Date("2018-04-01")

#################
## Collecting  ##
##   data      ##
#################

table_wholesample <- connect_to_database('algo1')

#################
## Objective ####
#################

table_wholesample <- objective_RJ_LJ_PS(table_wholesample)

######################
# Feature selection ##
######################

table_wholesample_sel <- table_wholesample %>%
  select(siret,periode,outcome,outcome_any,date_defaillance, cut_effectif,cut_growthrate, lag_effectif_missing,
         apart_last12_months, apart_consommee, apart_share_heuresconsommees,
         log_cotisationdue_effectif,
         log_ratio_dettecumulee_cotisation, indicatrice_dettecumulee,
         indicatrice_croissance_dettecumulee,
         nb_debits,
         delai, delai_sup_6mois, taux_marge, financier_ct, financier, delai_fournisseur, poids_frng, dette_fiscale)

#############################
## Missing data imputation ##
#############################

mids <-  impute_missing_data_BdF(table_wholesample_sel,seed)

tw_complete <- mice::complete(mids,1)


# TODO FIX ME
#provisoire: retrait NA et inf
tw_complete <- tw_complete %>% filter(!is.infinite(log_cotisationdue_effectif))

###########################
## Split train test #######
###########################

seed <- 10011
set.seed(seed)

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

##################
#### Model #######
##################

formula <-  cut_effectif + cut_growthrate + lag_effectif_missing +
  apart_last12_months + apart_consommee + apart_share_heuresconsommees +
  log_cotisationdue_effectif +
  log_ratio_dettecumulee_cotisation + indicatrice_dettecumulee


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


####################
### Prediction #####
####################

tw_complete_long <- mice::complete(mids,'long')
#provisoire: retrait NA et inf
tw_complete_long <- tw_complete_long %>% filter(!is.infinite(log_cotisationdue_effectif))

sample_actual <- tw_complete_long %>%
  filter(periode == actual_period |
           periode == actual_period %m-% months(1) |
           periode == actual_period %m-% months(2))

prob <- predict(randomForest, newdata = sample_actual,type = "prob")

prediction <- sample_actual %>%
  cbind(prob = prob$default) %>%
  group_by(siret) %>%
  summarize(prob = mean(prob), periode =actual_period) %>%
  arrange(desc(prob)) %>%
  as.data.frame()

