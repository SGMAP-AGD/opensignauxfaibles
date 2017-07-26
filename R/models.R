#' Compute AUC
#'
#' @param f a formula
#' @param df_train a training table
#' @param df_test a testing table
#'
#' @return a number
#' @export
#'
#' @examples
#' \dontrun{
#' compute_auc(f = formulas_0_12$m5, df_train = sample_train, df_test = sample_test)
#' }
#'
compute_auc <- function(f, df_train, df_test) {
  glm(formula = f, family = "binomial", data = df_train) %>%
    broom::augment(newdata = df_test,  type.predict = "response") %>%
    pROC::roc(
      lazyeval::f_new(lhs = f_lhs(f), rhs = quote(.fitted)),
      data = . ,
      smooth = FALSE) %>%
    .$auc
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
