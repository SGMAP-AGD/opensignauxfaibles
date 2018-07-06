mark_known_sirets <- function(df, name){
  sirets <-readLines(find_rstudio_root_file('..','data-raw',name))

  df <- df %>%
    mutate(connu = as.numeric(siret %in% sirets))

  return(df)
}
