#' Is siren
#'
#' @param siren a string vector which is suspected to include siren number
#'
#' @return a boolean vector with value TRUE if the string includes exactly 9 digits.
#' @export
#'
#' @examples
#' is_siren("2015")
#' is_siren("201512125")

is_siren <- function(siren) {
  stringr::str_detect(string = siren, pattern = "^[[:digit:]]{9}$")
}
