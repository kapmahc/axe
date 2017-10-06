package com.github.kapmahc.axe.nut.services;

import com.github.kapmahc.axe.nut.forms.InstallForm;
import com.github.kapmahc.axe.nut.forms.users.ResetPasswordForm;
import com.github.kapmahc.axe.nut.forms.users.SignInForm;
import com.github.kapmahc.axe.nut.forms.users.SignUpForm;
import com.github.kapmahc.axe.nut.helper.JwtHelper;
import com.github.kapmahc.axe.nut.helper.SecurityHelper;
import com.github.kapmahc.axe.nut.models.Log;
import com.github.kapmahc.axe.nut.models.Policy;
import com.github.kapmahc.axe.nut.models.Role;
import com.github.kapmahc.axe.nut.models.User;
import com.github.kapmahc.axe.nut.repositories.LogRepository;
import com.github.kapmahc.axe.nut.repositories.PolicyRepository;
import com.github.kapmahc.axe.nut.repositories.RoleRepository;
import com.github.kapmahc.axe.nut.repositories.UserRepository;
import org.springframework.context.MessageSource;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Propagation;
import org.springframework.transaction.annotation.Transactional;

import javax.annotation.Resource;
import java.io.UnsupportedEncodingException;
import java.security.NoSuchAlgorithmException;
import java.time.ZoneId;
import java.util.Date;
import java.util.Locale;
import java.util.Map;
import java.util.UUID;

@Service("nut.userService")
@Transactional(readOnly = true)
public class UserService {
    @Transactional(propagation = Propagation.REQUIRED)
    public void resetPassword(Locale locale, String ip, String token, ResetPasswordForm form) {
        Map<String, String> claim = jwtHelper.parse(token);
        User user = userRepository.findByUid(claim.get("uid"));
        if (user == null) {
            throw new IllegalArgumentException(messageSource.getMessage("nut.errors.email-not-exist", null, locale));
        }
        user.setPassword(securityHelper.password(form.getPassword()));
        log(user, ip, messageSource.getMessage("nut.logs.reset-password", null, locale));
    }

    @Transactional(propagation = Propagation.REQUIRED)
    public void unlock(Locale locale, String ip, String token) {
        Map<String, String> claim = jwtHelper.parse(token);
        User user = userRepository.findByUid(claim.get("uid"));
        if (user == null) {
            throw new IllegalArgumentException(messageSource.getMessage("nut.errors.email-not-exist", null, locale));
        } else if (user.getConfirmedAt() != null) {
            throw new IllegalArgumentException(messageSource.getMessage("nut.errors.not-lock", null, locale));
        }
        user.setLockedAt(null);
        userRepository.save(user);
        log(user, ip, messageSource.getMessage("nut.logs.unlock", null, locale));
    }

    @Transactional(propagation = Propagation.REQUIRED)
    public void confirm(Locale locale, String ip, String token) {
        Map<String, String> claim = jwtHelper.parse(token);
        User user = userRepository.findByUid(claim.get("uid"));
        if (user == null) {
            throw new IllegalArgumentException(messageSource.getMessage("nut.errors.email-not-exist", null, locale));
        } else if (user.getConfirmedAt() != null) {
            throw new IllegalArgumentException(messageSource.getMessage("nut.errors.already-confirm", null, locale));
        }
        user.setConfirmedAt(new Date());
        userRepository.save(user);
        log(user, ip, messageSource.getMessage("nut.logs.confirm", null, locale));
    }

    @Transactional(propagation = Propagation.REQUIRED)
    public User signUp(Locale locale, String ip, SignUpForm form) throws NoSuchAlgorithmException, UnsupportedEncodingException {
        User user = userRepository.findByProviderTypeAndProviderId(User.Type.EMAIL, form.getEmail());
        if (user == null) {
            throw new IllegalArgumentException(messageSource.getMessage("nut.errors.email-exist", null, locale));
        }
        user = addUser(form.getName(), form.getEmail(), form.getPassword());
        log(user, ip, messageSource.getMessage("nut.log.sign-up", null, locale));
        return user;
    }

