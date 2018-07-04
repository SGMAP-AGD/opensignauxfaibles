model_tune <-  function(model_fun, train, cv_folds, tune_grid){

  n <- length(tune_grid)
  perfs <- list()
  for (i in 1:nrow(tune_grid)){
      aux <- model_cv_eval(function(x,y) model_fun(x,y,tune_grid[i,]),train,cv_folds)
      browser()
      perfs[[i]] <- aux
  }
   browser()
  perfs_array <- lapply(perfs,function(x) x$AUCPR_default)
  best = tune_grid[which.max(perfs_array),]

    return (best)
}
