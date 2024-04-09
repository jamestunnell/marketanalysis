package models

func ValidateParam(p Param) []error {
	val := p.GetVal()

	errs := []error{}

	for _, c := range p.Constraints() {
		if err := c.Check(val); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}
