package network

import (
	"fmt"
	"net/url"
	"reflect"
)

func formatUrlAndGdprUrl(urlFormat string, pathParameters []*PathParameter, queryParameters []*QueryParameter) (*url.URL, *url.URL, error) {
	reqUrlStr, gdprUrlStr := formatUrlAndGdprUrlStrings(urlFormat, pathParameters)
	reqUrl, err := url.Parse(reqUrlStr)
	if err != nil {
		return nil, nil, err
	}

	gdprUrl, err := url.Parse(gdprUrlStr)
	if err != nil {
		return nil, nil, err
	}

	formatQueryParameters(reqUrl, gdprUrl, queryParameters)
	return reqUrl, gdprUrl, nil
}

func formatUrlAndGdprUrlStrings(urlFormat string, pathParameters []*PathParameter) (string, string) {
	if len(pathParameters) == 0 {
		return urlFormat, urlFormat
	}

	reqUrlStr := urlFormat
	gdprUrlStr := urlFormat

	size := len(pathParameters)
	values := make([]interface{}, size)
	gdprValues := make([]interface{}, size)

	for i, parameter := range pathParameters {
		switch reflect.TypeOf(parameter.Value).Kind() {
		case reflect.Ptr:
			if reflect.ValueOf(parameter.Value).IsNil() {
				continue
			}

			values[i] = url.QueryEscape(fmt.Sprintf("%v", reflect.ValueOf(parameter.Value).Elem()))

			if parameter.Sensitive {
				gdprValues[i] = sensitivePlaceholderValue
			} else {
				gdprValues[i] = values[i]
			}

		case reflect.String:
			values[i] = url.QueryEscape(parameter.Value.(string))
			gdprValues[i] = parameter.gdprValue().(string)

		default:
			values[i] = url.QueryEscape(fmt.Sprintf("%v", parameter.Value))
			gdprValues[i] = fmt.Sprintf("%v", parameter.gdprValue())
		}
	}

	reqUrlStr = fmt.Sprintf(reqUrlStr, values...)
	gdprUrlStr = fmt.Sprintf(gdprUrlStr, gdprValues...)
	return reqUrlStr, gdprUrlStr
}

func formatQueryParameters(reqUrl *url.URL, gdprUrl *url.URL, queryParameters []*QueryParameter) {
	if len(queryParameters) == 0 {
		return
	}

	reqUrlQuery := reqUrl.Query()
	gdprUrlQuery := gdprUrl.Query()

	for _, queryParameter := range queryParameters {
		if queryParameter == nil || queryParameter.Value == nil {
			continue
		}

		var (
			value     string
			gdprValue string
		)

		switch reflect.TypeOf(queryParameter.Value).Kind() {
		case reflect.Ptr:
			if reflect.ValueOf(queryParameter.Value).IsNil() {
				continue
			}

			value = fmt.Sprintf("%v", reflect.ValueOf(queryParameter.Value).Elem())

			if queryParameter.Sensitive {
				gdprValue = sensitivePlaceholderValue
			} else {
				gdprValue = value
			}

		case reflect.String:
			value = queryParameter.Value.(string)
			gdprValue = queryParameter.gdprValue().(string)

		default:
			value = fmt.Sprintf("%v", queryParameter.Value)
			gdprValue = fmt.Sprintf("%v", queryParameter.gdprValue())
		}

		if len(value) == 0 {
			continue
		}

		reqUrlQuery.Set(queryParameter.Name, value)
		gdprUrlQuery.Set(queryParameter.Name, gdprValue)
	}

	reqUrl.RawQuery = reqUrlQuery.Encode()
	gdprUrl.RawQuery = gdprUrlQuery.Encode()
}

