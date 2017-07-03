

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
    con = db$con,
    table = table) == TRUE) {

    base::message(base::paste0("Dropping ", table))

    dplyr::db_drop_table(db$con, table)

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
