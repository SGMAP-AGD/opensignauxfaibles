library("opensignauxfaibles")
library("dplyr")
database_signauxfaibles <- database_connect()
src_tbls(database_signauxfaibles)
periods <- as.character(seq(
  from = lubridate::ymd("2013-01-01"),
  to = lubridate::ymd("2017-03-01"),
  by = "month")
)

compute_wholesample_dettecumulee <- function(db, name, start, end) {
  db_drop_table_ifexist(db = db, table = name)
  periods <- as.character(seq(
    from = lubridate::ymd(start),
    to = lubridate::ymd(end),
    by = "month")
  )
  plyr::llply(
    .data = periods,
    .fun = function(x) {
      compute_sample_dettecumulee(
        db = db,
        .date = x) %>%
        dplyr::collect()
    }
  ) %>%
    dplyr::bind_rows() %>%
    copy_to(
      dest = db, name = name,
      indexes = list("numero_compte", "periode"),
      temporary = FALSE
      )
}

compute_wholesample_dettecumulee(
  db = database_signauxfaibles,
  name = "wholesample_dettecumulee",
  start = "2013-01-01",
  end = "2017-03-01"
  )

compute_sample_dettecumulee(db = database_signauxfaibles, .date = "2017-01-01") %>%
  collect()


compute_sample_dettecumulee(db = database_signauxfaibles, .date = "2017-01-01") %>%
  collect() %>%



compute_wholesample_meancotisation(
  db = database_signauxfaibles,
  name = "wholesample_meancotisation",
  start = "2013-01-01",
  end = "2017-03-01"
  )

compute_wholesample_prefilter_altares(
  db = database_signauxfaibles,
  name = "wholesample_prefilter_altares",
  start = "2013-01-01",
  end = "2017-03-01"
  )


compute_wholesample_altares(
  db = database_signauxfaibles,
  name = "wholesample_altares",
  start = "2013-01-01",
  end = "2017-03-01")


left_join(
  x = tbl(src = database_signauxfaibles, from = "wholesample_effectif"),
  y = tbl(src = database_signauxfaibles, from = "wholesample_altares"),
  by = c("siret", "periode")
  ) %>%
  anti_join(
    y = tbl(src = database_signauxfaibles, from = "wholesample_prefilter_altares"),
    by = c("siret", "periode")
    ) %>%
  inner_join(
    y = compute_sample_sirene(db = database_signauxfaibles),
    by = "siret"
    ) %>%
  inner_join(
    y = tbl(src = database_signauxfaibles, from = "wholesample_meancotisation"),
    by = c("numero_compte", "periode")
    ) %>%
  compute(name = "whole_sample")




