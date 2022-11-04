package main

import (
	"log"
	"strings"
	"testing"

	"github.com/chrj/smtpd"
	"github.com/emersion/go-msgauth/dkim"
	"github.com/emersion/go-msgauth/dmarc"
)

func TestMailHandler(t *testing.T) {
	env := smtpd.Envelope{
		Sender:     "hermanbanken@gmail.com",
		Recipients: []string{"herman@ptsm.hermanbanken.nl"},
		Data:       []byte("Received: from mail-ed1-f43.google.com ([209.85.208.43]) by localhost.localdomain\r\n\twith ESMTP;\r\n\tFri, 21 Oct 2022 14:59:59 +0000 (UTC)\r\nReceived: by mail-ed1-f43.google.com with SMTP id m16so6874440edc.4\n        for <herman@ptsm.hermanbanken.nl>; Fri, 21 Oct 2022 07:59:59 -0700 (PDT)\nDKIM-Signature: v=1; a=rsa-sha256; c=relaxed/relaxed;\n        d=gmail.com; s=20210112;\n        h=to:subject:message-id:date:from:mime-version:from:to:cc:subject\n         :date:message-id:reply-to;\n        bh=JImGYayboS+aa+UFjQ2Gcm2T5RCchiP3q2SJjhSlt64=;\n        b=kAFNBq9UfhEZ/MmTdxaWFXZH7MRE+zQt99skGwZc+YY6baL7nza/eRfnwpr4aQKmeI\n         UMaILR9s+GDSx0RdwLZx6wNh5tOEVDdDTV16potBEQe9r6opOjlDABfNt9B38vuaeYxq\n         Z8MqyeP38pQw0B9aD6U9Yu9/s2wG3a3tTJa1qNvQOlbkgSSdVvRNJM11Auzp4VPxxGN/\n         HZYVbDA727cErv84DPW6QKQWfyZURnH8NVcdDg1coapZsPz8l+LVNcAY2CAYUzaErEET\n         aMC8lR4X5cEbQn14+12lXsMP6SotNgY6B2AHgKdnur7w5e/1vO2LPFQSBayrTKexr/hd\n         5WJw==\nX-Google-DKIM-Signature: v=1; a=rsa-sha256; c=relaxed/relaxed;\n        d=1e100.net; s=20210112;\n        h=to:subject:message-id:date:from:mime-version:x-gm-message-state\n         :from:to:cc:subject:date:message-id:reply-to;\n        bh=JImGYayboS+aa+UFjQ2Gcm2T5RCchiP3q2SJjhSlt64=;\n        b=1Nf6TeO9Q0z260C3vTGck+cz97NC1KaJkF+5PviMv0X7dMFaDkLQyIIatcIMtHeQ0U\n         HC9g0hfcYtLq6RcbiWTWGGvI5fcvRBV9WLu4RRuT4pWRUSk/JE+r5lx4OrDHbRuSX+N6\n         h8ITEU0B2B0yio2stoN+/oKVU90XmGP0MCOyWMuXYYTYXUJ8ngq6Jn1o6K3JfKonyt9e\n         GwSBZbEm9tr5x0OJYByu9cvC3AUL3y3u+5o7MsfCRXs5Xuhfposur564X561qYOIC5iu\n         6P6Cpo9Vbibb4G/JzVpLwmeQJnjazmiJyL3As2G0QR8qOcWGe+wCZlUuHorgkRrA9/44\n         NIyQ==\nX-Gm-Message-State: ACrzQf3tO6ah60NFO9oJudD4UaFJxiUMY+qchnm1CFWFbAy8q5PelPfM\n\t/FIuWvmdzxji8QptjUkh0gycB22U1AQ/JxREm9PhMpO1\nX-Google-Smtp-Source: AMsMyM6OEg1WwCBzGXLdITVr+Xjrm+AnHetcjlxiCwKKDfWQ0+eAdeIP5J4LjRMNK4EaiX2AH3GnUPQ3k9bKLmkBG4w=\nX-Received: by 2002:a05:6402:501b:b0:459:df91:983 with SMTP id\n p27-20020a056402501b00b00459df910983mr17043924eda.85.1666364398530; Fri, 21\n Oct 2022 07:59:58 -0700 (PDT)\nMIME-Version: 1.0\nFrom: Herman Banken <hermanbanken@gmail.com>\nDate: Fri, 21 Oct 2022 16:59:47 +0200\nMessage-ID: <CAKZLU9s5yCEbwNigD-6KQzKod-bgcwXBe7BewyH7USrG_+TCUg@mail.gmail.com>\nSubject: Second email test\nTo: herman@ptsm.hermanbanken.nl\nContent-Type: multipart/alternative; boundary=\"0000000000000a734005eb8cb29a\"\n\n--0000000000000a734005eb8cb29a\nContent-Type: text/plain; charset=\"UTF-8\"\n\nFingerscrossed!\n\nMet vriendelijke groet,\n\nSpammer\n\n--0000000000000a734005eb8cb29a\nContent-Type: text/html; charset=\"UTF-8\"\n\n<div dir=\"ltr\">Fingerscrossed!<div><br clear=\"all\"><div><div dir=\"ltr\" class=\"gmail_signature\" data-smartmail=\"gmail_signature\"><div dir=\"ltr\"><div><div dir=\"ltr\">Met vriendelijke groet,<br><br><div>Spammer</div></div></div></div></div></div></div></div>\n\n--0000000000000a734005eb8cb29a--\n"),
	}
	r := strings.NewReader(string(env.Data))

	verifications, err := dkim.Verify(r)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range verifications {
		if v.Err == nil {
			log.Println("Valid signature for:", v.Domain)
		} else {
			t.Fail()
			log.Println("Invalid signature for:", v.Domain, v.Err)
		}
		rec, err := dmarc.Lookup(v.Domain)
		_ = err
		if rec.Policy == dmarc.PolicyReject {
			log.Println("TODO Reject email")
		}
	}
}
