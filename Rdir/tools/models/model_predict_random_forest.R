model_predict_random_forest <-
  function(formula,
           train_set,
           new_data = NULL,
           mtry = 4) {
    if (!is.null(new_data)) {
      assertthat::assert_that(nrow(train_set %>% semi_join(new_data, by = 'siret')) == 0)
    }

    if (is.data.frame(mtry)) {
      mtry <- mtry$mtry
    }
    set.seed(1900)

    assertthat::assert_that(length(mtry) == 1)

    ##
    ###
    ## Features
    ###
    ##

    data <-
      feature_engineering_random_forest(train_set, new_data, oversampling = TRUE)
    train_set <- data[[1]]
    new_data_long <- data[[2]]

    ##
    ###
    ## Model training
    ###
    ##

    ctrl <-
      trainControl(
        method = "none",
        classProbs = TRUE,
        summaryFunction = prSummary,
        savePredictions = "none"
      )

    grid <- expand.grid(mtry= mtry, splitrule = c("gini"), min.node.size=c(1))
    #grid <- expand.grid(.mtry = mtry)

    my_model  <- train(
      formula,
      data =  train_set %>%
        mutate(outcome = fct_relevel(outcome, c(
          'default', 'non_default'
        ))),
      method = 'ranger',
      metric = 'AUC',
      trControl = ctrl,
      tuneGrid = grid
    )


    if (!is.null(new_data)) {

      n_imputation <- n_distinct(new_data_long$.imp)

       # To avoid to predict several times complete cases
      aux  <- unique(new_data_long %>% select(-.imp)) %>%
        mutate(pred = caret_prob(my_model, sample = .))

      new_data_long <- new_data_long %>%
        left_join(aux)

      pred_long <- new_data_long$pred

      assertthat::assert_that(length(pred_long) %% n_imputation == 0)
      dim(pred_long) <- c(length(pred_long)/ n_imputation, n_imputation)

      pred_short <- pred_long %>%
        rowMeans()

      new_data_long <- new_data_long %>%
        filter(.imp == 1) %>%
        mutate(pred = pred_short) %>%
        select(siret,periode,pred)

      pred <- new_data %>%
        left_join(new_data_long, by = c('siret', 'periode')) %>%
        .$pred




    } else
      pred = NULL

    return(list(pred = pred, model = my_model))
  }
