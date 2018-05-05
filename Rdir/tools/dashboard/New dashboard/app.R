#########################
library(semantic.dashboard)
library(ggplot2)
library(dplyr)
#########################

# Look at selectize.

shiny_data <- table_wholesample
shiny_siret <- shiny_data %>% select(siret,raison_sociale)

#################################################################################################
# UI ############################################################################################
#################################################################################################
ui <- dashboardPage(
  dashboardHeader(title = "App signaux faibles"),

  dashboardSidebar(sidebarMenu(
    menuItem(tabName = "series", text = "Séries temporelles", icon = icon("home")),
    menuItem(tabName = "another", text = "Another Tab", icon = icon("heart"))
  )),
  dashboardBody(
    fluidRow(
      box(plotOutput("time_series")),
          # plotOutput("histogram"),
      box(htmlOutput("text_information")),
      box(
        title = "Controls",
        selectizeInput(
          inputId = "siret",
          label = "Siret ou raison sociale",
          choices = NULL
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
      )
    )
  )
)

#################################################################################################
# SERVER ########################################################################################
#################################################################################################
server <- function(input, output,session) {

  # SIRET selectize
  updateSelectizeInput(
    session,
    'siret',
    choices = shiny_siret %>% mutate(siret_rs = paste(raison_sociale, siret)),
    server = TRUE,
    options = list(
      valueField = 'siret',
      labelField = 'siret_rs',
      searchField = 'siret_rs',
      create = FALSE,
      maxItems = 1,
      searchConjunction = 'and',
      openOnFocus = FALSE,
      maxOptions = 10,
      render = I("{
        option: function(item,escape){
          return '<div>' +item.raison_sociale + ' <i>' + item.siret + '</i></div>'
        }
      }")
    )
  )



  # PLOT 1
  output$time_series <- renderPlot({
    my_data <- shiny_data %>%
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
    if (!is.na(my_data_slice$date_effet)) {
      plot <- plot +
        geom_vline(xintercept = as.Date(my_data_slice$date_effet),
                   colour = '#AF0000')
    }
    plot
  })

  # TEXT
  output$text_information <- renderText({
    table <- shiny_data %>%
      filter(siret == input$siret) %>%
      slice(1)

    raison_sociale <- table %>%
      select(raison_sociale) %>%
      as.character() %>%
      paste("Raison sociale:", .)

    siret <- paste('<i>',table$siret,'</i>')

    procedure_collective <-
      ifelse(
        is.na(table$date_effet),
        "Pas de procédure collective connue",
        paste("Procédure collective:", '<p style="color:#AF0000">',as.character(table$date_effet),'</p>')
      )


    departement <- paste("Departement:",table$code_departement)

    APE <- paste('APE:',table$code_ape,"<br/> <em> niveau 1 </em>:",table$libelle_naf_niveau1, "<br/> <em> niveau 5 </em>:",table$libelle_naf_niveau5)



    HTML(paste(
      raison_sociale,
      siret,
      departement,
      procedure_collective,
      APE,
      sep = '<br/>'))

  })
}

shinyApp(ui, server)
