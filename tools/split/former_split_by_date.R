former_split_by_date <- function(dateTrain = as.Date("2015-01-01"),dateCross =  as.Date("2016-01-01")){
sample_train <- table_wholesample %>%
  filter(periode == as.character(dateTrain)) %>%
  filter(effectif >= 10)
sample_cross <- table_wholesample %>%
  filter(periode == as.character(dateCross)) %>%
  filter(effectif >= 10)

return(list("train" = sample_train,"cross" = sample_cross, test = tibble()))
}
