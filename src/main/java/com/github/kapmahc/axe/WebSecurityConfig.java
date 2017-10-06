package com.github.kapmahc.axe;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpMethod;
import org.springframework.security.config.annotation.authentication.builders.AuthenticationManagerBuilder;
import org.springframework.security.config.annotation.method.configuration.EnableGlobalMethodSecurity;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configuration.WebSecurityConfigurerAdapter;
import org.springframework.security.config.http.SessionCreationPolicy;

@Configuration
@EnableWebSecurity
@EnableGlobalMethodSecurity(proxyTargetClass = true, prePostEnabled = true)
public class WebSecurityConfig extends WebSecurityConfigurerAdapter {
    @Override
    protected void configure(HttpSecurity http) throws Exception {

        http.sessionManagement()
                .invalidSessionUrl("/users/sign-in")
                .sessionCreationPolicy(SessionCreationPolicy.IF_REQUIRED)
                .maximumSessions(3);

        http
                .authorizeRequests()
                .antMatchers(
                        HttpMethod.GET,
                        "/",

"/users/logs",
                        "/install",
                        "/users/sign-in",
                        "/users/sign-up",
                        "/users/confirm",
                        "/users/confirm/{token}",
                        "/users/unlock",
                        "/users/unlock/{token}",
                        "/users/forgot-password",
                        "/users/reset-password/{token}",

                        "/assets/**"
                ).permitAll()
                .antMatchers(
                        HttpMethod.POST,
                        "/install",
                        "/users/sign-in",
                        "/users/sign-up",
                        "/users/confirm",
                        "/users/unlock",
                        "/users/forgot-password",
                        "/users/reset-password/{token}"
                ).permitAll()
                .antMatchers(
                        HttpMethod.DELETE,
                        "/users/sign-up"
                ).permitAll()
                .anyRequest().authenticated()
                .and().formLogin().loginPage("/users/sign-in").failureForwardUrl("/users/sign-in")
                .and().logout().logoutUrl("/users/sign-out").logoutSuccessUrl("/users/sign-in");
    }

    @Autowired
    public void configureGlobal(AuthenticationManagerBuilder auth) throws Exception {
        auth
                .inMemoryAuthentication()
                .withUser("user").password("password").roles("USER");
    }


}
