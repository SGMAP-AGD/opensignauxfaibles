model_predict_random_forest <- function(formula, train_set,new_data = NULL, mtry = 4) {
  if (!is.null(new_data)){
      assertthat::assert_that(nrow(train_set %>% semi_join(new_data, by = 'siret')) == 0)
  }

  if (is.data.frame(mtry)) {
    mtry <- mtry$mtry
  }
  set.seed(1900)

  assertthat::assert_that(length(mtry)==1)

  ctrl <-
    trainControl(
      method = "none",
      classProbs = TRUE,
      summaryFunction = prSummary,
      savePredictions = "none"
    )

  grid <- expand.grid(mtry= mtry, splitrule = c("gini"), min.node.size=c(1))
  #grid <- expand.grid(.mtry = mtry)
  my_model  <- train(formula,
                     data =  train_set %>%
                       mutate(outcome = fct_relevel(outcome,c('default','non_default'))),
                     method = 'ranger',
                     metric = 'AUC',
                     trControl = ctrl,
                     tuneGrid = grid,
                     na.action = "na.omit")

  if (!is.null(new_data)) {

    pred <- my_model %>% caret_prob(sample = new_data)

  } else pred = NULL

  return(list(pred = pred, model = my_model))
}
