package types

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"github.com/gotd/td/tg"
	"io"
	"strings"
)

type Contact struct {
	FirstName string
	LastName  string
	Phone     string
	Thumb     []byte
}

// ContactFromUser converts a tg.User to Contact
// Pass nil as image if it doesn't exist
func ContactFromUser(user *tg.User, image []byte) Contact {
	return Contact{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Thumb:     image,
	}
}

// AppendAsVCF writes the output of
func (c Contact) AppendAsVCF(w io.Writer) error {
	_, err := w.Write([]byte("BEGIN:VCARD\nVERSION:2.1\nN;CHARSET=UTF-8;ENCODING=QUOTED-PRINTABLE:"))
	if err != nil {
		return err
	}
	err = asUTF8Quoted(c.LastName, w)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(";"))
	if err != nil {
		return err
	}
	err = asUTF8Quoted(c.FirstName, w)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(";;;\nFN;CHARSET=UTF-8;ENCODING=QUOTED-PRINTABLE:"))
	if err != nil {
		return err
	}
	err = asUTF8Quoted(c.FirstName+" "+c.LastName, w)
	if err != nil {
		return err
	}
	err = writePhoneNumber(c.Phone, w)
	if err != nil {
		return err
	}
	err = writeProfilePhoto(c.Thumb, w)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("\nEND:VCARD\n"))
	return err
}

// asUTF8Quoted writes the input to output as UTF0-8 quoted
func asUTF8Quoted(in string, out io.Writer) error {
	for i := 0; i < len(in); i++ {
		var inByte = [1]byte{in[i]}
		var outBytes [3]byte
		hex.Encode(outBytes[1:], inByte[:])
		outBytes[0] = '='
		_, err := out.Write(bytes.ToUpper(outBytes[:]))
		if err != nil {
			return err
		}
	}
	return nil
}

// writePhoneNumber writes the PhoneNumber number of a user to output
// If the PhoneNumber number is "", it does nothing
func writePhoneNumber(phone string, out io.Writer) error {
	if phone == "" {
		return nil
	}
	// Write the prefix
	_, err := out.Write([]byte("\nTEL;CELL:"))
	if err != nil {
		return err
	}
	// Append a + to country code if possible
	if !strings.HasPrefix(phone, "0") && !strings.HasPrefix(phone, "+") {
		phone = "+" + phone
	}
	_, err = out.Write([]byte(phone))
	return err
}

// writeProfilePhoto writes the photo bytes to VCF file output
// It does nothing if the len(photo) is zero
func writeProfilePhoto(photo []byte, out io.Writer) error {
	if len(photo) == 0 {
		return nil
	}
	_, err := out.Write([]byte("\nPHOTO;ENCODING=BASE64;TYPE=JPEG:"))
	if err != nil {
		return err
	}
	encoder := base64.NewEncoder(base64.StdEncoding, out)
	_, err = encoder.Write(photo)
	if err != nil {
		return err
	}
	_ = encoder.Close()
	_, err = out.Write([]byte("\n"))
	return err
}
