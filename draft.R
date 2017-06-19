
table_naf <- import_table_naf(
  path = "data-raw/insee/naf/naf2008_5_niveaux.xls"
)

table_sirene <- dplyr::left_join(
  x = table_sirene,
  y = table_naf,
  by = c("code_apet_700" = "code_naf_niveau5")
) %>%
  dplyr::select(-nic)
