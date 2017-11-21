

#' Drop table if exist
#'
#' The default function db_drop_table in the dplyr package returns an error if the table already exist. This creates a lot of errors in the data manipulation pipe. db_drop_table_ifexist check if the table exist and drop the table if and only if it exist.
#'
#' @param db name of the database
#' @param table name of the table
#'
#' @export
#'
#' @examples
#' db_drop_table_ifexist(
#' db = database_signauxfaibles,
#' table = "table_periods"
#' )
#'
db_drop_table_ifexist <- function(db, table) {

  if (dplyr::db_has_table(
    con = db$obj,
    table = table) == TRUE) {

    base::message(base::paste0("Dropping ", table))

    dplyr::db_drop_table(db$obj, table)

  } else {

    base::message(base::paste0("Table ", table, " doesn't exist"))

  }

}

#' Database connect
#'
#' This function reads the file keys.json at the root of the directory and create a connection to the postgre database.
#'
#' The file keys.json should have the following format :
#'
#' {
#' "host": ["127.0.0.1"],
#' "dbname": ["databasename"],
#' "port": ["5433"],
#' "id":["login"],
#' "pw":["password"]
#' }
#'
#' @param file name of the file where keys are stored
#'
#' @return a connection to the signauxfaible database
#' @export
#'
#' @examples
#' database_signauxfaibles <- database_connect()
#'
database_connect <- function(file = "keys.json") {

  keys <- jsonlite::fromJSON(
    rprojroot::find_rstudio_root_file(file)
  )

  dplyr::src_postgres(
    host = keys$host,
    dbname = keys$dbname,
    port = keys$port,
    user = keys$id,
    password = keys$pw
  )

}


#' Map variables
#'
#' @param db database connexion
#'
#' @return a table
#' @export
#'
#' @examples
#' \dontrun{
#' map_variables(db = database_signauxfaibles)
#' }
map_variables <- function(db) {
  purrr::map(
    .x = dplyr::src_tbls(db),
    .f = function(x) {
      tibble::tibble(
        table = x,
        variables = dplyr::tbl(src = db, from = x) %>%
          dplyr::collect(n = 1) %>%
          names()
      )
    }
  ) %>%
  dplyr::bind_rows()

}



#' Has variable
#'
#' @param db database connexion
#' @param table table in the database
#' @param variable variable to be checked
#'
#' @return
#' @export
#'
#' @examples
#' \dontrun{
#' has_variable(db = database_signauxfaibles, table = "table_altares", variable = "siret")
#' }
has_variable <- function(db, table, variable) {
  dplyr::tbl(src = db, from = table) %>%
    dplyr::collect(n = 1) %>%
    names() %>%
    any(. == variable)
}

