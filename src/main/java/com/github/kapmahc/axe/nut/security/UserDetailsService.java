package com.github.kapmahc.axe.nut.security;

import com.github.kapmahc.axe.nut.models.Policy;
import com.github.kapmahc.axe.nut.models.Role;
import com.github.kapmahc.axe.nut.models.User;
import com.github.kapmahc.axe.nut.repositories.PolicyRepository;
import com.github.kapmahc.axe.nut.repositories.UserRepository;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.stereotype.Component;

import javax.annotation.Resource;
import java.util.HashSet;
import java.util.Set;

@Component
public class UserDetailsService implements org.springframework.security.core.userdetails.UserDetailsService {
    @Override
    public UserDetails loadUserByUsername(String email) throws UsernameNotFoundException {
        logger.debug("select user by {}", email);

        User user = userRepository.findByProviderTypeAndProviderId(User.Type.EMAIL, email);
        if (user == null) {
            throw new UsernameNotFoundException(email);
        }
        Set<GrantedAuthority> authorities = new HashSet<>();
        for (Policy p : policyRepository.findByUser(user)) {
            Role role = p.getRole();
            if (role.getResourceType() == null &&
                    role.getResourceId() == null &&
                    p.isEnable()) {
                authorities.add(new SimpleGrantedAuthority("ROLE_" + role.getName()));
            }
        }

        return new org.springframework.security.core.userdetails.User(
                user.getUid(), user.getPassword(),
                user.isConfirm(),
                true,
                true,
                !user.isLock(),
                authorities
        );
    }

    @Resource
    UserRepository userRepository;
    @Resource
    PolicyRepository policyRepository;
    private final static Logger logger = LoggerFactory.getLogger(UserDetailsService.class);
}
