package error

import "log"
/**
 * ErrPrint prints the error.
 * @param err error
 */
func ErrPrint(err error) {
	if err != nil {
		log.Printf("error: %v\n", err)
	}
}

func NumError(err error) {
	if err != nil {
		log.Printf("error: %v\n", err)
	}
}