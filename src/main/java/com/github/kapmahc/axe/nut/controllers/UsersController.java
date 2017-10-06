package com.github.kapmahc.axe.nut.controllers;

import com.github.kapmahc.axe.nut.forms.users.EmailForm;
import com.github.kapmahc.axe.nut.forms.users.ResetPasswordForm;
import com.github.kapmahc.axe.nut.forms.users.SignInForm;
import com.github.kapmahc.axe.nut.forms.users.SignUpForm;
import com.github.kapmahc.axe.nut.helper.EmailJobSender;
import com.github.kapmahc.axe.nut.helper.JwtHelper;
import com.github.kapmahc.axe.nut.helper.RequestHelper;
import com.github.kapmahc.axe.nut.models.User;
import com.github.kapmahc.axe.nut.repositories.UserRepository;
import com.github.kapmahc.axe.nut.services.UserService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.context.MessageSource;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.validation.BindingResult;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.servlet.mvc.support.RedirectAttributes;

import javax.annotation.Resource;
import javax.servlet.http.HttpServletRequest;
import javax.validation.Valid;
import java.io.IOException;
import java.io.UnsupportedEncodingException;
import java.security.NoSuchAlgorithmException;
import java.security.Principal;
import java.time.Duration;
import java.util.HashMap;
import java.util.Locale;
import java.util.Map;

import static com.github.kapmahc.axe.Flash.ERROR;
import static com.github.kapmahc.axe.Flash.NOTICE;
import static com.github.kapmahc.axe.nut.services.UserService.*;

@Controller("nut.usersController")
@RequestMapping(value = "/users")
public class UsersController {
    @GetMapping("/logs")
    public String getLogs(Principal principal) {
        return "nut/users/logs";
    }


    // ------------------------------------------

    // http://www.thymeleaf.org/doc/articles/springsecurity.html
    @GetMapping("/sign-in")
    public String getSignIn(SignInForm signInForm) {
        return "nut/users/sign-in";
    }

    // ------------------------------------------
    @GetMapping("/sign-up")
    public String getSignUp(SignUpForm form) {
        return "nut/users/sign-up";
    }

    @PostMapping("/sign-up")
    public String postSignUp(@Valid SignUpForm form, BindingResult result, final RedirectAttributes attributes, Locale locale, HttpServletRequest request) throws NoSuchAlgorithmException, UnsupportedEncodingException {
        if (requestHelper.check(result, attributes)) {
            if (form.getPassword().equals(form.getPasswordConfirmation())) {
                String ip = requestHelper.clientIp(request);
                try {
                    User user = userService.signUp(locale, ip, form);
                    sendEmail(request, locale, user, ACTION_CONFIRM);
                    attributes.addFlashAttribute(NOTICE, messageSource.getMessage("nut.users.confirm.notice", null, locale));
                    return "redirect:/users/sign-in";
                } catch (Exception e) {
                    e.printStackTrace();
                    attributes.addFlashAttribute(ERROR, e.getMessage());
                }
            } else {
                attributes.addFlashAttribute(ERROR, messageSource.getMessage("validators.passwords-not-match", null, locale));
            }
        }
        return "redirect:/users/sign-up";
    }

    @GetMapping("/confirm")
    public String getConfirm(EmailForm emailForm, Model model, Locale locale) {
        String act = ACTION_CONFIRM;
        model.addAttribute("action", act);
        model.addAttribute("title", messageSource.getMessage("nut.users." + act + ".title", null, locale));
        return "nut/users/email-form";
    }

    // ------------------------------------------
    @PostMapping("/confirm")
    public String postConfirm(HttpServletRequest request, @Valid EmailForm emailForm, BindingResult result, final RedirectAttributes attributes, Locale locale) throws IOException {
        if (requestHelper.check(result, attributes)) {
            User user = userRepository.findByProviderTypeAndProviderId(User.Type.EMAIL, emailForm.getEmail());
            if (user == null) {
                attributes.addFlashAttribute(ERROR, messageSource.getMessage("nut.errors.email-not-exist", null, locale));
            } else if (user.isConfirm()) {
                attributes.addFlashAttribute(ERROR, messageSource.getMessage("nut.errors.already-confirm", null, locale));
            } else {
                sendEmail(request, locale, user, ACTION_CONFIRM);
                attributes.addFlashAttribute(NOTICE, messageSource.getMessage("nut.users.confirm.notice", null, locale));
            }
        }
        return "redirect:/users/confirm";
    }

    @GetMapping("/confirm/{token}")
    public String getConfirm(@PathVariable("token") String token, final RedirectAttributes attributes, Locale locale, HttpServletRequest request) {

        String ip = requestHelper.clientIp(request);
        try {
            userService.confirm(locale, ip, token);
            attributes.addFlashAttribute(NOTICE, messageSource.getMessage("nut.users.confirm.success", null, locale));
        } catch (Exception e) {
            e.printStackTrace();
            attributes.addFlashAttribute(ERROR, e.getMessage());
        }

        return "redirect:/users/sign-in";
    }

