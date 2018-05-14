#
# This is a Shiny web application. You can run the application by clicking
# the 'Run App' button above.
#
# Find out more about building applications with Shiny here:
#
#    http://shiny.rstudio.com/
#

library(shiny)
library(ggplot2)
library(dplyr)

#################################################################################################
# UI
#################################################################################################

ui <- fluidPage(# Application title
  titlePanel("Séries temporelles"),

  sidebarLayout(
    sidebarPanel(
      textInput(
        inputId = "siret",
        value = "44250732300040",
        label = "SIRET"
      ),

      selectInput(
        "feature",
        "Feature to plot:",
        choices = c(
          "effectif",
          "mean_cotisation_due",
          "nb_debits",
          "heures_consommees",
          "ratio_dettecumulee_cotisation",
          "montant_part_ouvriere",
          "montant_part_patronale",
          "poids_frng",
          "taux_marge",
          "delai_fournisseur",
          "dette_fiscale",
          "financier_ct",
          "financier"

        )
      )
    ),

    # Show a plot of the generated distribution
    mainPanel(plotOutput("time_series"),
              # plotOutput("histogram"),
              htmlOutput("text_information"))
  ))

#################################################################################################
# SERVER
#################################################################################################
server <- function(input, output) {

  # PLOT 1
  output$time_series <- renderPlot({
    my_data <- table_wholesample %>%
      filter(siret == input$siret) %>%
      mutate(time = as.Date(periode))

    plot <- ggplot(my_data,
                   aes_string(x = "time", y = input$feature)) +
      geom_point() +
      geom_smooth(se = FALSE) +
      xlab("Année") +
      ylab(input$feature) +
      theme(text = element_text(size = 20))

    my_data_slice <-  my_data %>%
      slice(1)
    if (!is.na(my_data_slice$date_defaillance)) {
      plot <- plot +
        geom_vline(xintercept = as.Date(my_data_slice$date_defaillance),
                   colour = '#AF0000')
    }
    plot
  })

  # PLOT 2
#
#   output$histogram <- renderPlot({
#     my_data <- table_wholesample %>%
#       group_by(siret) %>%
#       top_n(n = 1,wt = as.Date(periode))
#
#     plot <- ggplot(my_data,
#                    aes_string(x = input$feature)) +
#       geom_histogram()
#       # geom_vline(xintercept = (my_data %>% filter(siret == input$siret))[,input$feature]) +
#       # theme(text = element_text(size = 20))
#     plot
#   })

  # TEXT
  output$text_information <- renderText({
    table <- table_wholesample %>%
      filter(siret == input$siret) %>%
      slice(1)

    raison_sociale <- table %>%
      select(raison_sociale) %>%
      as.character() %>%
      paste("Raison sociale:", .)

    procedure_collective <-
      ifelse(
        is.na(table$date_defaillance),
        "Pas de procédure collective connue",
        paste("Procédure collective:", '<p style="color:#AF0000">',as.character(table$date_defaillance),'</p>')
      )


    departement <- paste("Departement:",table$code_departement)

    APE <- paste('APE:',table$code_ape, "niveau 1:",table$libelle_naf_niveau1)



    HTML(paste(
      raison_sociale,
      departement,
      procedure_collective,
      APE,
      sep = '<br/>'))

  })
}

# Run the application
shinyApp(ui = ui, server = server)
