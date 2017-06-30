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

#' convert urssaf periods to standard format
#'
#' Urssaf store dates in a very specific format.
#' This can be years + 62 for annual values,
#' years + quarter + 0 for quarterly values
#' years + quarter + number of the month in the quarter for monthly values.
#'
#' @param .data a tibble with a period
#' @param .variable the name of the variable with urssaf periods
#' @param format can bee yyqm or yyyyqm
#'
#' @return an new tibble with two more columns : periodicity and period
#' @export
#'
#' @examples
#' convert_urssaf_periods_(.variable = ~ yyyyqm, format = "yyyyqm")
#'

convert_urssaf_periods_ <- function(.data, .variable, format = "yyqm") {
  pattern_ = "([[:digit:]]{2,4})([[:digit:]]{1})([[:digit:]]{1})"
  year_ = stringr::str_replace(
    string = lazyeval::f_eval(~ uq(.variable), data = .data),
    pattern = pattern_,
    replacement = "\\1") %>%
    as.numeric()
  if (format == "yyqm") {
    year_ <- ifelse(year_ <= 20, year_ + 2000, year_ + 1900)
  }
  quarter_ <- stringr::str_replace(
    string = lazyeval::f_eval(~ uq(.variable), data = .data),
    pattern = pattern_,
    replacement = "\\2"
  )
  month_ <- stringr::str_replace(
    string = lazyeval::f_eval(~ uq(.variable), data = .data),
    pattern = pattern_,
    replacement = "\\3"
  )
  periodicity_ <- ifelse(
    stringr::str_detect(
      string = lazyeval::f_eval(~ uq(.variable), data = .data),
      pattern = "[[:digit:]]{2,4}62$") == TRUE,
    "yearly",
    ifelse(
      month_ == "0",
      "quarterly",
      "monthly"
    )
  )
  period_ <- ifelse(
    periodicity_ == "yearly",
    as.character(year_),
    ifelse(
      periodicity_ == "quarterly",
      paste0(year_, "-Q", quarter_),
      paste0(year_, "-",
             stringr::str_pad(
               string = (as.numeric(quarter_) - 1) * 3 + as.numeric(month_),
               side = "left",
               width = 2,
               pad = "0")
      )
    )
  )
  return(dplyr::bind_cols(.data, tibble::tibble(periodicity = periodicity_, period = period_)))
}


#' Convert URSSAF weird dates to normal dates
#'
#' @param weird_date
#'
#' @return a date
#' @export
#'
#' @examples
#' convert_urssaf_date(weird_date = "1010115")
#'
convert_urssaf_date <- function(weird_date) {
  year <- weird_date %>%
    stringr::str_pad(string = .,
                     width = "7",
                     pad = "9",
                     side = "left") %>%
    substr(., 1, 3) %>%
    as.numeric() %>%
    `+`(1900) %>%
    as.character()
  month_day <- weird_date %>%
    stringr::str_pad(string = .,
                     width = "7",
                     pad = "9",
                     side = "left") %>%
    substr(., 4, 7)
  date <- lubridate::ymd(paste0(year, month_day))
  return(date)
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

#' Get table last n months
#'
#' @param .date a date as a string
#' @param .n_months an integer
#'
#' @return a tibble
#' @export
#'
#' @examples
#' get_table_last_n_months(.date = "2017-01-01", .n_months = 12)
#'
get_table_last_n_months <- function(.date, .n_months) {
  tibble::tibble(
    period = format(
      x = seq(
        from = lubridate::ymd(.date) %m-% months(.n_months - 1),
        to = lubridate::ymd(.date),
        by = "month"
      ),
      format = "%Y-%m"
    )
  )
}


#' Count NA
#'
#' count_na count the number of NA in a vector
#'
#' @param x a vector
#'
#' @return a tibble
#' @export
#'
#' @examples
#'
#' table_training$code_departement %>% count_na()
#'
#' plyr::ldply(.data = table_training, .fun = count_na)
#'
count_na <- function(x) {
  x %>%
    is.na() %>%
    factor() %>%
    forcats::fct_count()
}


#' Detect NA
#'
#' Detect NA takes a tibble and return the number and share of missing values for each variable
#'
#' @param table name of the input variable
#'
#' @return a tibble with 3 columns (variable, n_missing, share_missing)
#' @export
#'
#' @examples
#' dplyr::tbl(src = database_signauxfaibles, from = "table_training") %>%
#' dplyr::collect() %>%
#' detect_na()
#'
detect_na <- function(table) {
  table %>%
    plyr::ldply(.data = ., .fun = count_na, .id = "variable") %>%
    dplyr::group_by_(~ variable) %>%
    dplyr::mutate_(
      .dots = list(
        "share_missing" = lazyeval::interp(~ 100 * x / sum(x), x = quote(n))
      )
    ) %>%
    dplyr::filter_(
      .dots = list(~ f == TRUE)
    ) %>%
    dplyr::select_(
      .dots = list(~variable, "n_missing" = ~n, ~share_missing)
    )

}


#' Count infinite
#'
#' @param x a vector
#'
#' @return a tibble
#' @export
#'
#' @examples
#'
#' count_infinite(x = c(Inf, -))
#'
count_infinite <- function(x) {
  x %>%
    is.infinite() %>%
    factor() %>%
    forcats::fct_count()
}


#' Detect infinite values
#'
#' @param table a table
#'
#' @return a tibble
#' @export
#'
#' @examples
#'
detect_infinite <- function(table) {

  table %>%
    plyr::ldply(.data = ., .fun = count_infinite, .id = "variable") %>%
    dplyr::group_by_(~ variable) %>%
    dplyr::mutate_(
      .dots = list(
        "share_infinite" = lazyeval::interp(~ 100 * x / sum(x), x = quote(n))
      )
    ) %>%
    dplyr::filter_(
      .dots = list(~ f == TRUE)
    ) %>%
    dplyr::select_(
      .dots = list(~variable, "n_infinite" = ~n, ~share_infinite)
    )

}


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
