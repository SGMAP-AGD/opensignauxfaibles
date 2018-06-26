import_diane <- function(){

  wd <- find_rstudio_root_file('..','data-raw','diane');

  diane1 <- read.csv2(file.path(wd,'utf8_Export_Diane_1_1587.csv'))
  diane2 <- read.csv2(file.path(wd,'utf8_Export_Diane_1588_3174.csv'))
  diane3 <- read.csv2(file.path(wd,'utf8_Export_Diane_3175_4761.csv'))
  diane4 <- read.csv2(file.path(wd,'utf8_Export_Diane_4762_fin.csv'))

  return(rbind(diane1,diane2,diane3,diane4))

}
