normalize_df <- function(data, means = NULL, stds = NULL) {

    if (is.null(means) || is.null(stds)){
      data <- data %>%
          mutate_if(is.numeric, .funs = scale)
      means <- lapply(data, function(x) attr(x,'scaled:center')) %>%
        plyr::compact()
      stds <- lapply(data, function(x) attr(x,'scaled:scale')) %>%
        plyr::compact()
    } else {
      assertthat::assert_that(length(means) == length(stds))
      assertthat::assert_that(setequal(names(means),names(stds)))
      assertthat::assert_that(all(names(means) %in% names(data)))

      for (i in 1:length(means)){
        variable <- names(means)[i]
        data[variable] <- scale(data[variable], center = means[[variable]], scale = stds[[variable]])
      }
    }

  return(list(data = data,means =means,stds=stds))

}
