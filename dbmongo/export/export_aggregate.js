db = db.getSiblingDB('opensignauxfaibles')
db.getCollection("Features").aggregate(
	[
		{ 
			"$match" : {
				"_id.batch" : batch , 
				"_id.algo" : algo
			}
		}, 
		{ 
			"$unwind" : {
				"path" : "$value"
			}
		}, 
	   {
 $replaceRoot: { newRoot: "$value" } 	
},
		{ 
			"$project" : {
				"_id" : 0.0
			}
		}, 
		{ 
			"$out" : "to_export"
		}
	], 
	{ 
		"allowDiskUse" : false
	}
);

