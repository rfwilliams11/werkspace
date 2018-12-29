package mongo

import (
	"github.com/globalsign/mgo/bson"
	"strings"
	"fmt"
)

func Query(match bson.M, project []string, distinct string) []bson.M {
	completeQuery := []bson.M{match}

	if len(project) > 0 {
		completeQuery = append(completeQuery, Projection(project))
	}

	if len(distinct) > 0 {
		completeQuery = append(completeQuery, Unwind(distinct))
		completeQuery = append(completeQuery, Distinct(distinct))
	}

	return completeQuery
}

func Unwind(field string) bson.M {
	// We can only unwind to objects at the root level
	return bson.M{"$unwind":fmt.Sprintf("$%s", strings.Split(field,".")[0])}
}

func Distinct(field string) bson.M {
	return bson.M{"$group":bson.M{"_id":fmt.Sprintf("$%s", field)}}
}

func Projection(fields []string) bson.M {
	fieldsToProject := make(bson.M)

	for _, v := range fields {
		fieldsToProject[v] = 1
	}

	return bson.M{"$project":fieldsToProject}
}

func Match(fields map[string][]string) bson.M {
	if len(fields)== 0 {
		return nil
	}

	var andComponents []bson.M

	for key, values := range fields {
		if len(values) <= 0 {
			continue
		}

		splitValues := strings.Split(values[0], ",")

		var orComponents []bson.M
		for _, value := range splitValues {
			orComponents = append(orComponents, bson.M{key:value})
		}

		andComponents = append(andComponents, bson.M{"$or":orComponents})
	}

	return bson.M{"$match":  bson.M{"$and":andComponents}}
}