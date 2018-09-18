 monthly_to_yearly <- function(df_m){

   cat('Transforming monthly data to yearly data ... ')

   df_y <- df_m %>%
     collapse_by("1 year",
                 side = 'start',
                 start_date = floor_index(min(.$periode),'1 Y'),
                 clean = TRUE)
  indices <- group_indices(df_y,siret,periode)
   df_y <- df_y %>%
     mutate(index = indices) %>%
     arrange(index)

   df_y <- df_y[diff(c(0,df_y$index)) != 0, ]
   cat('FIX ME: not taking last known value !')

   return(df_y)
   cat('done','\n')



 }
