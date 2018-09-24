model_tune <-  function(model_fun, train, cv_folds, tune_grid, cross_validation = FALSE) {


  n <- length(tune_grid)
  perfs <- list()

  for (i in 1:nrow(tune_grid)) {
    if (cross_validation){
      aux <-
      model_cv_eval(function(x, y)
        model_fun(x, y, tune_grid[i, ]), train, cv_folds)
    } else {
      aux <- model_eval(function(x, y)
        model_fun(x, y, tune_grid[i, ]),
        train %>% slice(cv_folds[[1]]),
        train %>% slice(-cv_folds[[1]]))
    }
    perfs[[i]] <- aux
  }

  perfs_array <- lapply(perfs, function(x)
    x$AUCPR_default)
  best = tune_grid[which.max(perfs_array), ]

  history <- cbind(tune_grid = tune_grid, perf = perfs_array)
  ggplot(history, aes(x = tune_grid, y = perf)) +
    geom_line()

  return (list(history = history, best = best))
}
