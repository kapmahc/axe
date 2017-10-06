package com.github.kapmahc.axe.nut.controllers;

import com.github.kapmahc.axe.nut.forms.InstallForm;
import com.github.kapmahc.axe.nut.helper.RequestHelper;
import com.github.kapmahc.axe.nut.repositories.UserRepository;
import com.github.kapmahc.axe.nut.services.UserService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.context.MessageSource;
import org.springframework.stereotype.Controller;
import org.springframework.validation.BindingResult;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.servlet.mvc.support.RedirectAttributes;

import javax.annotation.Resource;
import javax.servlet.http.HttpServletRequest;
import javax.validation.Valid;
import java.io.UnsupportedEncodingException;
import java.security.NoSuchAlgorithmException;
import java.util.Locale;

import static com.github.kapmahc.axe.Flash.ERROR;
import static com.github.kapmahc.axe.Flash.NOTICE;

@Controller("auth.installController")
@RequestMapping(value = "/install")
public class InstallController {
    @GetMapping
    public String getInstall(InstallForm installForm) {
        checkDatabaseIsEmpty();
        return "nut/install";
    }

    @PostMapping
    public String postInstall(@Valid InstallForm installForm, BindingResult result, final RedirectAttributes attributes, Locale locale, HttpServletRequest request) throws NoSuchAlgorithmException, UnsupportedEncodingException {
        checkDatabaseIsEmpty();
        if (requestHelper.check(result, attributes)) {
            if (installForm.getPassword().equals(installForm.getPasswordConfirmation())) {
                String ip = requestHelper.clientIp(request);
                try {
                    userService.install(locale, ip, installForm);
                    attributes.addFlashAttribute(NOTICE, messageSource.getMessage("nut.install.success", null, locale));
                    return "redirect:/users/sign-in";
                } catch (Exception e) {
                    e.printStackTrace();
                    attributes.addFlashAttribute(ERROR, e.getMessage());
                }
            } else {
                attributes.addFlashAttribute(ERROR, messageSource.getMessage("validators.passwords-not-match", null, locale));
            }

        }

        return "redirect:/install";

    }

    private void checkDatabaseIsEmpty() {
        if (userRepository.count() > 0) {
            throw new IllegalArgumentException();
        }
    }

    @Resource
    RequestHelper requestHelper;
    @Resource
    UserService userService;
    @Resource
    MessageSource messageSource;
    @Resource
    UserRepository userRepository;

    private final static Logger logger = LoggerFactory.getLogger(InstallController.class);
}
