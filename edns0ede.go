package rdns

import (
	"github.com/miekg/dns"
)

type EDNS0EDETemplate struct {
	infoCode     uint16
	textTemplate *Template
}

func NewEDNS0EDETemplate(infoCode uint16, extraText string) (*EDNS0EDETemplate, error) {
	if infoCode == 0 && extraText == "" {
		return nil, nil
	}

	tpl, err := NewTemplate(extraText)
	if err != nil {
		return nil, err
	}

	return &EDNS0EDETemplate{
		infoCode:     infoCode,
		textTemplate: tpl,
	}, nil
}

// Apply executes the template for the EDNS0-EDE record text, e.g. replacing
// placeholders in the Text with Query names, then adding the EDE record to
// the given msg.
func (t *EDNS0EDETemplate) Apply(msg, q *dns.Msg) error {
	if t == nil {
		return nil
	}
	extraText, err := t.textTemplate.Apply(q)
	if err != nil {
		return err
	}
	ede := &dns.EDNS0_EDE{
		InfoCode:  t.infoCode,
		ExtraText: extraText,
	}
	msg.SetEdns0(4096, false)
	opt := msg.IsEdns0()
	opt.Option = append(opt.Option, ede)
	return nil
}
