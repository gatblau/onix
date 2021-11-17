/*
Onix Config Manager - Copyright (c) 2018-2020 by www.gatblau.org

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Contributors to this project, hereby assign copyright in their code to the
project, to be licensed under the same terms as the rest of the code.
*/
package org.gatblau.onix;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.gatblau.onix.conf.Config;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import javax.mail.*;
import javax.mail.internet.InternetAddress;
import javax.mail.internet.MimeBodyPart;
import javax.mail.internet.MimeMessage;
import javax.mail.internet.MimeMultipart;
import java.io.BufferedReader;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.UnsupportedEncodingException;
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.util.Date;
import java.util.Properties;
import java.util.stream.Collectors;

/*
  send emails
 */
@Service
public class Mailer {
    private final String UTF8 = StandardCharsets.UTF_8.toString();

    @Value("${wapi.smtp.from.pwd}")
    private char[] smtpFromPwd;
    private final Config cfg;
    private final Logger log;

    public Mailer(Config cfg) {
        this.cfg = cfg;
        this.log = LogManager.getLogger();
    }

    public void sendHtmlEmail(String toEmail, String subject, String body) {
        try {
            Session session = getSession();
            if (!cfg.isSmtpEnabled()) {
                String msg = String.format("cannot email to %s: email sending is disabled", toEmail);
                log.atWarn().log(String.format(msg));
                throw new RuntimeException(msg);
            }
            MimeMessage msg = new MimeMessage(session);
            //set message headers
            msg.addHeader("Content-type", "text/HTML; charset=UTF-8");
            msg.addHeader("format", "flowed");
            msg.addHeader("Content-Transfer-Encoding", "8bit");
            msg.setFrom(new InternetAddress(cfg.getSmtpFromUser(), cfg.getSmtpFromUser()));
            msg.setReplyTo(InternetAddress.parse(cfg.getSmtpFromUser(), false));
            msg.setSubject(subject, UTF8);
            Multipart multipart = new MimeMultipart("alternative");
            MimeBodyPart htmlPart = new MimeBodyPart();
            htmlPart.setContent(body, "text/html; charset=utf-8");
            multipart.addBodyPart(htmlPart);
            msg.setContent(multipart);
            msg.setSentDate(new Date());
            msg.setRecipients(Message.RecipientType.TO, InternetAddress.parse(toEmail, false));
            Transport transport = session.getTransport("smtps");
            transport.connect(cfg.getSmtpHost(), cfg.getSmtpPort(), cfg.getSmtpFromUser(), cfg.getSmtpFromPwd());
            transport.sendMessage(msg, msg.getAllRecipients());
            transport.close();
            log.atDebug().log(String.format("email sent successfully to %s, subject: %s", toEmail, subject));
        } catch (Exception e) {
            log.atError().log(String.format("failed to send email to '%s' with subject '%s': %s", toEmail, subject, e.getMessage()));
            throw new RuntimeException(e);
        }
    }

    public void sendResetPwdEmail(String toEmail, String subject, String username, String jwt) {
        sendHtmlEmail(toEmail, subject, getPwdResetTokenEmailHtml(username, jwt));
    }

    public void sendPwdChangedEmail(String toEmail, String subject, String username) {
        sendHtmlEmail(toEmail, subject, getPwdChangedEmailHtml(username));
    }

    public void sendNewAccountEmail(String toEmail, String subject, String username) {
        sendHtmlEmail(toEmail, subject, getNewAccountEmailHtml(username));
    }

    private Session getSession() {
        // if the host is not set then do not create a session
        if (!cfg.isSmtpEnabled()) {
            log.atWarn().log("email sending is disabled in the configuration");
            return null;
        }
        Properties props = new Properties();
        props.put("mail.smtp.host", cfg.getSmtpHost()); //SMTP Host
        props.put("mail.smtp.port", cfg.getSmtpPort()); //TLS Port
        props.put("mail.smtp.auth", cfg.isSmtpAuth()); //enable authentication
        props.put("mail.smtp.starttls.enable", cfg.isSmtpStartTLS()); //enable STARTTLS

        //create Authenticator object to pass in Session.getInstance argument
        Authenticator auth = new Authenticator() {
            //override the getPasswordAuthentication method
            protected PasswordAuthentication getPasswordAuthentication() {
            return new PasswordAuthentication(cfg.getSmtpFromUser(), new String(smtpFromPwd));
            }
        };
        return Session.getInstance(props, auth);
    }

    private String getPwdResetTokenEmailHtml(String username, String jwtToken) {
        String encJwtToken = null;
        try {
            encJwtToken = URLEncoder.encode(jwtToken, UTF8);
        } catch (UnsupportedEncodingException e) {
            log.atError().log("failed to encode password reset url");
            throw new RuntimeException(e);
        }
        String pwdResetUri = String.format("%s?token=%s", cfg.getSmtpPwdResetURI(), encJwtToken);
        ClassLoader classloader = Thread.currentThread().getContextClassLoader();
        InputStream inputStream = classloader.getResourceAsStream("mail/en/pwdReset.html");
        String html = new BufferedReader(new InputStreamReader(inputStream)).lines().collect(Collectors.joining("\n"));
        return String.format(html, username, pwdResetUri);
    }

    private String getPwdChangedEmailHtml(String username) {
        ClassLoader classloader = Thread.currentThread().getContextClassLoader();
        InputStream inputStream = classloader.getResourceAsStream("mail/en/pwdChanged.html");
        String html = new BufferedReader(new InputStreamReader(inputStream)).lines().collect(Collectors.joining("\n"));
        return String.format(html, username);
    }

    private String getNewAccountEmailHtml(String username) {
        ClassLoader classloader = Thread.currentThread().getContextClassLoader();
        InputStream inputStream = classloader.getResourceAsStream("mail/en/newAccount.html");
        String html = new BufferedReader(new InputStreamReader(inputStream)).lines().collect(Collectors.joining("\n"));
        return String.format(html, username, cfg.getSmtpPwdSetupURI());
    }
}
