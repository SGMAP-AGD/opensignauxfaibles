# Level (1-a) Wilson confidence interval for proportion
## WILSON, E. B. 1927. Probable inference, the law of succession,
## and statistical inference. Journal of the American Statistical
## Association 22: 209-212.
lowerWilsonBinCI <-  function(n, p, a=0.5) {
  # n sample size
  # share of positive outcome
  z <- qnorm(1-a/2,lower.tail=FALSE)
  l <- 1/(1+1/n*z^2)*(p + 1/2/n*z^2 +
                        z*sqrt(1/n*p*(1-p) + 1/4/n^2*z^2))
  return(l)
}
