objective_ratio_dettecumulee_cotisation_et_taux_marge <- function(wholesample,threshold){


  wholesample_temp <- wholesample %>%
    group_by(siret) %>%
    top_n(2,periode) %>%
    summarize(deux_taux_marge_negatifs = all(taux_marge < -0.01))

  wholesample <- wholesample %>%
    left_join(y = wholesample_temp, by = "siret")

  wholesample <- wholesample %>%
    mutate(outcome = (ratio_dettecumulee_cotisation > threshold) | (!is.na(deux_taux_marge_negatifs) & deux_taux_marge_negatifs == TRUE)) %>%
    group_by(siret) %>%
    mutate(outcome_any = any(outcome))

}
