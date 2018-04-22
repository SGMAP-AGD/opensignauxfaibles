compare_prediction_probabilities <-  function(sample1,sample2){

    sample1 <- sample1 %>% rename(prob1 = predicted_default)
    sample2 <- sample2 %>% rename(prob2 = predicted_default)

    join_samples <- sample1 %>%
      full_join(sample2,by = "siret") %>%
      mutate(prob2_prob1 = prob2 - prob1)

    return(join_samples)
}
