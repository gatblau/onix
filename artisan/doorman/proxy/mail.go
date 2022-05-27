/*
  Onix Config Manager - Artisan's Doorman Proxy
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package main

import (
    "crypto/tls"
    "github.com/gatblau/onix/artisan/doorman/core"
    gomail "gopkg.in/mail.v2"
    "strings"
)

func sendMail(notification core.NotificationMsg) error {
    // gets the sender information
    from, err := getEmailFrom()
    if err != nil {
        return err
    }
    smtpUser, err := getSmtpUser()
    if err != nil {
        return err
    }
    smtpPwd, err := getSmtpPwd()
    if err != nil {
        return err
    }
    // gets the receiver email address
    to := strings.Split(notification.Recipient, ",")
    // gets smtp server configuration
    smtpHost, err := getSmtpHost()
    if err != nil {
        return err
    }
    smtpPort, err := getSmtpPort()
    if err != nil {
        return err
    }

    m := gomail.NewMessage()
    m.SetHeader("From", from)
    m.SetHeader("To", to...)
    m.SetHeader("Subject", notification.Subject)
    if strings.Contains(notification.Content, "<html>") {
        m.SetBody("text/html", notification.Content)
    } else {
        m.SetBody("text/plain", notification.Content)
    }
    d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPwd)
    // This is only needed when SSL/TLS certificate is not valid on server.
    // In production this should be set to false
    d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
    // send email
    if dialErr := d.DialAndSend(m); dialErr != nil {
        return dialErr
    }
    return nil
}
