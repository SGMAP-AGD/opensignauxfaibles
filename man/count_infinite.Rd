% Generated by roxygen2: do not edit by hand
% Please edit documentation in R/utilities.R
\name{count_infinite}
\alias{count_infinite}
\title{Count infinite}
\usage{
count_infinite(x)
}
\arguments{
\item{x}{a vector}
}
\value{
a tibble
}
\description{
Count infinite
}
\examples{

count_infinite(x = c(Inf, 1, 2, 3, 0, -Inf))

}

