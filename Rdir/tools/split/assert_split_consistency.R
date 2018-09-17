assert_split_consistency <- function(train,validation, test = data.frame(siret = character())){

  assertthat::assert_that(nrow(train %>% semi_join(validation,by = 'siret')) == 0)
  assertthat::assert_that(nrow(train %>% semi_join(test,by = 'siret')) == 0)
  assertthat::assert_that(nrow(test %>% semi_join(validation,by = 'siret')) == 0)

  # train_sets <- lapply(cv_folds, function(x) slice(train,x))
  # cv_sets <- lapply(cv_folds, function(x) slice(train,-x))
  #
  # for (i in seq_along(train_sets)){
  #   assertthat::assert_that(nrow(train_sets[[i]] %>% semi_join(cv_sets[[i]],by = 'siret'))==0)
  # }


}
