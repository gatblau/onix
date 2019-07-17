package features;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.ComponentScan;
import org.springframework.context.annotation.Configuration;

@Configuration
@ComponentScan(basePackages = { "features.*, org.gatblau.*" })
public class Config {
    @Bean
    public Util util() {
        return new Util();
    }
}