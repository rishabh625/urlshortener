package handler

import "net/http"

func RedirectHandler() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// Implement
		http.RedirectHandler("", http.StatusPermanentRedirect)
	})
}

func GenerateShortURL() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodPost:
			// Implement
		default:
			writer.WriteHeader(http.StatusBadRequest)
		}

	})
}
