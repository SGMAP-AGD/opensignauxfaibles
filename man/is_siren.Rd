% Generated by roxygen2: do not edit by hand
% Please edit documentation in R/utilities.R
\name{is_siren}
\alias{is_siren}
\title{Is siren}
\usage{
is_siren(siren)
}
\arguments{
\item{siren}{a string vector which is suspected to include siren number}
}
\value{
a boolean vector with value TRUE if the string includes exactly 9 digits.
}
\description{
Is siren
}
\examples{
is_siren("2015")
is_siren("201512125")
}

