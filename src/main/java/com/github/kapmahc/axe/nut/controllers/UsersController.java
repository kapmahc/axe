package com.github.kapmahc.axe.nut.controllers;

import com.github.kapmahc.axe.nut.forms.users.EmailForm;
import com.github.kapmahc.axe.nut.forms.users.ResetPasswordForm;
import com.github.kapmahc.axe.nut.forms.users.SignInForm;
import com.github.kapmahc.axe.nut.forms.users.SignUpForm;
import com.github.kapmahc.axe.nut.helper.EmailJobSender;
import com.github.kapmahc.axe.nut.helper.JwtHelper;
import com.github.kapmahc.axe.nut.models.User;
import com.github.kapmahc.axe.nut.repositories.UserRepository;
import com.github.kapmahc.axe.nut.services.UserService;
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
import java.time.Duration;
import java.util.HashMap;
import java.util.Locale;
import java.util.Map;

@Controller("nut.usersController")
@RequestMapping(value = "/users")
public class UsersController {
    @GetMapping("/sign-in")
    public String getSignIn(SignInForm form) {
        return "nut/users/sign-in";
    }

    @PostMapping("/sign-in")
    public String postSignUp(@Valid SignInForm form, BindingResult result, final RedirectAttributes attributes, Locale locale, HttpServletRequest request) {
        // TODO
        return "redirect:/users/sign-in";
    }

    @GetMapping("/sign-up")
    public String getSignUp(SignUpForm form) {
        return "nut/users/sign-up";
    }

    @PostMapping("/sign-up")
    public String postInstall(@Valid SignUpForm form, BindingResult result, final RedirectAttributes attributes, Locale locale, HttpServletRequest request) throws NoSuchAlgorithmException, UnsupportedEncodingException {
        // TODO
        return "redirect:/users/sign-up";
    }

    @GetMapping("/confirm")
    public String getConfirm(EmailForm form, Model model) {
        model.addAttribute("action", "confirm");
        return "nut/users/email-form";
    }

    @PostMapping("/confirm")
    public String postConfirm(@Valid EmailForm form, BindingResult result, final RedirectAttributes attributes, Locale locale, HttpServletRequest request) {
        // TODO
        return "redirect:/users/sign-in";
    }

    @GetMapping("/confirm/{token}")
    public String getConfirm(@PathVariable("token") String token, final RedirectAttributes attributes, Locale locale, HttpServletRequest request) {

        return "nut/users/email-form";
    }

    @GetMapping("/unlock")
    public String getUnlock(EmailForm form, Model model) {
        model.addAttribute("action", "unlock");
        return "nut/users/email-form";
    }

    @PostMapping("/unlock")
    public String postUnlock(@Valid EmailForm form, BindingResult result, final RedirectAttributes attributes, Locale locale, HttpServletRequest request) {
        // TODO
        return "redirect:/users/sign-in";
    }

    @GetMapping("/unlock/{token}")
    public String getUnlock(@PathVariable("token") String token, final RedirectAttributes attributes, Locale locale, HttpServletRequest request) {

        return "nut/users/email-form";
    }

    @GetMapping("/forgot-password")
    public String getForgotPassword(EmailForm form, Model model) {
        model.addAttribute("action", "forgot-password");
        return "nut/users/email-form";
    }

    @PostMapping("/forgot-password")
    public String postForgotPassword(@Valid EmailForm form, BindingResult result, final RedirectAttributes attributes, Locale locale, HttpServletRequest request) {
        // TODO
        return "redirect:/users/sign-in";
    }

    @GetMapping("/reset-password/{token}")
    public String getResetPassword(@PathVariable("token") String token, ResetPasswordForm form, Model model) {
        return "nut/users/reset-password";
    }

    @PostMapping("/reset-password/{token}")
    public String postResetPassword(@PathVariable("token") String token, @Valid ResetPasswordForm form, BindingResult result, final RedirectAttributes attributes, Locale locale, HttpServletRequest request) {
        // TODO
        return "redirect:/users/reset-password";
    }

    private void sendEmail(Locale locale, User user, String action) throws IOException {
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
                new Object[]{"/users/" + action + "/" + token},
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


    public final static String ACTION_CONFIRM = "confirm";
    public final static String ACTION_UNLOCK = "unlock";
    public final static String ACTION_RESET_PASSWORD = "reset-password";
}
