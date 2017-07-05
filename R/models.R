#' DEPRECATED Area under curve
#'
#' @param model a model
#' @param table_test the name of the table including test data
#'
#' @return a scalar value with area under curve
#' @export
#'
#' @examples
#' \dontrun{
#' area_under_curve(model = model_1, table_test = table_test_preprocessed)
#' }
#'
#'
area_under_curve <- function(model, table, outcome) {
  table_test_augmented <- model %>%
    broom::augment(newdata = table,  type.predict = "response")
  roc_curve <- pROC::roc(outcome ~ .fitted, data = table_test_augmented , smooth = FALSE)
  return(roc_curve$auc)
}

#' DEPRECATED Plot roc curve
#'
#' @param model the name of a model
#' @param table the name of the test data
#'
#' @return a ggplot graph
#' @export
#'
#' @examples
#' \dontrun{
#' plot_roc_curve(model = model_departements, table = table_test_preprocessed)
#' }
#'

plot_roc_curve <- function(model, table) {
  table_test_augmented <- model %>% broom::augment(newdata = table, type.predict = "response")
  roc_curve <- pROC::roc(outcome ~ .fitted, table_test_augmented, smooth = FALSE)
  tibble::tibble(
    true_positive_rate = roc_curve$sensitivities,
    false_positive_rate = 1 - roc_curve$specificities
  ) %>%
    ggplot2::ggplot() +
    ggplot2::geom_path(
      mapping = ggplot2::aes(x = false_positive_rate, y =  true_positive_rate)
    ) +
    ggplot2::geom_abline(
      mapping = ggplot2::aes(intercept = 0, slope = 1)
    )
}
