package com.github.kapmahc.axe.nut.security;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.security.core.AuthenticationException;
import org.springframework.stereotype.Component;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;

@Component
public class AuthenticationFailureHandler implements org.springframework.security.web.authentication.AuthenticationFailureHandler {
    @Override
    public void onAuthenticationFailure(HttpServletRequest request, HttpServletResponse response, AuthenticationException exception) throws IOException, ServletException {
        String error = exception.getMessage();
        logger.error("auth failed: ", error);
        response.sendRedirect("/users/sign-in?error=" + error);
    }

    private final static Logger logger = LoggerFactory.getLogger(AuthenticationFailureHandler.class);
}
