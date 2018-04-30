plot_comparison <- function(df1,df2, table_wholesample,max_rank = 100){

  table_wholesample <- table_wholesample %>%
    mutate(periode = as.Date(periode))



  df1 <- inner_join(df1 %>% select(siret,periode,prob),
                    table_wholesample, by = c('siret','periode'))
  df2 <- inner_join(df2 %>% select(siret,periode,prob),
                    table_wholesample, by = c('siret','periode'))

  assert_that('prob' %in% colnames(df1) & 'prob' %in% colnames(df2))

  df <- bind_rows(df1,df2,.id = 'id') %>%
    group_by(id) %>%
    mutate(ranking = rank(desc(prob),ties.method = 'random')) %>%
    ungroup() %>%
    tidyr::complete(siret,id) %>%
    group_by(siret) %>%
    filter(any(ranking <= max_rank)) %>%
    arrange(siret,id) %>%
    summarize(effectif = first(na.omit(cut_effectif)),
              prob_diffs = prob[ id == 1 ] - prob[ id == 2 ],
              prob = first(na.omit(prob)),
              new_or_old = factor(
                ifelse(!any(is.na(ranking)) & max(ranking) <= max_rank,
                       'both',
                       ifelse( is.na(ranking[id == 1]) | (!is.na(ranking[id == 2]) & ranking[id == 2] > max_rank), 'only_new',
                               'only_old')),
                levels = c('both','only_new','only_old')
              ),
              ranking = first(na.omit(ranking))
    ) %>%
    mutate(prob_diffs = replace(prob_diffs, is.na(prob_diffs),0))


  p <- ggplot(df, aes(x=ranking, y = prob_diffs, size = as.numeric(effectif), color = new_or_old)) +
    geom_point(aes(text = paste('siret:',siret)))

  ggplotly(p)
}
