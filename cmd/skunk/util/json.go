package util

import "encoding/json"

func JsonEncode(data interface{}) (string, error) {
    jsonEnc, err := json.Marshal(data)
    if err != nil {
        return "", err
    }
    return string(jsonEnc), nil
}

func JsonDecode(jsonEnc string) (interface{}, error) {
    var jsonDec interface{}
    err := json.Unmarshal([]byte(jsonEnc), &jsonDec)
    if err != nil {
        return nil, err
    }
    return jsonDec, nil
}
