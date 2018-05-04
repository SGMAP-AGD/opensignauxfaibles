compare_vects <- function(data,field1,field2){

    smy <- data %>%
      group_by(siret) %>%
      summarize(a = first(!!sym(field1)), b = first(!!sym(field2)))
    return(list(only_first = smy$a & !smy$b, only_second = smy$b & !smy$a))
  }
