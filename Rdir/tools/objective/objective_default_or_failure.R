objective_default_or_failure <- function(data,n_months, threshold, lookback){

  data <- data %>%
    group_by(siret) %>%
    arrange(siret,periode) %>%
    mutate(defaut_urssaf = check_n_successive_defaults(ratio_dettecumulee_cotisation_12m,n_months,threshold),
           outcome_any = any(defaut_urssaf | outcome_0_12 == 'default'),
           outcome = with_lookback(defaut_urssaf | outcome_0_12 == 'default',lookback)) %>%
    ungroup()

}

check_n_successive_defaults <- function(data, n_months,threshold) {

  exceedance <- data >= threshold;
  max_consecutive <- sequence(rle(exceedance)$lengths) * exceedance
  return(max_consecutive >= n_months)
}

with_lookback <-  function(data,lookback) {
  output <- sapply(1:length(data),
                   FUN= function(x)
                     any(data[x:min(x+lookback,length(data))]))
  return(output)
}


