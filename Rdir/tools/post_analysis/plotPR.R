plotPR <- function(model, new_fig = TRUE, test = FALSE){
  if (new_fig)  plot(1, type = 'n', xlab = "recall", ylab = "precision", xlim = c(0,1), ylim = c(0,1))


  if (!test){
    perf <- h2o.performance(model, valid = TRUE)
    pred <- h2o.predict(model, validation.hex)
    true_res <- as.vector(as.numeric(validation.hex['outcome']))
  } else {
    perf <- h2o.performance(model, newdata = test.hex)
    pred <- h2o.predict(model, test.hex)
    true_res <- as.vector(as.numeric(test.hex['outcome']))
  }

  precision <- h2o.precision(perf)
  recall <- h2o.recall(perf)
  lines(recall$tpr,precision$precision, col = rgb(runif(5),runif(5),runif(5)) )

  F2 <- h2o.F2(perf)
  precision_F2 <- precision$precision[which.max(F2$f2)]
  recall_F2 <- recall$tpr[which.max(F2$f2)]
  points(recall_F2, precision_F2, col = 'red', pch= 4)
  text(recall_F2, precision_F2,'F2', pos = 4, col = 'red')

  F1 <- h2o.F1(perf)
  precision_F1 <- precision$precision[which.max(F1$f1)]
  recall_F1 <- recall$tpr[which.max(F1$f1)]
  points(recall_F1, precision_F1, col = 'green', pch= 4)
  text(recall_F1, precision_F1,'F1', pos = 4, col = 'green')


  cat(str(pr.curve(scores.class0 = as.vector(pred$default), weights.class0 = true_res)))
}
