AUCPR <- function(y_pred, y_true){

  PR <-  pr.curve(scores.class0 = y_pred,
                  weights.class0 =  as.numeric(y_true),curve = TRUE)

  return(PR$auc.integral)
}
