model_eval <- function(model_fun, train_set, validation_set) {
  train_set <- train_set

  output <- model_fun(train_set, validation_set)
  prediction <- output$pred

  AUCPR_failure <-  AUCPR(prediction, validation_set$failure)
  AUCPR_default <- AUCPR(prediction, validation_set$default)
  F1_failure <- pr.F1(prediction, validation_set$failure)
  F1_default <- pr.F1(prediction, validation_set$default)

  cat('Precision Recall AUC for default', AUCPR_default, '\n')
  cat('\n')

  return(
    list(
      AUCPR_failure = AUCPR_failure,
      AUCPR_default = AUCPR_default,
      F1_failure = F1_failure,
      F1_default = F1_default
    )
  )
}
