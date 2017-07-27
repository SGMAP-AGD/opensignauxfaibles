library("opensignauxfaibles")
library("broom")
library("lmtest")
library("ggplot2")
library("dplyr", quietly = TRUE)
database_signauxfaibles <- database_connect()
table_wholesample <- collect_wholesample(db = database_signauxfaibles, table = "wholesample")
sample_train <- table_wholesample %>%
  filter(periode == "2014-01-01")
sample_test <- table_wholesample %>%
  filter(periode == "2015-01-01")
sample_actual <- table_wholesample %>%
  filter(periode == "2017-01-01")

sample_train %>%
  count(outcome_0_12, cut_effectif) %>%
  group_by(cut_effectif) %>%
  mutate(share = 100 * n / sum(n))

glm(
  data = sample_train,
  formula = outcome_0_12 ~ cut_effectif,
  family = binomial
  ) %>%
  tidy()

m0 <- glm(
  data = sample_train,
  formula = outcome_0_12 ~ cut_effectif,
  family = binomial
  )

m1 <- glm(
  data = sample_train,
  formula = outcome_0_12 ~ cut_effectif + cut_growthrate + lag_effectif_missing,
  family = binomial
  )

m0 %>% glance()
m1 %>% glance()

- 2 * (-428.7177 - -419.3278)
1 - pchisq(q = - 2 * (-428.7177 - -419.3278), df = 5)

library("lazyeval")
plot2 <- function(f, data) {
  data %>% select_(f_lhs(f = f))
  }
plot2(f =  outcome_0_12 ~ cut_effectif + cut_growthrate + lag_effectif_missing, data = sample_train)

plot_liftcurve <- function(f, df_train, df_test) {
  glm(
    data = df_train,
    formula = f,
    family = binomial
    ) %>%
    augment(type.predict = "response", newdata = df_test) %>%
    select_(f_lhs(f), ~ .fitted) %>%
    arrange(desc(.fitted)) %>%
    mutate(
      cumden = cumsum(1*(f_lhs(f) == "default")) / sum(1*(f_lhs(f) == "default")),
      rownumber = row_number() / n()
    )
}

plot_liftcurve(f = outcome_0_12 ~ cut_effectif + cut_growthrate + lag_effectif_missing, df_train = sample_train, df_test = sample_test)

  %>%
    ggplot() +
    geom_line(mapping = aes(x = rownumber, y = cumden)) +
    geom_abline(mapping = aes(intercept = 0, slope = 1)
  }

glm(
  data = sample_train,
  formula = outcome_0_12 ~ cut_effectif + cut_growthrate + lag_effectif_missing,
  family = binomial
  ) %>%
  augment(type.predict = "response", newdata = sample_actual) %>%
  select(outcome_0_12, .fitted) %>%
  arrange(desc(.fitted)) %>%
  mutate(
    cumden = cumsum(1*(outcome_0_12 == "default")) / sum(1*(outcome_0_12 == "default")),
    rownumber = row_number() / n()
    ) %>%
  ggplot() +
  geom_line(mapping = aes(x = rownumber, y = cumden)) +
  geom_abline(mapping = aes(intercept = 0, slope = 1))


table_prediction$.fitted %>%

library(ROCR)
predictions <- prediction(predictions = table_prediction$.fitted, labels = table_prediction$outcome_0_12)
performance(predictions, "lift","rpp") %>% plot()

gain <- performance(predictions, "tpr", "rpp")
plot(gain, main = "Gain Chart")

glm(
  data = sample_train,
  formula = outcome_0_12 ~ cut_effectif + cut_growthrate + lag_effectif_missing +
    apart_last12_months + apart_consommee + apart_share_heuresconsommees +
    log_cotisationdue_effectif +
    log_ratio_dettecumulee_cotisation + indicatrice_dettecumulee,
  family = binomial
  ) %>%
  augment(
    type.predict = "response",
    newdata = sample_actual
    ) %>%
  select(siret, raison_sociale, .fitted) %>%
  arrange(desc(.fitted)) %>%
  slice(1:10)

glm(
  formula = f_lhs(f = formulas_0_12$m0) ~ cut_effectif,
  family = "binomial", data = sample_train)

f_rhs(f = )

compute_auc <- function(f, df_train, df_test) {
  glm(formula = f, family = "binomial", data = df_train) %>%
    broom::augment(newdata = df_test,  type.predict = "response") %>%
    pROC::roc(
      lazyeval::f_new(lhs = f_lhs(f), rhs = quote(.fitted)),
      data = . ,
      smooth = FALSE) %>%
    .$auc
  }
compute_auc(f = formulas_0_12$m5, df_train = sample_train, df_test = sample_test)

plot_roc <- function(f, df_train, df_test) {
  roc_curve <- glm(formula = f, family = "binomial", data = df_train) %>%
    broom::augment(newdata = df_test,  type.predict = "response") %>%
    pROC::roc(
      lazyeval::f_new(lhs = f_lhs(f), rhs = quote(.fitted)),
      data = . ,
      smooth = FALSE)

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
plot_roc(f = formulas_0_12$m5, df_train = sample_train, df_test = sample_test)




glm(formula = formulas_0_12[[1]], family = "binomial", data = df_train) %>%
  broom::augment(newdata = df_test,  type.predict = "response") %>%
  pROC::roc(
    lazyeval::f_new(lhs = f_lhs(f), rhs = quote(.fitted)),
    data = . ,
    smooth = FALSE)

tibble::tibble(
  true_positive_rate = roc_curve$sensitivities,
  false_positive_rate = 1 - roc_curve$specificities
)

compute_roc <- function(f, df_train, df_test) {
  glm(formula = f, family = "binomial", data = df_train) %>%
    broom::augment(newdata = df_test,  type.predict = "response") %>%
    pROC::roc(
      lazyeval::f_new(lhs = f_lhs(f), rhs = quote(.fitted)),
      data = . ,
      smooth = FALSE)
  }

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
plot_roc(f = formulas_0_12$m5, df_train = sample_train, df_test = sample_test)

