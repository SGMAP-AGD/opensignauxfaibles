average_12m <- function(vec){
  sapply(
    1:length(vec),
    FUN = function(x){
      res <- tail(vec[1:x], 12)
      if (sum(!is.na(res)) > 3)
        return(mean(res , na.rm = TRUE))
      else return(NA)
    }
  )
}
