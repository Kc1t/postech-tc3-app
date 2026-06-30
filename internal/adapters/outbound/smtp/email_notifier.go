package smtpadapter

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/fiap/postech-tc1/internal/domain/entities"
)

type EmailNotifier struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func New(host, port, username, password, from string) *EmailNotifier {
	return &EmailNotifier{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (n *EmailNotifier) NotifyStatusChange(_ context.Context, toEmail, toName string, orderCode int, status entities.OrderStatus) error {
	if n.host == "" {
		return nil
	}

	subject, body, ok := buildContent(toName, orderCode, status)
	if !ok {
		return nil
	}

	msg := []byte(
		"From: " + n.from + "\r\n" +
			"To: " + toEmail + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n" +
			"\r\n" +
			body,
	)

	auth := smtp.PlainAuth("", n.username, n.password, n.host)
	return smtp.SendMail(n.host+":"+n.port, auth, n.from, []string{toEmail}, msg)
}

func buildContent(name string, code int, status entities.OrderStatus) (subject, body string, ok bool) {
	switch status {
	case entities.StatusAwaitingApproval:
		subject = fmt.Sprintf("OS #%d — Orçamento pronto para aprovação", code)
		body = fmt.Sprintf("Olá, %s!\n\nSeu orçamento da OS #%d está pronto.\nAcesse o sistema para aprovar ou recusar.\n\nObrigado,\nOficina Workshop", name, code)
	case entities.StatusInExecution:
		subject = fmt.Sprintf("OS #%d — Orçamento aprovado, serviços iniciados", code)
		body = fmt.Sprintf("Olá, %s!\n\nSeu orçamento da OS #%d foi aprovado e os serviços foram iniciados.\n\nObrigado,\nOficina Workshop", name, code)
	case entities.StatusFinished:
		subject = fmt.Sprintf("OS #%d — Veículo pronto para retirada", code)
		body = fmt.Sprintf("Olá, %s!\n\nSeu veículo está pronto! A OS #%d foi finalizada.\nPasse na oficina para retirada.\n\nObrigado,\nOficina Workshop", name, code)
	case entities.StatusDelivered:
		subject = fmt.Sprintf("OS #%d — Entrega confirmada", code)
		body = fmt.Sprintf("Olá, %s!\n\nA entrega do seu veículo referente à OS #%d foi confirmada.\n\nObrigado pela preferência!\nOficina Workshop", name, code)
	default:
		return "", "", false
	}
	return subject, body, true
}
