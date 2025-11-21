package models

type UserRequestDto struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserUpdateNameDTO struct {
	Name string `json:"name"`
}

func (u UserUpdateNameDTO) Validate() error {
	if len(u.Name) < 4 {
		return ErrShortName
	}
	return nil
}

func (u UserRequestDto) Validate() error {

	switch {
	case len(u.Name) < 4:
		return ErrShortName

	case len(u.Name) > 35:
		return ErrLongName

	case len(u.Password) < 8:
		return ErrShortPass

	case len(u.Password) > 18:
		return ErrLongPass

	default:
		return nil
	}
}