    @Transactional(propagation = Propagation.REQUIRED)
    public User signIn(Locale locale, String ip, SignInForm form) {
        User user = userRepository.findByProviderTypeAndProviderId(User.Type.EMAIL, form.getEmail());
        if (user == null) {
            throw new IllegalArgumentException(messageSource.getMessage("nut.errors.email-not-exist", null, locale));
        }
        if (!securityHelper.check(form.getPassword(), user.getPassword())) {
            String msg = messageSource.getMessage("nut.errors.bad-password", null, locale);
            log(user, ip, msg);
            throw new IllegalArgumentException(msg);
        }
        if (user.getConfirmedAt() == null) {
            throw new IllegalArgumentException(messageSource.getMessage("nut.errors.not-confirm", null, locale));
        }
        if (user.getLockedAt() != null) {
            throw new IllegalArgumentException(messageSource.getMessage("nut.errors.is-lock", null, locale));
        }
        return user;
    }

    @Transactional(propagation = Propagation.REQUIRED)
    public void install(Locale locale, String ip, InstallForm form) throws NoSuchAlgorithmException, UnsupportedEncodingException {
        localeService.set(locale, "site.title", form.getTitle());
        localeService.set(locale, "site.subhead", form.getSubhead());
        User user = addUser(form.getName(), form.getEmail(), form.getPassword());
        log(user, ip, messageSource.getMessage("nut.logs.sign-up", null, locale));
        user.setConfirmedAt(new Date());
        userRepository.save(user);
        log(user, ip, messageSource.getMessage("nut.logs.confirm", null, locale));
        for (String n : new String[]{Role.ADMIN, Role.ROOT}) {
            allow(user, n, 20);
            log(user, ip, messageSource.getMessage("nut.logs.allow", new Object[]{n, null, null}, locale));
        }
    }


    private User addUser(String name, String email, String password) throws NoSuchAlgorithmException, UnsupportedEncodingException {
        User it = new User();
        it.setName(name);
        it.setEmail(email);
        it.setProviderId(email);
        it.setUid(UUID.randomUUID().toString());
        it.setProviderType(User.Type.EMAIL);
        it.setGravatarLogo();
        it.setPassword(securityHelper.password(password));
        userRepository.save(it);
        return it;
    }

    public void log(User user, String ip, String message) {
        Log it = new Log();
        it.setIp(ip);
        it.setMessage(message);
        it.setUser(user);
        logRepository.save(it);
    }

    public void deny(User user, String role) {
        deny(user, role, null, null);
    }

    public void deny(User user, String role, String resourceType, Long resourceId) {
        Role r = getRole(role, resourceType, resourceId);
        Policy p = policyRepository.findByUserAndRole(user, r);
        if (p == null) {
            return;
        }
        policyRepository.delete(p);
    }

    public void allow(User user, String role, int years) {
        Date begin = new Date();
        ZoneId zone = ZoneId.systemDefault();
        Date end = Date.from(
                begin.toInstant()
                        .atZone(zone)
                        .toLocalDateTime()
                        .plusYears(years)
                        .atZone(zone)
                        .toInstant()
        );
        allow(user, role, null, null, new Date(), end);
    }

    public void allow(User user, String role, String resourceType, Long resourceId, Date begin, Date end) {
        Role r = getRole(role, resourceType, resourceId);
        Policy p = policyRepository.findByUserAndRole(user, r);
        if (p == null) {
            p = new Policy();
            p.setRole(r);
            p.setUser(user);
        }
        p.setBegin(begin);
        p.setEnd(end);
        policyRepository.save(p);
    }

    private Role getRole(String name, String resourceType, Long resourceId) {
        Role it = roleRepository.findByNameAndResourceTypeAndResourceId(name, resourceType, resourceId);
        if (it != null) {
            return it;
        }
        it = new Role();
        it.setName(name);
        it.setResourceType(resourceType);
        it.setResourceId(resourceId);
        roleRepository.save(it);
        return it;
    }

    @Resource
    UserRepository userRepository;
    @Resource
    RoleRepository roleRepository;
    @Resource
    PolicyRepository policyRepository;
    @Resource
    LogRepository logRepository;
    @Resource
    LocaleService localeService;
    @Resource
    MessageSource messageSource;
    @Resource
    SecurityHelper securityHelper;
    @Resource
    JwtHelper jwtHelper;
}
