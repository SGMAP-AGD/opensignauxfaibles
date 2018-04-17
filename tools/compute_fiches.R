
compute_pdf <- function(inputpath) {
  library("opensignauxfaibles")
  library("readr")
  library("rmarkdown")
  df_liste <-
    read_csv(path, col_types = list(siret = col_character()))
  for (k in df_liste$siret[5:9]) {
    cat(k, "\n")
    rmarkdown::render(
      input = "inst/fiche_visite.Rmd",
      output_format = "pdf_document",
      params = list(siret = k),
      output_file = paste0("output/fiches/fiche_", k, ".pdf")
    )
  }
}
