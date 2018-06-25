# define custom callback class

AUCPR_keras_callback <- function(){
  AUCPR_keras <- R6::R6Class("AUCPR_keras_callback",
                           inherit = KerasCallback,

                           public = list(

                             AUCPR = NULL,

                             on_epoch_end = function(epoch, logs = list()) {
                               X_val = self.validation_data[0]
                               Y_val = self.validation_data[1]
                               Y_predict = model.predict(X_val)
                               self$AUCPR <- c(self$AUCPR,
                                               AUCPR(
                                                 y_true = Y_val,
                                                 y_pred = Y_predict)
                                               )
                               cat(self$AUCPR)
                             }
                           ))

return(AUCPR$new())
}

