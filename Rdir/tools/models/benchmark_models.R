benchmark_models <- function() {

  results = data_frame(type = character(0), auc = numeric(0))

  log_paul_antoine <- h2o.glm(
    model_id = 'log_paul_antoine',
    x = x_paul_antoine,
    y = y,
    training_frame = train.hex,
    fold_column = "fold_column",
    seed = 123,
    family = 'binomial',
    score_each_iteration = TRUE,
    keep_cross_validation_predictions = TRUE,
    keep_cross_validation_fold_assignment = TRUE
  )

  results[nrow(results)+1,] = list(
    "Algorithme logistique", h2o.auc(h2o.performance(log_paul_antoine,test.hex))
  )

  log_bdf <- h2o.glm(
    model_id = 'log_banque_de_france',
    x = x_minimal,
    y = y,
    training_frame = train.hex,
    fold_column = "fold_column",
    balance_classes = TRUE,
    seed = 123,
    family = 'binomial',
    score_each_iteration = TRUE,
    keep_cross_validation_predictions = TRUE,
    keep_cross_validation_fold_assignment = TRUE
  )

  results[nrow(results)+1,] = list(
    "Ajout données banque de France", h2o.auc(h2o.performance(log_bdf,test.hex))
  )


  model_RF <- h2o.randomForest(model_id = 'randomForest',
                      x= x_minimal,
                      y =y,
                      training_frame = train.hex,
                      balance_classes = TRUE,
                      fold_column = "fold_column",
                      ntrees = 200,
                      seed = 123,
                      max_depth = 5,
                      mtries = 4,
                      score_each_iteration = TRUE,
                      stopping_rounds = 4,
                      stopping_tolerance = 0.0001
  )

  results[nrow(results)+1,] = list(
    "Forêt aléatoire", h2o.auc(h2o.performance(model_RF, test.hex))
  )

  model_new_variables <- h2o.randomForest(model_id = 'randomForest',
                               x= c(x_minimal, new_variables),
                               y =y,
                               training_frame = train.hex,
                               balance_classes = TRUE,
                               fold_column = "fold_column",
                               ntrees = 200,
                               seed = 123,
                               score_each_iteration = TRUE,
                               stopping_rounds = 4,
                               stopping_tolerance = 0.0001
  )

  results[nrow(results)+1,] = list(
    "Nouvelles variables (délais, consolidation entreprise, etc.)", h2o.auc(h2o.performance(model_new_variables, test.hex))
  )

  model_lgb <- h2o.xgboost(
    model_id = 'model_lgb',
    x= c(x_minimal, new_variables),
    y =y,
    training_frame = train.hex,
    tree_method = "hist",
    grow_policy = "lossguide",
    fold_column = "fold_column",
    learn_rate = 0.05,
    max_depth = 4,
    ntrees = 120,
    seed = 123,
    score_each_iteration = TRUE,
    stopping_rounds = 4,
    stopping_tolerance = 0.001,
    quiet_mode =  TRUE,
    keep_cross_validation_predictions = TRUE,
    keep_cross_validation_fold_assignment = TRUE)

  results[nrow(results)+1,] = list(
    "Light Gradient Boosting", h2o.auc(h2o.performance(model_lgb, test.hex))
  )

  model_diane <- h2o.xgboost(
    model_id = 'LGB_diane',
    x= c(x_minimal,
         new_variables,
         ratios_diane),
    y =y,
    training_frame = train.hex,
    tree_method = "hist",
    grow_policy = "lossguide",
    fold_column = "fold_column",
    learn_rate = 0.1,
    max_depth = 4,
    ntrees = 120,
    seed = 123,
    score_each_iteration = TRUE,
    stopping_rounds = 4,
    stopping_tolerance = 0.001,
    quiet_mode =  TRUE,
    keep_cross_validation_predictions = TRUE,
    keep_cross_validation_fold_assignment = TRUE)

  results[nrow(results)+1,] = list(
    "Ratios financiers (Diane)", h2o.auc(h2o.performance(model_diane, test.hex))
  )

  best_model <- h2o.xgboost(
    model_id = 'best_model',
    x= x_distrib,
    y =y,
    training_frame = train.hex,
    tree_method = "hist",
    grow_policy = "lossguide",
    fold_column = "fold_column",
    learn_rate = 0.05,
    max_depth = 4,
    ntrees = 120,
    seed = 123,
    score_each_iteration = TRUE,
    stopping_rounds = 4,
    stopping_tolerance = 0.001,
    quiet_mode =  TRUE,
    keep_cross_validation_predictions = TRUE,
    keep_cross_validation_fold_assignment = TRUE)

  results[nrow(results)+1,] = list(
    "Nouvelles variables 2 (variations, comparaisons par secteur)", h2o.auc(h2o.performance(best_model,test.hex))
  )

  ensemble <-  h2o.stackedEnsemble(
    x = x_distrib,
    y = y,
    training_frame = train.hex,
    model_id = "ensemble_2",
    base_models = list(best_model,model_lgb))

  h2o.auc(h2o.performance(ensemble,test.hex))


  dataToPlot <- results %>%
    mutate(type = fct_rev(factor(type,level = type)))

  ggplot(dataToPlot, aes(x = type, y = auc, group = 1)) +
    geom_point(stat = 'summary', fun.y = sum) +
    stat_summary(fun_y = sum, geom = 'line') +
    theme(
    axis.text.x = element_text(size = 15),
    axis.text.y = element_text(size = 15),
    axis.title.x = element_text(size = 15),
    axis.title.y = element_blank()
  ) +
    coord_flip()


}
