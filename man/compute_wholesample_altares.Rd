% Generated by roxygen2: do not edit by hand
% Please edit documentation in R/compute_samples.R
\name{compute_wholesample_altares}
\alias{compute_wholesample_altares}
\title{Compute whole sample altares}
\usage{
compute_wholesample_altares(db, name, start, end)
}
\arguments{
\item{db}{a database}

\item{name}{name of the table}

\item{start}{start date}

\item{end}{end date}
}
\value{
a table in the database
}
\description{
Compute whole sample altares
}
\examples{

\dontrun{
compute_wholesample_altares(
db = database_signauxfaibles,
name = "wholesample_altares",
start = "2013-01-01",
end = "2017-03-01")
}

}