    // ------------------------------------------
    @GetMapping("/unlock")
    public String getUnlock(EmailForm emailForm, Model model, Locale locale) {
        String act = ACTION_UNLOCK;
        model.addAttribute("action", act);
        model.addAttribute("title", messageSource.getMessage("nut.users." + act + ".title", null, locale));
        return "nut/users/email-form";
    }

    @PostMapping("/unlock")
    public String postUnlock(HttpServletRequest request, @Valid EmailForm emailForm, BindingResult result, final RedirectAttributes attributes, Locale locale) throws IOException {
        if (requestHelper.check(result, attributes)) {
            User user = userRepository.findByProviderTypeAndProviderId(User.Type.EMAIL, emailForm.getEmail());
            if (user == null) {
                attributes.addFlashAttribute(ERROR, messageSource.getMessage("nut.errors.email-not-exist", null, locale));
            } else if (!user.isLock()) {
                attributes.addFlashAttribute(ERROR, messageSource.getMessage("nut.errors.not-lock", null, locale));
            } else {
                sendEmail(request, locale, user, ACTION_UNLOCK);
                attributes.addFlashAttribute(NOTICE, messageSource.getMessage("nut.users.unlock.notice", null, locale));
            }
        }
        return "redirect:/users/unlock";
    }

    @GetMapping("/unlock/{token}")
    public String getUnlock(@PathVariable("token") String token, final RedirectAttributes attributes, Locale locale, HttpServletRequest request) {
        String ip = requestHelper.clientIp(request);
        try {
            userService.unlock(locale, ip, token);
            attributes.addFlashAttribute(NOTICE, messageSource.getMessage("nut.users.unlock.success", null, locale));
        } catch (Exception e) {
            e.printStackTrace();
            attributes.addFlashAttribute(ERROR, e.getMessage());
        }

        return "redirect:/users/sign-in";
    }

    // ------------------------------------------
    @GetMapping("/forgot-password")
    public String getForgotPassword(EmailForm emailForm, Model model, Locale locale) {
        String act = "forgot-password";
        model.addAttribute("action", act);
        model.addAttribute("title", messageSource.getMessage("nut.users." + act + ".title", null, locale));
        return "nut/users/email-form";
    }

    @PostMapping("/forgot-password")
    public String postForgotPassword(HttpServletRequest request, @Valid EmailForm emailForm, BindingResult result, final RedirectAttributes attributes, Locale locale) throws IOException {
        if (requestHelper.check(result, attributes)) {
            User user = userRepository.findByProviderTypeAndProviderId(User.Type.EMAIL, emailForm.getEmail());
            if (user == null) {
                attributes.addFlashAttribute(ERROR, messageSource.getMessage("nut.errors.email-not-exist", null, locale));
            } else {
                sendEmail(request, locale, user, ACTION_RESET_PASSWORD);
                attributes.addFlashAttribute(NOTICE, messageSource.getMessage("nut.users.forgot-password.notice", null, locale));
            }
        }
        return "redirect:/users/forgot-password";
    }

    // ------------------------------------------
    @GetMapping("/reset-password/{token}")
    public String getResetPassword(ResetPasswordForm resetPasswordForm) {
        return "nut/users/reset-password";
    }

    @PostMapping("/reset-password/{token}")
    public String postResetPassword(@PathVariable("token") String token, @Valid ResetPasswordForm resetPasswordForm, BindingResult result, final RedirectAttributes attributes, Locale locale, HttpServletRequest request) {
        if (requestHelper.check(result, attributes)) {
            if (resetPasswordForm.getPassword().equals(resetPasswordForm.getPasswordConfirmation())) {
                String ip = requestHelper.clientIp(request);
                try {
                    userService.resetPassword(locale, ip, token, resetPasswordForm);
                    attributes.addFlashAttribute(NOTICE, messageSource.getMessage("nut.users.reset-password.success", null, locale));
                    return "redirect:/users/sign-in";
                } catch (Exception e) {
                    e.printStackTrace();
                    attributes.addFlashAttribute(ERROR, e.getMessage());
                }
            } else {
                attributes.addFlashAttribute(ERROR, messageSource.getMessage("validators.passwords-not-match", null, locale));
            }
        }

        return "redirect:/users/reset-password/" + token;
    }

    // ------------------------------------------
    private void sendEmail(HttpServletRequest request, Locale locale, User user, String action) throws IOException {
        Map<String, String> claim = new HashMap<>();
        claim.put("uid", user.getUid());
        claim.put("act", action);
        String token = jwtHelper.generate(claim, Duration.ofHours(1));


        String subject = messageSource.getMessage(
                "nut.emails." + action + ".subject",
                null,
                locale
        );
        String body = messageSource.getMessage(
                "nut.emails." + action + ".body",
                new Object[]{requestHelper.home(request), token},
                locale
        );
        emailJobSender.send(user.getEmail(), subject, body);
    }

    @Resource
    UserRepository userRepository;
    @Resource
    UserService userService;
    @Resource
    EmailJobSender emailJobSender;
    @Resource
    JwtHelper jwtHelper;
    @Resource
    MessageSource messageSource;

    @Resource
    RequestHelper requestHelper;


    private final static Logger logger = LoggerFactory.getLogger(UsersController.class);
}
