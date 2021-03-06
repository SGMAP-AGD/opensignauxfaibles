% Generated by roxygen2: do not edit by hand
% Please edit documentation in R/compute_samples.R
\name{compute_sample_dettecumulee_12m}
\alias{compute_sample_dettecumulee_12m}
\title{Calcul de la dette cumulée à 12 mois}
\usage{
compute_sample_dettecumulee_12m(db, .date)
}
\arguments{
\item{db}{a database}

\item{.date}{a date}
}
\value{
a table in the database
}
\description{
Calcul de la dette cumulée à 12 mois
}
\examples{

 \dontrun{
 compute_sample_dettecumulee12M(db = database_signauxfaibles, .date = "2017-01-01")
 }
}

