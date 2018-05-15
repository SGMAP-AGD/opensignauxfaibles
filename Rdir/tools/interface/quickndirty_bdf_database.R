quickndirty_bdf_database <- function(){
  dbconnection <- mongo(collection = 'Entreprise', db = 'opensignauxfaibles', verbose = TRUE, url = 'mongodb://localhost:27017')

  data <- dbconnection$aggregate(
    '[
    {
    "$addFields": {
    "batch": {"$objectToArray": "$value.batch"}
    }
    },
    {
    "$project":{
    "value.batch":0
    }
    },
     {
     "$unwind": {
     "path" : "$batch",
     "includeArrayIndex" : "arrayIndex",
     "preserveNullAndEmptyArrays" : false
     }
     },
{
      "$project": {
        "batch.v.compact" : 0
      }
    },
     {
     "$addFields": {
     "types": {"$objectToArray": "$batch.v"}
     }
     },
    {
     "$project": {
     "batch": 0
     }
     },
     {
     "$unwind":{
     "path":"$types",
     "includeArrayIndex":"arrayIndex",
     "preserveNullAndEmptyArrays" : false
     }
     },
     {
     "$addFields":{
     "data":{
     "$objectToArray":"$types.v"
     }
     }
     },
     {
     "$project":{
     "types.v" : 0
     }
     },
     {
     "$unwind":{
     "path":"$data",
    "includeArrayIndex":"arrayIndex",
     "preserveNullAndEmptyArrays":false
     }
   },
     {
    "$addFields":{
     "value":"$data.v"
     }
     },
     {
     "$project":{
     "value":1
     }
     }
     ]')$value

  return(data)
}
