package creator

import "github.com/grig-iv/anki-card-creator/anki"

func CheckModel() (bool, error) {
	names, err := anki.ModelNames()
	if err != nil {
		return false, err
	}

	for _, n := range names {
		if n == modelName {
			return true, nil
		}
	}

	return false, nil
}
