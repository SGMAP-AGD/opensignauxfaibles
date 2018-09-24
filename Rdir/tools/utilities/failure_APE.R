failure_APE <- function(train_data, ..., ape_levels){



  arguments <-  list(...)

  aux_wilson <- function(x){
    lowerWilsonBinCI(length(x), sum(x == 'default')/length(x))
  }

  for (i in seq_along(ape_levels)){

    codes <- train_data %>%
      mutate(codes = substr(code_ape,1,ape_levels[i])) %>%
      group_by(codes) %>%
      summarize(wilson_low_limit = aux_wilson(outcome))

    aux_join <- function(data,.codes){
      data <- data %>%
        mutate(codes = substr(code_ape,1,ape_levels[i])) %>%
        left_join(.codes, by = 'codes') %>%
        select(-codes)
    }

    train_data <- train_data %>% aux_join(codes)
    arguments <- lapply(arguments, aux_join,.codes = codes)
  }


  return(list(train_data, arguments))




}
