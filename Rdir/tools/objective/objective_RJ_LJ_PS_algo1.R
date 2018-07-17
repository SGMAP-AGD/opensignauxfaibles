objective_RJ_LJ_PS <- function(data){

  data <- data %>%
    mutate(failure = outcome_0_12 == "default",
           failure_any = !is.na(date_defaillance))


  return(data)

}
