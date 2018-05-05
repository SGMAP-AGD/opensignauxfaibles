impute_missing_data_quick_n_dirty <- function(df){

  tempdata <-  mice(df %>% select(-date_effet),
               m = 1,
              maxit = 2,
              printFlag = FALSE)

df <- complete(tempdata,1)
return(df)
}
