package webclient

import "net/url"

// mapToUrlValues Преобразует map[string][]sring в url.Values
func mapToUrlValues(data map[string][]string) url.Values {
	// result := url.Values{}
	//
	// for key, values := range data {
	// 	for _, v := range values {
	// 		result.Add(key, v)
	// 	}
	// }
	//
	// return result

	return url.Values(data)
}
