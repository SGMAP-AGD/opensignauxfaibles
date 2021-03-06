% Generated by roxygen2: do not edit by hand
% Please edit documentation in R/import_data.R
\name{import_table_sirene}
\alias{import_table_sirene}
\title{Import table SIREN}
\usage{
import_table_sirene(path)
}
\arguments{
\item{path}{the path of the siren database in a SAS format}

\item{db}{a database connexion}
}
\value{
a tibble with siren variables
}
\description{
Import table SIREN
}
\examples{
\dontrun{
import_table_sirene(path = "data-raw/raw-data/direccte/bfc.sas7bdat")
}

}

