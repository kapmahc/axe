package com.github.kapmahc.axe.nut.security;

import com.github.kapmahc.axe.nut.helper.RequestHelper;
import com.github.kapmahc.axe.nut.models.User;
import com.github.kapmahc.axe.nut.repositories.UserRepository;
import com.github.kapmahc.axe.nut.services.UserService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.security.core.Authentication;
import org.springframework.stereotype.Component;

import javax.annotation.Resource;
import javax.servlet.ServletException;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.util.Locale;

@Component
public class AuthenticationSuccessHandler implements org.springframework.security.web.authentication.AuthenticationSuccessHandler {
    @Override
    public void onAuthenticationSuccess(HttpServletRequest request, HttpServletResponse response, Authentication authentication) throws IOException, ServletException {
        User user = userRepository.findByUid(authentication.getName());
        logger.info("sign in success: {}@{}", user.getProviderId(), user.getProviderType());

        request.getSession().setAttribute("name", user.getName());
        Locale locale = request.getLocale();
        String ip = requestHelper.clientIp(request);
        userService.signIn(user, locale, ip);

        response.sendRedirect("/");
    }

    @Resource
    UserRepository userRepository;
    @Resource
    UserService userService;
    @Resource
    RequestHelper requestHelper;
    private final static Logger logger = LoggerFactory.getLogger(AuthenticationSuccessHandler.class);
}
