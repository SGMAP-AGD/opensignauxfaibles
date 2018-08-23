quantile_APE <- function(data, variable_names, levels = 1) {

  assertthat::assert_that(all(c('code_naf_niveau1', 'code_ape') %in% names(data)))
  assertthat::assert_that(all(levels >= 1 & levels <= 5))

  # FIX ME: problème de nom dans le mutate_at si une seule variable
  assertthat::assert_that(length(variable_names) > 1)

  for (i in seq_along(levels)) {
    level = levels[i]

    if (level == 1) {
      data <- data %>%
        mutate(.target = as.factor(code_naf_niveau1))
    } else {
      data <- data %>%
        mutate(.target = as.factor(substr(code_ape, 1, levels)))
    }

    aux_fun <- function(x){
      if (all(is.na(x))) return (NA_real_)
      mecdf = ecdf(x)(x)
      return(mecdf)
    }

    my_fun <- c("aux_fun")
    names(my_fun) <- paste0("distrib_APE",level)

    data <- data %>%
      group_by(.target) %>%
      mutate_at(
        .vars = variable_names,
        .funs = funs_(my_fun)) %>%
      ungroup() %>%
      select(-.target)
  }

  return(data)


}
