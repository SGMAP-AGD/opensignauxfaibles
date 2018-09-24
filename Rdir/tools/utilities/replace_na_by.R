replace_na_by <- function(name,data,na_value) {
  data[is.na(data[,name]), name] = na_value
  return(data)
}
