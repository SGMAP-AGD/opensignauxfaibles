cross_validated_model_fit <- function(model_fun, train_set, cv_folds){

  n = length(cv_folds)
  AUCPR_aux_failure <- numeric(length = n)
  AUCPR_aux_default <- numeric(length = n)
  F1_aux_failure <- numeric(length = n)
  F1_aux_default <- numeric(length = n)

  for (i in seq_along(cv_folds)) {

    aux_train <- sample_train %>%
      slice(cv_folds[[i]])

    aux_cv <- sample_train %>%
      slice(-cv_folds[[i]])

    prediction <- model_fun(aux_train,aux_cv)

    AUCPR_aux_failure[i] <-  AUCPR(prediction, aux_cv$failure)
    AUCPR_aux_default[i] <- AUCPR(prediction, aux_cv$default)
    F1_aux_failure <- pr.F1(prediction,aux_cv$failure)
    F1_aux_default <- pr.F1(prediction,aux_cv$default)

  }
  return(list(
    AUCPR_failure = mean(AUCPR_aux_failure),
    AUCPR_default = mean(AUCPR_aux_default),
    F1_failure = mean(F1_aux_failure),
    F1_default = mean(F1_aux_default)))
}
