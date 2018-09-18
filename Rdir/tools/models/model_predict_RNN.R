model_predict_RNN <- function(train_set, new_data) {
  validation_frac <- 0.30
  validation_sirets <-
    tail(unique(train_set$siret),
         validation_frac * n_distinct(train_set$siret))
  validation_set <- train_set %>%
    filter(siret %in% validation_sirets)

  train_set <- train_set %>%
    filter(!siret %in% validation_sirets)


  assertthat::assert_that(nrow(validation_set %>% semi_join(train_set, by = 'siret')) == 0)
  assertthat::assert_that(nrow(train_set) == n_distinct(train_set %>% select(siret, periode)))
  assertthat::assert_that(nrow(validation_set) == n_distinct(validation_set %>% select(siret, periode)))
  assertthat::assert_that(nrow(new_data) == n_distinct(new_data %>% select(siret, periode)))
  ##
  ###
  ##############
  ## Features ##
  ##############
  ###
  ##

  ##
  ###
  #############################
  ## Parmaters and variables ##
  #############################
  ###
  ##
  lookback_monthly_data = 12

  lookback_yearly_data = 3

  na_value = 0

  var_monthly <- c(
    'effectif',
    'cotisation',
    'apart_heures_consommees',
    'apart_motif_recours',
    'etat_proc_collective',
    'montant_part_ouvriere',
    'montant_part_patronale',
    'effectif_entreprise',
    'apart_entreprise',
    'delai',
    'duree_delai',
    'montant_echeancier'
  )

  var_yearly <- c(
    'poids_frng',
    'taux_marge',
    'delai_fournisseur',
    'dette_fiscale',
    'financier_court_terme',
    'frais_financier',
    'age'
  )

  var_other <- c(
    'activite_saisonniere',
    'indice_monoactivite',
    'productif',
    'nbr_etablissements_connus'
  )

  if (!is.null(new_data)) {
    assertthat::assert_that(nrow(train_set %>% semi_join(new_data, by = 'siret')) == 0)
  }
  seed = 13579
  set.seed(seed)

  .list <-
    feature_engineering_RNN(
      train_set = train_set,
      validation_set = validation_set,
      new_data = new_data,
      oversampling = FALSE
    )

  # Oversampling fait dans le générateur !!

  train_set <- .list[[1]]
  validation_set <- .list[[2]]
  new_data <- .list[[3]]


  ##
  ###
  ###########################
  ## Transform data format ##
  ###########################
  ###
  ##

  aux_gen_fun <- function(sirets, data, batch_size, oversampling) {
    data_full <- data %>% filter(siret %in% sirets)
    return(
      generator(
        data = data_full,
        var_yearly = var_yearly,
        lookback_yearly_data = lookback_yearly_data,
        var_monthly = var_monthly,
        lookback_monthly_data = lookback_monthly_data,
        var_other = var_other,
        sirets = unlist(sirets),
        date_inf = '2014-01-01',
        date_sup = '2016-01-01',
        seed = seed,
        batch_size = batch_size,
        oversampling = oversampling
      )
    )
  }


  assertthat::assert_that(nrow(validation_set %>% semi_join(train_set, by = 'siret')) == 0)

  train_gen <-
    aux_gen_fun(unique(train_set$siret), train_set, 128, oversampling = TRUE)
  validation_gen  <-
    aux_gen_fun(unique(validation_set$siret),
                validation_set,
                128,
                oversampling = FALSE)


  if (!is.null(new_data)) {
    new_data_RNN <- df_to_RNN_input(
      siret_period_pairs = new_data %>% select(siret, periode, outcome),
      df_m = new_data,
      df_y = monthly_to_yearly(new_data),
      var_monthly = var_monthly,
      lookback_monthly_data = lookback_monthly_data,
      var_yearly = var_yearly,
      lookback_yearly_data = lookback_yearly_data,
      var_other = var_other,
      na_value = na_value,
      verbose = TRUE
    )
  }



  ##
  ###
  ###########
  ## Model ##
  ###########
  ###
  ##
  monthly_input <-
    layer_input(shape = list(lookback_monthly_data, length(var_monthly)),
                name = "monthly_input")

  monthly_NN <- monthly_input %>%
    layer_gru(units = 32,
              recurrent_dropout = 0.5,
              dropout = 0.1)

  yearly_input <-
    layer_input(shape = list(lookback_yearly_data, length(var_yearly)),
                name = "yearly_input")

  yearly_NN <- yearly_input %>%
    layer_simple_rnn(units = 32,
                     recurrent_dropout = 0.5,
                     dropout = 0.1)


  ape_input <- layer_input(shape = list(1), name = "ape_input")

  ape_NN <- ape_input %>%
    layer_embedding(input_dim = 522,
                    output_dim = 16,
                    name = 'embedding') %>%
    layer_flatten()

  other_input <-
    layer_input(shape = list(length(var_other)), name = "other_input")

  other_NN <- other_input %>%
    layer_dense(units = 16, activation = 'relu')


  concatenated <-
    layer_concatenate(list(yearly_NN, monthly_NN, ape_NN, other_NN))

  output <-  concatenated %>%
    layer_dense(units = 128, activation = 'relu') %>%
    layer_dense(units = 128, activation = 'relu') %>%
    layer_dense(units = 128, activation = 'relu') %>%
    layer_dense(units = 128, activation = 'relu') %>%
    layer_dense(units = 64, activation = 'relu') %>%
    layer_dense(units = 1, activation = 'sigmoid')

  model <- keras_model(list(monthly_input,
                            yearly_input,
                            ape_input,
                            other_input),
                       output)

  ##
  ###
  ###############
  ## Callbacks ##
  ###############
  ###
  ##

  callback_list <- list(
    callback_early_stopping(monitor = "val_loss",
                            patience = 3),
    callback_model_checkpoint(
      filepath = find_rstudio_root_file('pipes', 'RNNs', 'RNN.h5'),
      monitor = "val_loss",
      save_best_only = TRUE
    )
  )

  ##
  ###
  #############
  ## Control ##
  #############
  ###
  ##

  model %>% compile(
    # optimizer = optimizer_rmsprop(lr = 0.0001),
    optimizer = optimizer_adam(lr = 0.0001),
    loss = "binary_crossentropy",
    metrics = c('acc')
  )

  ##
  ###
  #########
  ## Run ##
  #########
  ###
  ##


  history <- model %>% fit_generator(
    train_gen,
    steps_per_epoch = 40,
    epochs = 50,
    validation_data = validation_gen,
    validation_steps = 15,
    verbose = TRUE,
    callbacks  = callback_list
  )

  ##
  ###
  ###############
  ## Load best ##
  ###############
  ###
  ##


  best <-
    load_model_hdf5(find_rstudio_root_file('pipes', 'RNNs', 'RNN.h5'))
  pred <- c()
  if (!is.null(new_data)) {
    pred <- best %>%
      predict(new_data_RNN[[1]])
  }

  return(list(pred = pred,
              model = best))


}
