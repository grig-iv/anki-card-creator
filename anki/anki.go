package deck

import "github.com/grig-iv/anki-card-creator/ankiConnect"

func CheckModel() (bool, error) {
	names, err := ankiConnect.ModelNames()
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
