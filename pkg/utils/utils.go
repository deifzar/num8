package utils

import (
	"context"
	"deifzar/num8/pkg/log8"
	"encoding/json"
	"fmt"
	"net"

	"github.com/itchyny/gojq"
)

func IsValidIPAddress(ip string) bool {
	ipAddress := net.ParseIP(ip)
	return ipAddress != nil
}

func RunGoJQQuery(JSONInput any, query *gojq.Query) (any, error) {
	var JSONOutput any
	iter := query.RunWithContext(context.Background(), JSONInput)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			if err, ok := err.(*gojq.HaltError); ok && err.Value() == nil {
				log8.BaseLogger.Debug().Stack().Msg(err.Error())
				log8.BaseLogger.Warn().Msg("NO HALT Error iterating JSON object after applying 'jq' filter")
				return nil, err
			}
			log8.BaseLogger.Debug().Stack().Msg(err.Error())
			log8.BaseLogger.Warn().Msg("HALT Error iterating the Burp sitemap JSON object after applying 'jq' filter")
		}
		b, err := json.Marshal(v)
		if err != nil {
			log8.BaseLogger.Debug().Stack().Msg(err.Error())
			log8.BaseLogger.Warn().Msgf("Error with Marshal function. Origing: JSON object after applying 'jq' filter")
			return nil, err
		}
		err = json.Unmarshal(b, &JSONOutput)
		if err != nil {
			log8.BaseLogger.Debug().Stack().Msg(err.Error())
			log8.BaseLogger.Warn().Msgf("Error with Unmarshing function. Origing: JSON object after applying 'jq' filter")
			return nil, err
		}
	}
	return JSONOutput, nil
}

func EqualAny(a, b interface{}) bool {
	return fmt.Sprint(a) == fmt.Sprint(b)
}
