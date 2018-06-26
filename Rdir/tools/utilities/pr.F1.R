pr.F1 <- function(predicted, outcome){
  PR <-  pr.curve(scores.class0 = predicted,
                  weights.class0 =  as.numeric(outcome),
                  curve = TRUE)

  return( max(2 * PR$curve[, 1] * PR$curve[, 2] / (PR$curve[, 1] + PR$curve[, 2]), na.rm = TRUE))
}
