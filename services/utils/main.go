package utils

import "github.com/joho/godotenv"

func Env(path string) {

	err := godotenv.Load(path)
	if err != nil {
		panic(err)
	}
}

func Contains(slice []string, item string) bool {

	set := make(map[string]struct{}, len(slice))

	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]

	return ok
}
