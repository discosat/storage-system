package dam

import "log"

func QueryParser(requestData map[string]interface{}) bool {
	log.Println("Initiating query parser with request data: ", requestData)

	if imgID, ok := requestData["img_id"].(string); ok && imgID == "1" {
		log.Println("Returning id: ", imgID)
		return true
	}

	return false
}
