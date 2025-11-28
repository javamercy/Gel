package utilities

import "errors"

func Run(errors ...error) error {
	for _, err := range errors {
		if err != nil {
			return err
		}
	}
	return nil
}

func RunAll(rules ...error) error {
	var errs []error
	for _, err := range rules {
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
