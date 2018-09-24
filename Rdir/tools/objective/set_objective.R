set_objective <- function (data,objective){

data[['outcome']] <- fct_recode(
    as.factor(data[[objective]]),
    default = "TRUE",
    non_default = "FALSE"
  ) %>%
    fct_relevel(c("default","non_default"))

cat("L'entraînent se fait désormais sur l'objectif suivant:",objective)
assertthat::assert_that(!any(is.na(data$outcome)))
return(data)
}
