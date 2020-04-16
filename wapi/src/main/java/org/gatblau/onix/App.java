/*
Onix Config Manager - Copyright (c) 2018-2019 by www.gatblau.org

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

import org.gatblau.onix.conf.Config;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import org.springframework.retry.backoff.FixedBackOffPolicy;
import org.springframework.retry.policy.SimpleRetryPolicy;
import org.springframework.retry.support.RetryTemplate;

import javax.annotation.PostConstruct;
import javax.annotation.PreDestroy;

@SpringBootApplication
public class App {
    private final Config cfg;

    public App(Config cfg) {
        this.cfg = cfg;
    }

    public static void main(String[] args) {
        SpringApplication.run(App.class, args);
    }

    @PostConstruct
    private void onStartUp(){
    }

    @PreDestroy
    private void onShutDown(){
    }

    /*
      defines retry policies for acquiring a broker session
     */
    @Bean
    RetryTemplate retrySession() {
        RetryTemplate retryTemplate = new RetryTemplate();

        FixedBackOffPolicy fixedBackOffPolicy = new FixedBackOffPolicy();
        fixedBackOffPolicy.setBackOffPeriod(cfg.getEventsClientBackOffPeriod());
        retryTemplate.setBackOffPolicy(fixedBackOffPolicy);

        SimpleRetryPolicy retryPolicy = new SimpleRetryPolicy();
        retryPolicy.setMaxAttempts(cfg.getEventsClientRetries());
        retryTemplate.setRetryPolicy(retryPolicy);

        return retryTemplate;
    }
}
