objective_ratio_dettecumulee_cotisation <- function(wholesample,threshold){

  wholesample <- wholesample %>%
    mutate(outcome = (ratio_dettecumulee_cotisation > threshold)) %>%
    group_by(siret) %>%
    mutate(outcome_any = any(outcome))

}
