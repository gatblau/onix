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
import javax.mail.internet.MimeMessage;
import java.util.Date;
import java.util.Properties;

/*
  send emails
 */
@Service
public class Mailer {
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
            MimeMessage msg = new MimeMessage(session);
            //set message headers
            msg.addHeader("Content-type", "text/HTML; charset=UTF-8");
            msg.addHeader("format", "flowed");
            msg.addHeader("Content-Transfer-Encoding", "8bit");
            msg.setFrom(new InternetAddress("no_reply@onix.com", "NoReply-Onix"));
            msg.setReplyTo(InternetAddress.parse("no_reply@onix.com", false));
            msg.setSubject(subject, "UTF-8");
            msg.setText(body, "UTF-8");
            msg.setSentDate(new Date());
            msg.setRecipients(Message.RecipientType.TO, InternetAddress.parse(toEmail, false));
            log.atDebug().log("message is ready");
            Transport.send(msg);
            log.atDebug().log(String.format("email sent successfully to %s, subject: %s", toEmail, subject));
        } catch (Exception e) {
            log.atError().log(String.format("failed to send email to '%s' with subject '%s': %s", toEmail, subject, e.getMessage()));
        }
    }

    private Session getSession() {
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
}
