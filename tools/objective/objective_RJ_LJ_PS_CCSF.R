objective_RJ_LJ_PS_CCSF <- function(wholesample){
  library(dplyr)
  ccsf <-
    dplyr::tbl(src = database_signauxfaibles, from = "wholesample_ccsv") %>%
    mutate(CCSF = TRUE) %>%
    collect()

  wholesample <- wholesample %>%
    dplyr::left_join(ccsf,
                     by = c("numero_compte", "periode"))

  wholesample <- wholesample %>%
    mutate(outcome = (outcome_0_12 == "default" |
                        !is.na(CCSF))) %>%
    select(-CCSF)

  return(wholesample)

  }
