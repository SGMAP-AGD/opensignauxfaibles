objective_RJ_LJ_PS <- function(wholesample){
  library(dplyr)

  wholesample <- wholesample %>%
    mutate(outcome = (outcome_0_12 == "default"),
           outcome_any = !is.na(date_defaillance))

  return(wholesample)

}
